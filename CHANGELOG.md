## Unreleased

### Added

- **Patter Core Engine Gateway & Omnichannel Campaign Engine in Go 1.26**
  - **Evolution Go Decoupled Architecture**: High-performance Go microservice with Gin framework (`cmd/server/main.go` and `pkg/`).
  - **White-Label Instance & WhatsApp Proxy Engine (`pkg/instance/`, `pkg/whatsapp/`)**: Session lifecycle management, Base64 QR Code / Pairing Code generation (`GET /api/v1/gateway/whatsapp/qrcode`), status tracking (`WORKING`, `DISCONNECTED`), and webhook interception.
  - **Omnichannel Campaign Engine (`pkg/campaign/`)**:
    - **`AudienceResolver`**: Multi-criteria contact filtering (`contactCategoryIds`, `AND`/`OR` match logic, required tags, excluded tags, opt-out check, and strict E.164 phone validation).
    - **Anti-Spam Jitter Calculator**: Cryptographically secure random delay generator (`calculateJitter`) with ±30% random delay variation.
    - **Multi-Account Round-Robin Session Rotation**: Automatic session rotation across multi-account active sessions (`SessionNames`).
    - **Dynamic Variable Interpolation**: Placeholder string replacement for `{{first_name}}`, `{{company}}`, `{{email}}`, etc.
    - **Multi-Channel Dispatch**: Support for `WHATSAPP` (Evolution Go proxy), `SMS`/`MMS` (multi-providers), and `VOICE` (Twilio/Telnyx AI Agent call dialer).
    - **Real-Time Lifecycle & Metrics**: Status tracking (`DRAFT`, `SCHEDULED`, `RUNNING`, `PAUSED`, `COMPLETED`, `FAILED`) and live counters (`totalContacts`, `sentCount`, `deliveredCount`, `readCount`, `failedCount`).
  - **Interactive Swaggo OpenAPI Documentation**: Self-documenting controllers exposed at `/swagger/index.html`.
  - **Embedded Web Dashboard (`web/index.html`)**: Interactive Tailwind CSS & Alpine.js frontend embedded in the Gin server for managing white-label instances, WhatsApp sessions, QR code scanning, campaign launches, and real-time execution statistics.
  - **Technical Documentation (`docs/technical_doc.md`)**: Comprehensive technical architecture documentation covering Evolution Go transposition, campaign algorithms, API specs, and dashboard usage.

- **Fish Audio voice provider suite** — TTS (S2.1-Pro / S2-Pro) and ASR land in
  both SDKs, wired end-to-end with real pricing. One `FISH_AUDIO_API_KEY`
  covers both directions of the call:
  - **`FishAudioTTS`**: HTTP streaming synthesis via `POST /v1/tts`, default
    model `s2.1-pro` (83 languages, multi-speaker via `<|speaker:N|>` markers,
    natural-language `[bracket]` expression control). Output defaults to
    PCM_S16LE @ 16 kHz so chunks feed the pipeline with no transcoding, and
    `latency` defaults to `balanced` (~300 ms). Fish selects the model with a
    **request header**, so switching generations is a constructor argument.
    `for_twilio()` requests PCM directly at 8 kHz — the carrier wire rate — so
    the pipeline skips the resample (Fish has no native G.711, so the μ-law
    encode still runs); `for_telnyx()` keeps the 16 kHz pipeline rate. Every
    unset knob is omitted from the request so Fish applies its own documented
    default rather than one the SDK guessed.
  - **`FishAudioWebSocketTTS`**: opt-in low-latency transport over
    `wss://api.fish.audio/v1/tts/live`, pairing with `s2-pro`'s ~100 ms
    time-to-first-audio and skipping the per-utterance HTTP setup. Fish serves
    **only `s1` and `s2-pro`** on that socket, so the constructor rejects
    `s2.1-pro` up front with a pointer back to the HTTP adapter instead of
    failing mid-call. Frames are MessagePack (`start`→`text`→`flush`→`stop`),
    which is the sole reason the `[fish_audio]` extra pulls in `ormsgpack`
    (TypeScript: `@msgpack/msgpack`, an optional dependency).
  - **`FishAudioSTT`**: batch transcription via `POST /v1/asr`. Fish exposes no
    streaming socket, so the adapter buffers ~2 s windows and emits one final
    transcript each — no interim partials, same shape as `WhisperSTT`. A short
    tail is silence-padded up to Fish's documented 1-second floor rather than
    dropped, so the last words of an utterance are never lost. Server errors
    are logged and yield no transcript; they never raise into a live call.
  - **Pricing** (official, docs.fish.audio): TTS $15.00 per 1M UTF-8 **bytes**
    (`s2.1-pro-free` free), ASR $0.36 per audio hour. The adapters meter UTF-8
    byte length rather than character count — exact for latin scripts and
    correctly ~3× higher for CJK.
  `libraries/python/getpatter/providers/fish_audio_{tts,ws_tts,stt}.py`,
  `libraries/typescript/src/providers/fish-audio-{tts,ws-tts,stt}.ts`,
  `tts/fish_audio.*`, `stt/fish_audio.*`, `pricing.*`, `telemetry/stack.*`,
  init-wizard catalogues, and docs pages under
  `docs/{python,typescript}-sdk/providers/fish-audio-{tts,stt}.mdx`. Beta:
  validated against the Fish Audio API spec with real local HTTP/WebSocket
  servers in the test suite; not yet exercised on a live call.

