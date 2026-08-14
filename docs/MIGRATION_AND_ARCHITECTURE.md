# Patter Vocal Engine : Manuel de Migration et d'Architecture Go (Golang)

Ce document décrit en détail la migration intégrale du moteur d'agent vocal open-source **Patter** (initialement en TypeScript) vers un microservice natif et ultra-performant en Go.

L'architecture, les choix de conception et la structure modulaire de ce microservice s'inspirent strictement du référentiel de référence **`evolution-foundation/evolution-go`**.

---

## 1. Architecture Globale du Système

Le microservice Go s'occupe exclusivement du routage audio bidirectionnel à très faible latence entre les opérateurs de téléphonie (Carriers) et les API de LLM Multimodaux en temps réel (OpenAI/Groq).

```
 ┌──────────────────────┐                     ┌──────────────────────┐
 │  Twilio / Telnyx /   │   WebSockets (PCM)  │  Vocal Engine (Go)   │
 │        Plivo         │◄───────────────────►│  [3-Goroutines Loop] │
 └──────────────────────┘                     └──────────┬───────────┘
                                                         │
                                                         │ WebSockets (Realtime API)
                                                         ▼
                                              ┌──────────────────────┐
                                              │  OpenAI / Groq LLM   │
                                              └──────────────────────┘
                                                         │
                                                         │ Asynchronous Tool Calls
                                                         ▼
 ┌──────────────────────┐    AMQP / RabbitMQ   ┌──────────────────────┐
 │   Frontend / SaaS    │◄────────────────────┤  Main NestJS Backend │
 │      Dashboard       │                     └──────────────────────┘
 └──────────────────────┘
```

---

## 2. Structure des Dossiers (Style `evolution-go`)

La structure du projet sépare strictement les préoccupations métier en packages Go modulaires et isolés :

```
vocal-engine/
├── cmd/
│   └── vocal-engine/
│       └── main.go              # Point d'entrée de l'application, bootstrap et fermeture propre.
├── pkg/
│   ├── config/
│   │   └── config.go            # Chargement centralisé des variables d'environnement.
│   ├── engine/
│   │   ├── agent.go             # Structure de l'agent, profils voix, et repositories d'états.
│   │   ├── carrier.go           # Abstraction unifiée et implémentations Twilio, Telnyx, Plivo.
│   │   ├── stt.go               # Abstraction et connecteurs Speech-to-Text (Deepgram, AssemblyAI).
│   │   ├── tts.go               # Abstraction et connecteurs Text-to-Speech (ElevenLabs, Cartesia, OpenAI).
│   │   ├── llm.go               # Abstraction et connecteurs LLM (OpenAI, Anthropic, Groq).
│   │   ├── fallback.go          # Gestionnaire de repli et failover séquentiel mid-call.
│   │   ├── stream.go            # Serveur WebSocket Twilio, pipeline temps-réel, barge-in et mutexes.
│   │   └── pipeline_orchestrator.go # Orchestrateur Mode Pipeline (STT -> Fallback LLM -> TTS) avec VAD.
│   ├── rabbitmq/
│   │   └── client.go            # Client AMQP asynchrone et non-bloquant pour les Tool Calls.
│   ├── telephony/
│   │   ├── call.go              # Contrôleur d'appels sortants et de génération TwiML de streaming.
│   │   └── agent_controller.go  # Contrôleur CRUD des profils d'agents et des webhooks d'appels entrants.
│   └── routes/
│       └── routes.go            # Configuration du routeur Gin, du Swagger, et des fichiers statiques.
├── docs/
│   ├── docs.go                  # Code Swagger généré par Swaggo.
│   ├── swagger.json             # Spécifications OpenAPI de l'API.
│   └── MIGRATION_AND_ARCHITECTURE.md # Ce manuel complet.
├── web/
│   └── index.html               # Dashboard SaaS d'administration en Tailwind CSS.
├── go.mod                       # Déclarations des modules et dépendances Go.
└── go.sum                       # Sommes de contrôle des dépendances.
```

---

