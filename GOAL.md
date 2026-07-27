# gokata — Goal & Methodology

## Path
~/Documents/personalProjects/gokata/

## Perfil del alumno
- **8+ años** desarrollo, **6+ años en Go** (desde 2019)
- Senior Backend Engineer
- Roles en: Globant (Disney, Zynga), Qubika (Imprint fintech), Encora, Modak, One Logic Soft, Lone Wolf, **Mercado Libre (Mercado Pago)**, Ceiba
- Stack real: Go, gRPC, REST, AWS (Lambda, ECS, Fargate, DynamoDB, Step Functions), K8s, Docker, Kafka, PostgreSQL, MySQL
- Ya trabaja con concurrencia, microservicios, event-driven, CQRS, producción real
- Inglés B2

### Lo que NO aparece en el CV (gaps para enfocar)
- Go runtime internals: scheduler (G/M/P), GC mark-sweep, memory model, escape analysis
- Profiling con pprof / trace / benchstat
- Patrones avanzados de concurrencia: singleflight, errgroup+semaphore, lock-free
- Sistemas distribuidos teóricos: Raft, CAP/PACELC, MVCC, Linearizability, Lease Read
- System design para ranking/recommendation (el target ezCater)
- Rate limiting, circuit breakers, leader election, log replication

## Metodología

### Capa 1 — Socrática (Show → Explain → Expand → Build)
- Show: código, patrón, estructura
- Explain: vos decís qué ves
- Expand: profundizo según tu nivel
- Build: construímos/modificamos

### Capa 2 — Narrativa (Story → Problem → Need → Evolution → Solution → Extend → Contrast)
- Story: contexto real
- Problem: qué falla
- Need: por qué necesitamos algo mejor
- Evolution: cómo se llegó a la solución
- Solution: aplicamos rápido
- Extend: dónde más sirve
- Contrast: enfoques opuestos

### Formato de proyecto
- `kata/N-nombre/hints.md` + `*_test.go` como vehículo
- `GOAL.md` como plan vivo

## Curriculum

### Fase 1 — Go Internals & Advanced Concurrency (fundación)
| # | Challenge | Cluster | Skills clave |
|---|-----------|---------|-------------|
| 01 | Worker Pool with Backpressure | concurrency | errgroup, worker_pool, context, heap |
| 09 | G/M/P Scheduler Deep Dive | runtime | gmp_scheduler, goroutine_lifecycle, pprof |
| 10 | GC & Memory Model Lab | runtime | gc, escape_analysis, memory_model |
| 11 | Fan-Out/Fan-In Data Pipeline | concurrency | fan_out_fan_in, pipeline, context |
| 08 | Rate Limiter + Circuit Breaker Proxy | concurrency | rate_limiter, circuit_breaker, context |
| 13 | Lock-Free Patterns with atomic | concurrency | atomic, ring_buffer, sync_map, memory_model |

### Fase 2 — Architecture, Communication & Caching
| # | Challenge | Cluster | Skills clave |
|---|-----------|---------|-------------|
| 02 | Cache Stampede Prevention with Singleflight | concurrency | singleflight, escape_analysis, sync_map |
| 03 | gRPC API Gateway with Middleware | communication | grpc, protocol_buffers, graceful_shutdown |
| 04 | DDD + Hexagonal Architecture Service | architecture | ddd, clean_architecture, repository_pattern |
| 05 | Event Sourcing + CQRS Order Service | architecture | event_sourcing, cqrs, kafka |
| 06 | Saga Choreography with Compensation | distributed_systems | saga, transactional_outbox, inbox_pattern |
| 12 | Fuzzing & Benchmarking Mastery | testing | fuzzing, benchmark, race_detector |

### Fase 3 — Production, Distributed Systems & Infrastructure
| # | Challenge | Cluster | Skills clave |
|---|-----------|---------|-------------|
| 14 | Raft KV Store (Leader Election + Log Replication) | distributed_systems | raft, distributed_locking, cap_pacelc |
| 15 | CI/CD + Terraform Deploy Pipeline | infrastructure | ci_cd, terraform, docker, kubernetes |
| 07 | Full System — DDD + ES + gRPC on K8s | infrastructure | kubernetes, docker, observability, nats |

## Progreso
| # | Challenge | Estado |
|---|-----------|--------|
| 00 | GOAL.md | ✅ definido |
| 01 | workerpool | hints + tests (sin implementar) |
| 02 | singleflight-cache | hints + tests (sin implementar) |
| 03-15 | resto | ⏳ pendientes |
