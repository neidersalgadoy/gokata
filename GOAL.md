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

## Curriculum propuesto (por validar)

### Fase 1 — Go Internals & Advanced Concurrency (fundación)
| # | Kata | Tema |
|---|------|------|
| 01 | Worker pool | errgroup + semaphore + backpressure + cancelación |
| 02 | Scheduler deep dive | G/M/P, goroutine lifecycle, network poller, `runtime` pkg |
| 03 | GC & Memory | GC mark-sweep, escape analysis, stack vs heap, `pprof` |
| 04 | Lock-free patterns | `atomic`, `sync.Map`, `sync.Pool`, CAS, false sharing |
| 05 | Benchmark & profile | `benchstat`, `pprof`, `trace`, allocation optimization |

### Fase 2 — Distributed Systems (conectar teoría con práctica real)
| # | Kata | Tema |
|---|------|------|
| 06 | Log replication | Append-only log, WAL, fsync, durability, crash recovery |
| 07 | Raft distilled | Leader election, term, quorum, split-brain, lease read |
| 08 | MVCC & Isolation | MVCC, snapshot isolation, serializable, write skew |
| 09 | Consistent caching | `singleflight`, cache stampede, write-through, write-behind |
| 10 | Rate limiting | Token bucket, sliding window, distributed rate limiter (Redis) |

### Fase 3 — System Design for ezCater (ranking y recommendation)
| # | Kata | Tema |
|---|------|------|
| 11 | Ranking service | Feature store, scoring pipeline, top-N heap, real-time vs batch |
| 12 | Recommendation system | Collaborative filtering, embedding-based, cold start |
| 13 | Event-driven SAGA | Orquestación vs coreografía, compensating transactions, Kafka |
| 14 | Full system mock | Design "restaurant ranking for ezCater" — modelo, escala, trade-offs |

## Progreso
| # | Kata | Estado |
|---|------|--------|
| 00 | GOAL.md | ✅ definido |
| 01 | Worker pool | hints + tests (sin implementar) |
| 02+ | Scheduler deep dive | ⏳ |
