# otelcol-genai-safe

**OpenTelemetry Collector Processor for GenAI** — privacy-by-default redaction, cost/token metrics extraction, and simple runaway/loop detection.

> Drop this processor into any OTel Collector pipeline and instantly harden your GenAI telemetry — no application code changes required.
>
> ---
>
> ## Why collector-side?
>
> | Concern | App-side SDK | Collector Processor (this) |
> |---|---|---|
> | Prompt/completion redaction | Every team re-invents it | One config, all services |
> | Secret leakage (API keys in spans) | Easy to miss | Denylist regex catches it |
> | Token/cost metrics | Scattered implementations | Normalised, unified |
> | Runaway agent detection | Not usually done | Built-in loop detector |
> | Deployment | Redeploy every app | Config change, no redeploy |
>
> ---
>
> ## Features (MVP)
>
> ### 1. Redaction / Hardening
> - **Modes:** `drop`, `hash`, `hash_and_preview`, `truncate`
> - - SHA-256 hash of original value (with configurable salt) for correlation
>   - - Preview first N characters (configurable) for debugging
>     - - Denylist regex scrub across ALL string attributes (catches leaked API keys, bearer tokens, etc.)
>      
>       - ### 2. Token & Cost Metrics Extraction
>       - - Reads `gen_ai.usage.prompt_tokens`, `gen_ai.usage.completion_tokens` (and `llm.*` variants)
>         - - Normalises into `genai.tokens.prompt`, `genai.tokens.completion`, `genai.tokens.total`
>           - - Extracts `genai.cost.usd` when present
>             - - Periodic summary logging (real metric emission in week 2)
>              
>               - ### 3. Runaway / Loop Detection
>               - - Regex-based tool span name matching
>                 - - Counts repeated tool/function calls per scope batch
>                   - - Flags `genai.risk.loop_suspected = true` on all spans when threshold exceeded
>                    
>                     - ---
>
> ## Quick Start
>
> ### 1. Clone & build
>
> ```bash
> git clone https://github.com/nostalgicskinco/OpenTelemetry-Collector-Processor-GenAI-Perfect.git
> cd OpenTelemetry-Collector-Processor-GenAI-Perfect
> go mod tidy
> go build -o otelcol-genai-safe ./cmd/otelcol-custom
> ```
>
> ### 2. Run with Docker Compose (demo stack)
>
> ```bash
> cd examples
> docker compose up --build
> ```
>
> This starts:
> - **OTel Collector** with `genaisafe` processor on ports 4317/4318
> - - **Jaeger** UI at http://localhost:16686
>   - - **Prometheus** at http://localhost:9090 scraping collector metrics
>    
>     - ### 3. Send a test trace
>    
>     - ```bash
>       # Using otel-cli or any OTLP sender:
>       otel-cli span \
>         --endpoint localhost:4317 \
>         --name "chat_completion" \
>         --attrs "gen_ai.prompt=Tell me the secret password for sk-abc123xyz,gen_ai.usage.prompt_tokens=150,gen_ai.usage.completion_tokens=42"
>       ```
>
> Then check Jaeger — you'll see:
> - `gen_ai.prompt` truncated/hashed
> - - `gen_ai.prompt.hash` added for correlation
>   - - Secret patterns caught by denylist
>     - - Token counts normalised
>      
>       - ---
>
> ## Configuration
>
> ```yaml
> processors:
>   genaisafe:
>     redact:
>       mode: "hash_and_preview"   # drop | hash | hash_and_preview | truncate
>       preview_chars: 48
>       salt: "change-me-in-prod"
>       keys:
>         - "gen_ai.prompt"
>         - "gen_ai.completion"
>       denylist_regex:
>         - "(?i)api[_-]?key"
>         - "(?i)sk-[a-z0-9]{20,}"
>     metrics:
>       enable: true
>       emit_interval: "10s"
>       token_attr_candidates:
>         - "gen_ai.usage.prompt_tokens"
>         - "gen_ai.usage.completion_tokens"
>     loop_detection:
>       enable: true
>       tool_span_name_regex: "(?i)tool|function|call"
>       repeat_threshold: 6
> ```
>
> ---
>
> ## Before / After Example
>
> **Before (raw span attributes):**
> ```
> gen_ai.prompt = "What is the API key? My key is sk-abc123secretkey456..."
> gen_ai.usage.prompt_tokens = 150
> gen_ai.usage.completion_tokens = 42
> authorization = "Bearer eyJhbGciOi..."
> ```
>
> **After (processed by genaisafe):**
> ```
> gen_ai.prompt = "What is the API key? My key is sk-abc123secr..."
> gen_ai.prompt.hash = "a3f2b8c1d4e5..."
> gen_ai.usage.prompt_tokens = 150
> gen_ai.usage.completion_tokens = 42
> genai.tokens.prompt = 150
> genai.tokens.completion = 42
> genai.tokens.total = 192
> authorization = "[REDACTED]"
> ```
>
> ---
>
> ## Repo Structure
>
> ```
> .
> ├── cmd/otelcol-custom/main.go          # Custom collector entry point
> ├── processor/genaisafeprocessor/
> │   ├── config.go                        # Configuration types
> │   ├── factory.go                       # Factory + registration
> │   ├── processor.go                     # Core processing logic
> │   ├── redact.go                        # Redaction engine
> │   ├── loop.go                          # Loop/runaway detector
> │   └── metrics.go                       # Token/cost metrics extractor
> ├── examples/
> │   ├── docker-compose.yaml              # Demo stack
> │   ├── otelcol-config.yaml              # Sample collector config
> │   └── prometheus.yaml                  # Prometheus scrape config
> ├── Dockerfile                           # Multi-stage build
> ├── go.mod
> └── README.md
> ```
>
> ---
>
> ## Threat Model
>
> This processor addresses collector-side observability security for GenAI workloads:
>
> 1. **Prompt/Completion Data Leakage** — Full prompt text in traces can leak PII, secrets, proprietary data. Redaction modes prevent this while preserving correlation hashes.
>
> 2. 2. **Secret Exposure** — API keys, bearer tokens, and other secrets accidentally included in span attributes are caught by configurable denylist regex patterns.
>   
>    3. 3. **Runaway Agent Loops** — Autonomous AI agents can enter infinite tool-calling loops. The detector flags these with `genai.risk.loop_suspected` for alerting.
>      
>       4. 4. **Cost Visibility** — Without normalised token/cost attributes, it's hard to monitor GenAI spend. The metrics extractor provides unified counters.
>         
>          5. ---
>         
>          6. ## Roadmap
>         
>          7. - [x] MVP: Redaction engine (4 modes + denylist)
> - [x] MVP: Token/cost extraction + normalisation
> - [ ] - [x] MVP: Simple loop detection
> - [ ] - [ ] Week 2: Real trace-to-metrics emission (Prometheus counters)
> - [ ] - [ ] Week 2: YAML allowlists per attribute path
> - [ ] - [ ] Week 3: Per-trace loop detection (not just per-batch)
> - [ ] - [ ] Week 3: Same-prompt-hash-repeated detection
> - [ ] - [ ] Week 4: Structured PII masking (email, phone, SSN patterns)
> - [ ] - [ ] Week 4: Connector mode (traces -> metrics pipeline)
>
> - [ ] ---
>
> - [ ] ## License
>
> - [ ] Apache License 2.0 — see [LICENSE](LICENSE).
>
> - [ ] ---
>
> - [ ] ## Contributing
>
> - [ ] PRs welcome! See the roadmap above for what's coming next. The highest-impact contribution right now is real metric emission (trace-to-metrics connector/processor).