- **xAI (Grok) voice provider suite** — all three xAI Voice APIs land in both
  SDKs, wired end-to-end with real pricing:
  - **`XaiRealtime` engine (Grok Voice Agent API)**: real-time speech-to-speech
    over `wss://api.x.ai/v1/realtime` (default model `grok-voice-latest`, voice
    `eve`). The protocol is OpenAI-Realtime-GA-compatible, so the adapter
    subclasses the GA adapter and inherits audio streaming, barge-in, and the
    full tool-calling bridge (`transfer_call`/`end_call` included). xAI-specific
    session knobs are exposed opt-in: `reasoning_effort` (`"high"`/`"none"`),
    VAD `threshold`/`prefix_padding_ms`/`idle_timeout_ms`, ASR
    `language_hint`/`keyterms`, output `speed`, pronunciation `replace` map,
    session `resumption`, and raw `server_tools` passthrough for xAI
    server-side tools (`web_search`, `x_search`, `mcp`, `file_search`).
  - **`XaiSTT`**: streaming speech-to-text over `wss://api.x.ai/v1/stt`
    (binary frames, `transcript.created` handshake, chunk/utterance finals,
    Smart Turn end-of-turn detection, keyterm biasing, diarization) plus a
    batch `transcribe()` helper for `POST /v1/stt` (word timestamps, ITN
    formatting, 25 languages).
  - **`XaiTTS`**: one-shot streaming synthesis via `POST /v1/tts` with the
    26-voice Grok roster, `for_twilio()` (native G.711 µ-law 8 kHz) /
    `for_telnyx()` (PCM16 16 kHz) carrier presets, and a
    `create_custom_voice()` helper for `POST /v1/custom-voices` voice cloning.
  - **Pricing** (official, docs.x.ai): Realtime $0.05/min, TTS $15.00/1M chars,
    STT $0.20/hr streaming ($0.10/hr batch). `calculate_realtime_cost` /
    `calculateRealtimeCost` gain an optional duration argument to support
    per-minute realtime billing (token-based path unchanged).
  `libraries/python/getpatter/providers/xai_{stt,tts,realtime}.py`,
  `libraries/typescript/src/providers/xai-{stt,tts,realtime}.ts`,
  `engines/xai.*`, `pricing.*`, docs pages under
  `docs/{python,typescript}-sdk/providers/xai-*.mdx`. Beta: validated against
  the xAI API spec and mocked protocol tests; not yet exercised on a live call.

### Fixed