## 3. Abstraction Unifiée et Moteurs du Système

Pour reproduire la parité complète de Patter TS, le système introduit des interfaces Go découplées et extensibles :

### A. Carrier Interface (`pkg/engine/carrier.go`)
Abstrait les opérateurs téléphoniques pour l'initialisation des appels et la génération de balises de flux XML :
```go
type Carrier interface {
	GetType() CarrierType
	InitiateCall(ctx context.Context, from, to, callbackURL string) (string, string, error)
	GenerateStreamResponse(wsURL string) (string, error)
}
```
*   **TwilioCarrier** : Génère du TwiML avec `<Connect><Stream>` pour connecter l'appel au websocket du moteur Go.
*   **TelnyxCarrier** : Génère du TeXML équivalent.
*   **PlivoCarrier** : Génère du Plivo XML équivalent avec `<Stream keepCallActive="true">`.

### B. STT & TTS Providers (`pkg/engine/stt.go`, `pkg/engine/tts.go`)
*   **STTProvider** : Gère la diffusion de flux audio PCM vers les API cloud pour en extraire des transcriptions temps réel.
*   **TTSProvider** : Convertit les flux texte générés par le LLM en blocs audio compressés (G711 u-law/PCM) prêts à être injectés dans l'appel téléphonique.

### C. Fallback Manager (`pkg/engine/fallback.go`)
Implémente la mécanique de **Fallback** d'origine de Patter TS : si l'API du LLM primaire (ex: OpenAI) subit une panne de réseau ou un crash de service en plein appel, le gestionnaire intercepte l'erreur et redirige de manière séquentielle et transparente la complétion vers les LLM secondaires définis dans la chaîne de repli (`fallbackChain` : ex: Groq, Anthropic), assurant une résilience totale :
```go
func (fm *FallbackManager) ExecuteWithFallback(ctx context.Context, primaryID string, fallbackChain []string, ...) (*LLMResponse, string, error)
```

---

## 4. Pipeline Mode vs. Realtime Mode

Le Vocal Engine Go gère de manière transparente les deux modes de traitement voix :

### 1. Le "Pipeline Mode" (`pkg/engine/pipeline_orchestrator.go`)
Exécute la chaîne séquentielle **STT -> Fallback LLM -> TTS** :
*   **Détection d'activité vocale (VAD)** : Implémente un décodeur de niveau d'énergie efficace (Root Mean Square - RMS) directement sur les trames audio PCM u-law de Twilio, évitant de surcharger les API cloud si l'utilisateur ne parle pas.
*   **Orchestration des flux** : Les trames audio validées par le VAD sont diffusées vers le STT. Une fois la transcription finale générée, l'orchestrateur interroge le LLM avec le `FallbackManager` et envoie la réponse au TTS pour streaming audio sortant.

### 2. Le "Realtime Mode" (`pkg/engine/stream.go`)
Exécute la boucle ultra-basse latence sur des connexions WebSockets persistantes avec l'API Realtime d'OpenAI/Groq :
*   **Goroutine 1 (Lecture Twilio)** : Lit les paquets WebSocket de Twilio, extrait le payload Base64, et le décode.
*   **Goroutine 2 (Lien LLM Realtime)** : Maintient le canal bidirectionnel persistant avec le modèle, envoie les blocs audio décodés, et intercepte les retours (Audio, Barge-In, Tool Calls).
*   **Goroutine 3 (Écriture Twilio & Interruption)** : Gère l'envoi asynchrone et sécurisé des trames audio vers Twilio.

---

## 5. Gestion Strict des Performances et Barge-In

Pour répondre aux contraintes exigeantes de la téléphonie temps-réel, le microservice met en œuvre des mécanismes de concurrence avancés :

