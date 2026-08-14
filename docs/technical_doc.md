# Technical Documentation: Patter Core Engine Gateway & Campaign Engine

This technical documentation outlines the design, architecture, and API specification of the **Patter Core Engine Gateway & Omnichannel Campaign Engine** implemented in **Go 1.26** (compatible with Go 1.24+), inspired by the modular architecture of [Evolution Go](https://github.com/evolution-foundation/evolution-go).

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
│   │   ├── controller/
│   │   ├── dto/
│   │   └── service/
│   ├── whatsapp/                    # WhatsApp Cloud & Evolution Go Whatsmeow Proxy
│   │   ├── controller/
│   │   ├── dto/
│   │   └── service/
│   ├── campaign/                    # Omnichannel Campaign Engine (WhatsApp, SMS, Voice)
│   │   ├── controller/
│   │   ├── dto/
│   │   └── service/
│   ├── voice/                       # Inbound/Outbound Telephony & AI Agent Dialer
│   │   ├── controller/
│   │   └── service/
│   ├── messaging/                   # SMS/MMS Multi-Providers Orchestration
│   │   ├── controller/
│   │   └── service/
│   └── brain/                       # AI RAG & LeMUR v3 Context Dispatcher
│       ├── controller/
│       └── service/
├── web/                             # Embedded Tailwind/Alpine.js Web Dashboard
└── docs/                            # Technical documentation & Swaggo OpenAPI specs
```

---

## 🛡️ 2. Multi-Tenant White-Label Isolation

Each client or white-label reseller is isolated via a unique `tenant_id` supplied in the request header (`X-Tenant-ID`) or query parameter.

### Tenant Context Middleware (`pkg/core/middleware.go`)

- **Header Extraction**: Intercepts `X-Tenant-ID` and `X-API-Key`.
- **Isolation Scope**: Instance IDs, WhatsApp session keys, campaign records, and metrics are strictly partitioned under `tenant_id`.

---

## 📲 3. White-Label WhatsApp Session Engine (`pkg/whatsapp/`)

Patter acts as an abstraction proxy over Evolution Go and Whatsmeow engine.

### Key Capabilities

1. **Session Lifecycle**: Connect, check status (`WORKING`, `DISCONNECTED`, `CONNECTING`), list, and disconnect sessions per tenant.
2. **White-Label QR Code & Pairing Code**:
   - Endpoint: `GET /api/v1/gateway/whatsapp/qrcode`
   - Returns Base64 QR Code string and Pairing Code (`123-456`) directly consumable by white-label frontends.
3. **Webhook Interception**: Inbound messages are intercepted, associated with `tenant_id`, and dispatched asynchronously to the tenant's CRM/Brain module.

---

## 🚀 4. Omnichannel Campaign Engine (`pkg/campaign/`)

The campaign engine is designed for high-throughput, anti-spam resilience, and flexible audience targeting across **WhatsApp**, **SMS/MMS**, and **Voice Call Agent Dialing**.

### 4.1 Dynamic Audience Resolution (`AudienceResolver`)

Filters contact lists based on configurable rules:
- **Phone Validation**: Enforces strict `E.164` format compliance (`+15551234567`).
- **Opt-Out Check**: Automatically skips contacts flagged as `opted_out`.
- **Match Logic (`AND` / `OR`)**: Evaluates `ContactCategoryIDs` and `RequiredTags`.
- **Excluded Tags**: Drops any contact possessing excluded tags.

### 4.2 Anti-Spam Jitter Calculator (`CalculateJitter`)

To prevent triggering telecom provider spam filters, the engine introduces a random delay (jitter) between dispatches:
- **Formula**: Variation of **±30%** around the specified `RandomDelaySec`.
- **Randomization**: Uses cryptographically secure random integers (`crypto/rand`).

```go
func (s *campaignService) CalculateJitter(baseDelay int) time.Duration {
	if baseDelay <= 0 {
		return 0
	}
	delta := int64(float64(baseDelay) * 0.3)
	if delta <= 0 {
		return time.Duration(baseDelay) * time.Second
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(delta*2))
	actualDelay := int64(baseDelay) - delta + n.Int64()
	return time.Duration(actualDelay) * time.Second
}
```

### 4.3 Multi-Account Round-Robin Session Rotation

When multiple sender accounts or WhatsApp instances are attached to a campaign (`SessionNames`), the engine distributes dispatches across sessions using a round-robin strategy:

```go
currentSession := cfg.SessionNames[i % sessionCount]
```

### 4.4 Variable Interpolation

Replaces dynamic placeholders in template content with contact variables:
- Template: `"Bonjour {{first_name}}, découvrez l'offre de {{company}}!"`
- Resolved: `"Bonjour Alice, découvrez l'offre de Acme Inc!"`

### 4.5 Campaign Execution Lifecycle & Metrics

- **States**: `DRAFT` ➔ `SCHEDULED` ➔ `RUNNING` ➔ `PAUSED` ➔ `COMPLETED` / `FAILED`.
- **Real-Time Counters**: `sentCount`, `deliveredCount`, `readCount`, `failedCount`.

---

## 📡 5. REST API Specifications

### Key Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/gateway/instances` | Create white-label tenant instance |
| `GET` | `/api/v1/gateway/instances` | List tenant instances |
| `POST` | `/api/v1/gateway/whatsapp/connect` | Register/connect WhatsApp session |
| `GET` | `/api/v1/gateway/whatsapp/qrcode` | Get Base64 QR Code / Pairing code |
| `GET` | `/api/v1/gateway/whatsapp/sessions` | List active WhatsApp sessions |
| `POST` | `/api/v1/gateway/campaigns` | Create and trigger omnichannel campaign |
| `POST` | `/api/v1/gateway/voice/campaign` | Trigger AI Voice campaign |
| `GET` | `/api/v1/gateway/campaigns` | List tenant campaigns |
| `GET` | `/api/v1/gateway/campaigns/:id` | Get campaign status & counters |
| `POST` | `/api/v1/gateway/campaigns/:id/pause` | Pause running campaign |
| `POST` | `/api/v1/gateway/campaigns/:id/resume` | Resume paused campaign |
| `POST` | `/api/v1/gateway/voice/call` | Initiate single voice call |
| `POST` | `/api/v1/gateway/messaging/sms` | Send single SMS/MMS |
| `POST` | `/api/v1/gateway/brain/query` | Query RAG & AI Agent Brain |

Interactive Swaggo UI is available at `/swagger/index.html`.

---

## 🖥️ 6. Embedded Web Dashboard

An embedded web dashboard is served at `/web` (and `/` redirect) built with **Tailwind CSS** and **Alpine.js**:
- **Sessions Tab**: Manage WhatsApp Evolution Go sessions, inspect QR codes and pairing codes.
- **Campaigns Tab**: Trigger campaigns, configure anti-spam jitter, monitor real-time sent/delivered/read counters, pause/resume execution.
- **Instances Tab**: View white-label tenant instance statuses.