- **`gemini-3.1-flash-live-preview` is now actually usable** (field-debugged on
  the Pillar demo line; fixed in both SDKs):
  - Setup: `enableAffectiveDialog`/`proactivity` are dated-2.5-native-audio-only
    knobs — on the flash-live family the server never sends `setupComplete`
    ("session ready timeout"). The TypeScript adapter now drops them (with a
    warning) on flash-live models. No-op in Python, which does not expose these
    knobs and so never sent them.
  - Tool calls: in the default BLOCKING function-call mode 3.1 halts audio at
    the `toolCall` and never resumes after the response (silent stall until the
    caller hangs up; hit on real calls where the model tools mid-conversation).
    Both SDKs now opt flash-live into the async mode (`behavior: NON_BLOCKING`
    on the declaration + `scheduling: INTERRUPT` on the response), and log the
    tool round-trip at info.
  - Adds `GeminiLiveModel.LIVE_3_1_FLASH_PREVIEW` to the Python model enum,
    matching the exported `GEMINI_LIVE_3_1_FLASH_PREVIEW` constant in TypeScript.
  `libraries/typescript/src/providers/gemini-live.ts`,
  `libraries/python/getpatter/providers/gemini_live.py`.

## 0.7.1 (2026-07-03)

### Added

- **Opt-in destination policy for the built-in `transfer_call` tool**
  (hardening for GH issue #205 — prompt-injected toll fraud). The transfer
  destination is chosen by the LLM, which is steerable by caller speech, and
  the only gate was E.164 *format* — a successful prompt injection could
  direct a billable outbound leg to any attacker-chosen (premium-rate)
  number. New `Patter.agent()` options `transfer_allowed_numbers` /
  `transfer_allowed_prefixes` (Py) — `transferAllowedNumbers` /
  `transferAllowedPrefixes` (TS) — restrict destinations to an exact-number
  allowlist and/or E.164 prefix allowlist (union). Enforced at every guard
  site BEFORE the carrier REST call (OpenAI Realtime `function_call`,
  Pipeline built-in handler, ElevenLabs ConvAI client tool), rejecting with
  the standard `{"error": ..., "status": "rejected"}` envelope so the agent
  keeps the call. Unset (default) preserves today's behaviour byte-identical;
  an empty list denies all transfers. Allowlist entries are validated at
  `agent()` construction (fail fast on typos).
  `libraries/python/getpatter/{client,models,stream_handler}.py`,
  `telephony/common.py`; `libraries/typescript/src/{client,types,stream-handler}.ts`.
- **Opt-in wall-clock outbound audio pacing (pipeline mode).** `Agent.paced_output`
  (Py) / `AgentOptions.pacedOutput` (TS) routes outbound carrier audio through a
  fixed 20 ms / 160-byte (μ-law) frame grid instead of the default event-driven
  burst send. A precomputed silence frame (μ-law `0xFF`×160 / PCM16 zero) fills
  gaps so an empty queue never stalls the clock and the first tick pays no encode
  cost (no cold-start latency), and a drift-corrected monotonic-deadline loop
  keeps the stream wall-clock-aligned — eliminating post-pause bursts. On a
  barge-in the pacer's queued backlog is dropped in the same beat as the carrier
  flush. Default off ⇒ byte-identical to today. New `getpatter/audio/pacer.py` /
  `src/audio/pacer.ts` (`OutboundFramePacer`), wired in `stream_handler`.
- **Krisp noise-cancellation denoiser selection (bring-your-own-license).**
  `Agent.denoiser` (Py) / `AgentOptions.denoiser` (TS) selects a Krisp model by
  its stable id — `"krisp-viva-tel-v2"` (VIVA telephony, NC session) or
  `"krisp-bvc-o-pro-v3"` (Background Voice Cancellation, BVC session) — resolved
  to an `AudioFilter` in the existing pre-STT chain slot. Ships **zero** Krisp
  binaries/models: it loads the operator's own `krisp_audio` SDK +
  `KRISP_VIVA_SDK_LICENSE_KEY` + `KRISP_MODELS_DIR`, with clear fail-fast errors
  when the SDK/license/model is absent. `denoiser` and `audio_filter` are parallel
  opt-ins (explicit `audio_filter` wins). New `providers/denoiser.py` /
  `providers/denoiser.ts` registry.
- **Text-based semantic end-of-turn detector (NAMO Turn Detector v1, Apache-2.0).**
  New `NamoTurnDetector` (`providers/namo_turn_detector.py` / `.ts`) implements the
  existing `TurnDetectorProvider` over the rolling conversation transcript (ONNX,
  runtime-loaded, `onnxruntime` + HF tokenizer as optional deps). The `predict`
  contract gained a backward-compatible keyword-only `transcript` argument, so the
  audio-native smart-turn detector is unchanged; the pipeline now assembles the
  last few turns + in-flight utterance and passes them on each end-of-turn check.
  Opt-in via `agent.turn_detector`; default behaviour unchanged. The tokenizer
  loader accepts a `revision=` argument (`PATTER_NAMO_REVISION` env var), so a
  Hugging Face Hub repo id `tokenizer_path` can be pinned to an immutable
  commit/tag rather than a mutable branch.
- **Opt-in barge-in re-delivery of the un-heard remainder (pipeline mode).**
  `Agent.redeliver_interrupted` (Py) / `AgentOptions.redeliverInterrupted` (TS):
  when the agent is interrupted mid-answer and the un-heard remainder still
  matters, the next turn gets a one-shot system nudge so the LLM answers the
  caller's new message first and then, only if still relevant, resumes the
  unfinished idea naturally (rather than silently dropping it). The
  worth-resuming decision is a pluggable `RedeliveryPolicy` (default: gate on a
  non-trivial un-heard remainder). New `services/redelivery.py` / `.ts`. No-op in
  realtime mode. Default off ⇒ zero behaviour change.
- **Soniox v5 real-time STT adopted (both SDKs).** `SonioxModel` gains
  `STT_RT_V5 = "stt-rt-v5"` and the default model flips `stt-rt-v4` →
  `stt-rt-v5` (the GA v5 model, released 2026-06-16); `v4`/`v3`/`v2` stay
  reachable for back-compat. v5 is fully API-compatible (same WS endpoint,
  in-band `api_key`, request/token-stream shape), so this is a model-id swap
  plus two new opt-in v5 endpoint controls (`endpoint_sensitivity` /
  `endpointSensitivity` in `[-1.0, 1.0]`; `endpoint_latency_adjustment_level` /
  `endpointLatencyAdjustmentLevel` in `0|1|2|3`), omitted from the wire when
  unset so existing behavior is byte-identical. Real-time rate unchanged at
  $0.12/hr ≈ $0.002/min (https://soniox.com/pricing).
- **Soniox TTS provider (both SDKs).** `SonioxTTS` over the REST one-shot
  `https://tts-rt.soniox.com/tts` bytes endpoint (`tts-rt-v1`, default voice
  `Adrian`, shares the `SONIOX_API_KEY` credential with Soniox STT).
  `forTwilio` / `forTelnyx` emit `pcm_mulaw` @ 8 kHz natively for carrier-wire
  passthrough (no resampling). Wired end-to-end: package-root export
  (`SonioxTTS`), `providers.soniox_tts()` / `sonioxTts()` config helper,
  `_create_tts_from_config` `"soniox_tts"` branch, and a pricing entry modeled
  per-1k-chars from the $0.70/hr headline (~$0.013/1k chars; native billing is
  token-based — https://soniox.com/pricing). Reuses the existing `[soniox]`
  extra (aiohttp).
- **Sarvam AI TTS provider for Indian languages (both SDKs).** `SarvamTTS` over
  the REST `https://api.sarvam.ai/text-to-speech` endpoint — Bulbul v3 (default)
  / v2 across 11 Indian languages (Hindi, Bengali, Tamil, Telugu, Kannada,
  Malayalam, Marathi, Gujarati, Punjabi, Odia, Indian English) plus code-mixed
  text. `forTwilio` / `forTelnyx` emit `mulaw` @ 8 kHz natively (carrier-wire
  passthrough). Wired end-to-end: package-root export (`SarvamTTS`),
  `providers.sarvam()` config helper, `_create_tts_from_config` `"sarvam"`
  branch, a new `[sarvam]` pyproject extra (aiohttp), and a per-1k-chars pricing
  entry — Bulbul v3 ≈ $0.036/1k (Rs 30/10k), v2 ≈ $0.018/1k (Rs 15/10k); INR is
  authoritative, USD indicative (https://www.sarvam.ai/api-pricing).
- **Engine markers shipped: `GeminiLive`, `GeminiCascade`, and a native
  `InworldRealtime`.** 0.7.0 exported the Gemini pipeline factory + adapter but
  not the `GeminiLive` engine marker, so `new GeminiLive({...})` / the engine
  import failed at runtime; and there was no native Inworld realtime engine.
  Now `import { GeminiLive, GeminiCascade, InworldRealtime } from "getpatter"`
  resolves at runtime. The richer `GeminiLiveAdapter` (mulaw8↔PCM telephony
  transcode, apiVersion auto-detect, `GEMINI_LIVE_3_1_FLASH_PREVIEW`) plus
  `GeminiSTT`/`GeminiTTS` are integrated; `InworldRealtimeAdapter` subclasses the
  OpenAI Realtime adapter and overrides only the transport
  (`wss://api.inworld.ai/v1/realtime`, Bearer/JWT) via Inworld's OpenAI-Realtime
  migration path. Backward compatible — OpenAIRealtime / OpenAIRealtime2 /
  Ultravox / ConvAI unchanged.

### Fixed

- **Gemini Live spoke the model's "thinking" before the reply** (native-audio,
  every turn). The Live setup now always sends `thinkingConfig` with a voice
  default of `thinkingBudget: 0` (OFF), opt-in via `thinking`/`thinkingBudget`
  on `GeminiLiveOptions`; and the adapter defensively drops any `thought===true`
  part from BOTH the audio stream and the text transcript (never logged).
- **`gemini-3.1-flash-live-preview` dropped the call with a silent "session
  ready timeout".** It is a native-audio model served only on `v1alpha`, but its
  id lacks the `native-audio` token so the heuristic chose `v1beta`. New
  `geminiRequiresV1Alpha()` selects `v1alpha` for native-audio AND
  `flash-live-preview` (explicit `apiVersion` overrides); connect failures now
  throw a clear, actionable error (model id + resolved apiVersion + root causes)
  instead of dead air.

### Changed

- **Barge-in "heard prefix" is now word-accurate.** On an interruption, the
  conversation history is truncated to what the caller actually heard down to a
  word boundary mid-sentence, instead of rounding up to whole sentences from a
  byte estimate. Per-sentence segments now carry their playout duration, and the
  played position is taken from the best available source — Twilio
  carrier-confirmed per-sentence marks, else the wall-clock pacer's emitted
  position (when `paced_output`), else the byte estimate. An internal heard/unsaid
  split backs the new re-delivery feature. Telnyx/Plivo (no reliable marks) use
  the pacer/estimate tiers.
- **`@google/genai` peerDependency bumped `^0.3.0` → `>=2.0.0`** (kept optional)
  for the 2.x unified SDK, so `npm install getpatter` alongside
  `@google/genai@^2.x` no longer needs `--legacy-peer-deps`.
- **`create-getpatter` launcher realigned to the SDK version (lockstep).** The
  launcher pins the `getpatter` version it bootstraps to its own version, so it
  now tracks the SDK again (`npm create getpatter` provisions the current SDK).

### Security

- **Media-stream WebSockets are now authenticated (both SDKs, all carriers) —
  fixes #204.** The carrier webhook HTTP routes were signature-validated, but
  the media-stream WS endpoints (`/ws/stream/…` Twilio, `/ws/telnyx/stream/…`,
  `/ws/plivo/stream/…`) accepted any peer with an attacker-chosen `call_id`
  (the only control was a DoS per-IP cap). Anyone who could reach the public
  host could open the socket, send a `start` frame, and drive a full
  STT→LLM→TTS session on the operator's provider keys (toll fraud) or converse
  to extract the agent's system prompt + tool list. Now the signature-validated
  webhook (which builds the stream URL/TwiML) mints a high-entropy per-call
  token, delivers it on each carrier's existing custom channel (Twilio
  `<Parameter>` → `customParameters`, Telnyx URL query, Plivo `extra_headers`),
  and the WS handler constant-time-validates it **before** opening any provider
  session — an unauthenticated peer is closed (WS 1008) with no provider connect
  and no TTS. Fail-closed by default via the new `require_stream_auth` /
  `requireStreamAuth` option (opt-out for operators serving custom TwiML, which
  then logs a warning). The token is never logged. The webhook mints only for a
  legitimate (signed) carrier, so the token is a shared secret an unsigned
  attacker cannot obtain. Backward compatible for the standard `serve()`
  inbound + outbound path (the SDK mints/embeds/validates transparently). Also
  hardened the forgeable-`X-Forwarded-For` per-IP cap with a global concurrent-
  WS backstop, documented as DoS-cap-only now that auth is enforced separately.

## 0.7.0 (2026-06-30)