### A. Gestion du Barge-In (Interruption de parole)
Si l'humain coupe la parole à l'IA pendant qu'elle parle :
1. L'API Realtime signale l'événement `input_audio_buffer.speech_started` dans la Goroutine 2.
2. Le moteur Go envoie immédiatement une commande client `response.cancel` vers le WebSocket OpenAI/Groq pour forcer le modèle à stopper instantanément sa génération.
3. Un signal de contrôle est envoyé sur `interruptionChan`.
4. La Goroutine 3 intercepte ce signal, **vide instantanément** tous les buffers et files d'attente audio sortants en cours, et transmet une instruction de contrôle XML `clear` à Twilio pour interrompre immédiatement le flux audio physique dans le combiné de l'utilisateur.

### B. Absence de Mutex Bloquant sur le Flux Audio
L'échange de trames s'appuie exclusivement sur des Go Channels dimensionnés à l'aide de patterns `select` non-bloquants, éliminant tout verrou ou goulot d'étranglement CPU.

### C. Optimisation du Garbage Collector via `sync.Pool`
Pour éviter de saturer le GC en allouant des dizaines de milliers de trames de bytes par seconde, le décodage et le traitement Base64 réutilisent de manière intensive des trames pré-allouées de `8192` octets stockées dans un pool thread-safe réutilisable :
```go
var bufferPool = sync.Pool{
	New: func() interface{} { return make([]byte, 8192) },
}
```

### D. Sécurité des Websockets (Concurrence d'Écriture)
Les connexions Gorilla WebSocket ne tolèrent pas d'écritures concurrentes. Le moteur Go assure une sécurité absolue :
*   Les écritures vers Twilio sont centralisées dans une unique goroutine d'écriture (Goroutine 3).
*   Les écritures vers l'API OpenAI Realtime sont protégées par un mutex de synchronisation d'écriture thread-safe `rtWriteMu` dans le helper `safeWriteRtMsg`.

---

## 6. Intégration RabbitMQ Asynchrone (Tool Calling)

Lorsque le LLM demande l'exécution d'un outil en plein appel, le microservice intercepte la demande :
1. L'événement `response.function_call_arguments.done` est décodé.
2. Pour ne pas perturber ni ajouter de latence sur le flux de streaming audio, le moteur Go délègue instantanément l'exécution en lançant une goroutine isolée.
3. Cette goroutine publie un message structuré persistant dans la file RabbitMQ `agent_tool_calls` :
   ```go
   go func() {
       _ = se.publisher.PublishToolCall(ctx, tenantID, callSID, toolName, arguments)
   }()
   ```
4. Votre backend NestJS principal consomme ce message AMQP, exécute le code métier de l'outil (par exemple, des requêtes SQL ou l'appel de serveurs MCP), et peut mettre à jour l'appel si nécessaire, de manière asynchrone et transparente.

---

## 7. Guide de Démarrage et de Déploiement

### Prérequis
*   Go (version 1.25+ recommandée ou Go 1.22+)
*   Une file RabbitMQ (optionnelle, mode fallback inclus)
*   Des clés d'accès Twilio et OpenAI/Groq

### Installation Locale
1. Cloner le dépôt et configurer les variables d'environnement :
   ```bash
   cp .env.example .env
   ```
2. Éditer le fichier `.env` :
   ```env
   SERVER_PORT=8080
   TWILIO_ACCOUNT_SID=AC...
   TWILIO_AUTH_TOKEN=your_auth_token
   TWILIO_PHONE_NUMBER=+15551234567
   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   RABBITMQ_QUEUE=agent_tool_calls
   REALTIME_API_URL=wss://api.openai.com/v1/realtime?model=gpt-4o-realtime-preview-2024-10-01
   REALTIME_API_KEY=sk-proj-...
   ```
3. Compiler le projet :
   ```bash
   go build -v ./...
   ```
4. Lancer les tests unitaires :
   ```bash
   go test -v ./...
   ```
5. Démarrer le microservice :
   ```bash
   go run cmd/vocal-engine/main.go
   ```
6. Accéder au tableau de bord d'administration :
   *   **SaaS Dashboard** : Ouvrir `http://localhost:8080/` dans un navigateur.
   *   **Interactive OpenAPI Swagger UI** : Naviguer vers `http://localhost:8080/swagger/index.html`.
