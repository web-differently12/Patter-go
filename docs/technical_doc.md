# Technical Documentation: Patter Core Engine Gateway, Campaign Engine & Unified RAG Router

This technical documentation outlines the design, architecture, and API specification of the **Patter Core Engine Gateway, Omnichannel Campaign Engine & Unified RAG Router** implemented in **Go 1.26** (compatible with Go 1.24+), inspired by the modular architecture of [Evolution Go](https://github.com/evolution-foundation/evolution-go).

---

## 🏛️ 1. Architecture Overview

The system follows a clean, decoupled **Controller → Service → Driver** architecture using the **Gin Web Framework**, structured for high performance, concurrent handling, and seamless white-label multi-tenancy.

### Package Directory Structure

```
patter-go/
├── cmd/
│   └── server/
│       └── main.go                  # Entrypoint, Gin setup, Middleware, Swaggo & Static UI
├── pkg/
│   ├── core/                        # AuthMiddleware, TenantContext, slog Logger, API Response
│   ├── instance/                    # Tenant White-Label Instance lifecycle
│   ├── whatsapp/                    # WhatsApp Cloud & Evolution Go Whatsmeow Proxy
│   ├── campaign/                    # Omnichannel Campaign Engine (WhatsApp, SMS, Voice)
│   ├── rag/                         # Unified RAG Router & Automatic Voice Escalation
│   ├── skills/                      # Reusable System Skills (HumanTransferSkill with oral consent)
│   ├── voice/                       # Inbound/Outbound Telephony & AI Agent Dialer
│   ├── messaging/                   # SMS/MMS Multi-Providers Orchestration
│   └── brain/                       # AI RAG & LeMUR v3 Context Dispatcher
├── web/                             # Embedded Tailwind/Alpine.js Web Dashboard
└── docs/                            # Technical documentation & Swaggo OpenAPI specs
```

---

## 🛡️ 2. Multi-Tenant White-Label Isolation

Each client or white-label reseller is isolated via a unique `tenant_id` supplied in the request header (`X-Tenant-ID`) or query parameter.

---

## 🧠 3. Unified RAG Router & Automatic Voice Escalation (`pkg/rag/`)

The **Unified RAG Router** connects to multiple vector search providers (**Built-in PGVector**, **Google Vertex Search / Cloud Discovery Engine**, **Pinecone**, **Qdrant**, or **Custom HTTP Endpoints**) and applies automatic call transfer policy when documents are missing or confidence is low (`< 0.75`).

### 3.1 Decision Flow

```
                      PROSPECT ASKS QUESTION
                                │
                                ▼
                   UNIFIED RAG VECTOR SEARCH
             (PGVector / Google Vertex / Custom API)
                                │
             ┌──────────────────┴──────────────────┐
             ▼ (Score >= 0.75)                     ▼ (Score < 0.75)
    LLM Context Injection               Instant Call Escalation
    Voice Reply Streamed (<150ms)       Trigger SIP/TwiML Transfer
                                        Briefing sent to human agent
```

### 3.2 Configuration (`KnowledgeBaseConfig` & `FallbackTransferPolicy`)

```go
type FallbackTransferPolicy struct {
	EnableAutoTransfer  bool    `json:"enable_auto_transfer"`
	ConfidenceThreshold float32 `json:"confidence_threshold"`  // Default 0.75
	TransferDestination string  `json:"transfer_destination"`  // E.164 phone or SIP URI
	TransferType        string  `json:"transfer_type"`         // "WARM" or "COLD"
}
```

---

## 🛠️ 4. System Skill: `HumanTransferSkill` (`pkg/skills/`)

The `HumanTransferSkill` exposes the `request_human_transfer` OpenAPI tool definition to the LLM and enforces a **two-step oral user consent** rule:
1. When `user_consented == false`, the skill returns `CONSENTEMENT_REQUIS` instructing the LLM to ask the prospect: *"Souhaitez-vous que je vous mette en relation avec mon collègue spécialiste ?"*
2. When `user_consented == true`, the skill executes `TransferCall` via Twilio/SIP and returns `TRANSFERT_INITIÉ`.

---

## 🚀 5. Omnichannel Campaign Engine (`pkg/campaign/`)

The campaign engine handles high-throughput dispatches across **WhatsApp**, **SMS/MMS**, and **Voice Call Agent Dialing**:
- **AudienceResolver**: Filters contacts by tags/categories, checks `opted_out` status, and enforces strict `E.164` validation (`+15551234567`).
- **Anti-Spam Jitter Calculator**: Applies ±30% random delay variation using `crypto/rand`.
- **Multi-Account Round-Robin**: Rotates dispatches across `SessionNames`.
- **Variable Interpolation**: Replaces `{{first_name}}`, `{{company}}`, `{{email}}`, etc.

---

## 📡 6. REST API Specifications

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/gateway/instances` | Create white-label tenant instance |
| `GET` | `/api/v1/gateway/whatsapp/qrcode` | Get Base64 QR Code / Pairing code |
| `POST` | `/api/v1/gateway/campaigns` | Create and trigger omnichannel campaign |
| `POST` | `/api/v1/gateway/voice/campaign` | Trigger AI Voice campaign |
| `POST` | `/api/v1/gateway/brain/rag/search` | Search Unified RAG with auto-transfer rule |
| `POST` | `/api/v1/gateway/brain/transfer` | Execute human transfer skill |

Interactive Swaggo UI is available at `/swagger/index.html`.
