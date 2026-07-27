# SYSTEM AUDIT — gokata Tutor Infrastructure

> **Date:** 2026-07-26
> **Auditor:** OpenCode
> **Scope:** Tutor infrastructure only (YAML schemas, CLI tools, types, agent constitution)
> **Excludes:** Kata implementations (student work)

---

## EXECUTIVE SUMMARY

**Rating: 🟡 CONDITIONAL PASS — 2 CRITICAL, 5 MAJOR, 3 MINOR**

The tutor scaffolding is well-designed conceptually but has a **critical data-corruption bug in the core tool** (`complete-challenge`) and **zero test coverage across all infrastructure code**. The curriculum design and skill taxonomy are excellent and industry-aligned. The system is usable only with active agent supervision.

---

## 1. INFRASTRUCTURE INVENTORY

| Component | Type | Status | Lines |
|-----------|------|--------|-------|
| `AGENTS.md` | Agent constitution | 🟢 Complete | 99 |
| `SKILLS.yaml` | Skill catalog (immutable) | 🟢 Complete | 361 |
| `PROGRESS.yaml` | Progress overlay | 🟢 Complete | 253 |
| `CHALLENGE_REGISTRY.yaml` | Challenge catalog | 🟢 Complete | 208 |
| `CHALLENGE_LOG.md` | History log | 🟡 1 entry only | 5 |
| `GOAL.md` | Study plan + profile | 🟢 Complete | 76 |
| `SURVEY.md` | Knowledge tree | 🟢 Complete | 57 |
| `cmd/validate/main.go` | Consistency linter | 🟡 No tests | 268 |
| `cmd/next-challenge/main.go` | Challenge recommender | 🟡 No tests | 104 |
| `cmd/complete-challenge/main.go` | Completion tool | 🔴 BUGGY | 84 |
| `internal/types/types.go` | Shared types | 🟢 Complete | 94 |
| `go.mod` | Module definition | 🟡 Minor | 7 |
| `kata/01-workerpool/hints.md` | Kata context | 🟢 Complete | 49 |
| `kata/01-workerpool/*_test.go` | Kata tests | 🟢 Complete | 139 |

---

## 2. 🔴 CRITICAL FINDINGS

### C1. `complete-challenge` corrupts CHALLENGE_REGISTRY.yaml

**File:** `cmd/complete-challenge/main.go:53-54`

The regex-based text manipulation marks MULTIPLE challenges as `completed` per single execution:

```go
pat := fmt.Sprintf("(?m)^  %s:\n(?:.|\\n)*?^    status: pending", regexp.QuoteMeta(id))
statusRe := regexp.MustCompile(pat)
```

The `(?:.|\n)*?` pattern unpredictably crosses blank-line boundaries between YAML entries. Empirical test: running on `01-workerpool` also set `02-singleflight-cache` to `completed`.

**Downstream bugs from same tool:**
- **No idempotency guard** — running twice on the same challenge duplicates CHALLENGE_LOG.md entries and double-adds points to PROGRESS.yaml
- **No pre-validation** — doesn't run `go test ./kata/<id>/` before marking complete; student could mark a failing kata as "done"
- **`[horas]` arg parsed but never used** (line 17 declares it, line never reads it)
- **CHALLENGE_LOG.md is plain Markdown** — no machine-readable format exists; history cannot be parsed by tools

**Solution:** Use YAML marshal/unmarshal (already available via `gopkg.in/yaml.v3`), modify struct, re-serialize. Never do text-based YAML manipulation with regex.

### C2. Zero tests on all infrastructure tooling

| Tool | Tests | Risk |
|------|-------|------|
| `cmd/validate/` | 0 | Linter bugs go undetected; corrupt data is accepted as valid |
| `cmd/next-challenge/` | 0 | Wrong recommendations silently misdirect study |
| `cmd/complete-challenge/` | 0 | Data corruption (see C1) |

All three tools write to or read from the same YAML files. Without tests, a bug in any tool propagates through the system undetected.

---

## 3. 🟡 MAJOR FINDINGS

### M1. `go.mod` marks direct dependency as indirect (cosmetic, affects tooling)

`golang.org/x/sync` is imported directly by `kata/01-workerpool/ranker.go` but was listed as `// indirect`. Fixed by `go mod tidy`. Indicates module hygiene was not reviewed.

### M2. No centralized entry point

No `Makefile`, `Taskfile.yml`, or `justfile`. Agent must remember full `go run ./cmd/X/` paths. Compounds cognitive load during the orchestration loop (AGENTS.md step 2, 5, 6).

```
# Suggested:
validate:   go run ./cmd/validate/
next:       go run ./cmd/next-challenge/
complete:   go run ./cmd/complete-challenge/
test-kata:  go test ./kata/$(ID)/ -v
```

### M3. `challenges: []` arrays in PROGRESS.yaml are never populated

Each `ProgressEntry` has a `challenges []string` field, but `complete-challenge` never appends to it. After completing a challenge, there's no trace of WHICH challenge earned points for WHICH skill. Irreversible audit trail gap.

### M4. CHALLENGE_LOG.md is human-only, not machine-readable

The agent cannot parse the history to calculate trends, averages, or cumulative time. All history is lost to automation.

**Solution:** Dual format — append to `CHALLENGE_LOG.yaml` (structured, append-only array) and regenerate `CHALLENGE_LOG.md` as a view.

### M5. No point decay mechanism

Spaced repetition is described in AGENTS.md but only partially implemented. `next-challenge` checks if reinforcement is due, but PROGRESS.yaml points never decay. A skill at 8/10 stays 8/10 forever even if unpracticed for a year. The algorithm nudge is advisory only — it doesn't correct the record.

---

## 4. ⬜ MINOR FINDINGS

### m1. AGENTS.md typo: `PROGRAMS.yaml` (line 52)

Should be `PROGRESS.yaml`. Minor but will confuse anyone reading the constitution literally.

### m2. YAML indentation inconsistency

- `SKILLS.yaml` and `PROGRESS.yaml` use 2-space indent for keys
- `CHALLENGE_REGISTRY.yaml` uses 3-space indent for challenge IDs, 4-space for fields

This is why the `complete-challenge` regex uses `^  %s:` (2 spaces) for PROGRESS.yaml but `^    status:` (4 spaces) for CHALLENGE_REGISTRY.yaml. Fragile.

### m3. No `.gitignore` for generated files

`SYSTEM_AUDIT.md` should be in `.gitignore` if auto-generated. Also no `.golangci.yml` config.

---

## 5. CURRICULUM AUDIT: SKILL COVERAGE VS MARKET

### 5.1 Market requirements (Google L5, Meta E5, Uber Sr, ezCater Sr — 2026 data)

| Domain | Market demand | gokata coverage |
|--------|--------------|-----------------|
| Concurrency (errgroup, worker pool, context) | ✅ Essential | ✅ 11 skills, 6 challenges |
| Runtime internals (GMP, GC, escape) | ✅ Senior differentiator | ✅ 6 skills, 2 challenges |
| Distributed systems (Raft, CAP, MVCC) | ✅ Senior/Staff | ✅ 8 skills, 2 challenges |
| Architecture (DDD, CQRS, ES, SAGA) | ✅ Senior | ✅ 8 skills, 3 challenges |
| Production (graceful shutdown, observability) | ✅ Essential | ✅ 5 skills, spread across challenges |
| Testing (table-driven, fuzz, race, pprof) | ✅ Essential | ✅ 7 skills, 2 challenges |
| Communication (gRPC, Kafka, NATS) | ✅ Senior | ✅ 5 skills, 2 challenges |
| **System design (whiteboard, trade-offs, scaling)** | 🔴 **#1 reason seniors fail** | ❌ **0 challenges** |
| **Databases (SQL, transactions, optimization)** | ✅ ezCater explicit requirement | ❌ **0 challenges** |
| Infrastructure (K8s, Docker, Terraform, CI/CD) | ✅ Varies by role | 🟡 4 skills, 2 challenges (pending) |

### 5.2 Strengths vs competitors

- **Runtime internals coverage** is unique — no major interview prep resource covers GMP, GC, escape analysis at this depth
- **Raft + SAGA/Outbox/Inbox/DLQ** chain is rare and valuable for event-driven architecture roles
- **Spaced repetition** as a learning mechanic is absent from all competitor systems (LeetCode, Exercism, Codewars)

### 5.3 Critical gaps

1. **System Design: zero katas.** The GOAL.md Phase 3 identifies it (ezCater ranking service) but CHALLENGE_REGISTRY.yaml has no entries. This is the single highest-weighted interview category for L5+.
2. **Databases: zero katas.** ezCater job posts explicitly require "deep knowledge of SQL / relational databases, Postgres." Student SURVEY.md self-identifies this as a weakness.
3. **No culminating integration kata.** After 15 isolated challenges, there's no "build the ranking service" finale that exercises all skills together.

---

## 6. PEDAGOGICAL DESIGN AUDIT

### 6.1 Prerequisite chain integrity

```
worker_pool / errgroup (01) ──────┬── singleflight (02)
                                  ├── rate_limiter+cb (08)
                                  ├── fan_out+pipeline (11)
                                  ├── grpc (03) ──┬── event_sourcing (05) ── saga (06)
                                  │              └── full_system (07)
                                  ├── gmp_scheduler (09) ── gc+memory (10)
                                  └── atomic+sync_map (13)

ddd (04) ──┬── event_sourcing (05) ── saga (06) ── full_system (07)
           └── repository_pattern (via DDD)

Everything → raft (14), ci_cd (15)
```

**All prerequisite chains are valid.** No circular dependencies. Skills with pending prereqs are correctly blocked (`singleflight` requires `errgroup` + `context`, currently blocked because errgroup=0).

### 6.2 next-challenge algorithm correctness

The recommender scores each eligible challenge:
- NEW skill (0 pts) → +10 per skill
- REINFORCE due → +5 per skill
- Returns highest-scoring eligible challenge

**Current output:**
```
→ 06-saga-choreography — Saga Choreography with Compensation
  Skills:
    + saga (4) — NEW
    + transactional_outbox (2) — NEW
    + inbox_pattern (2) — NEW
    + dead_letter_queue (2) — NEW
    + kafka (1) — REINFORCE (cur: 2)
```

Score: 10×4 + 5×1 = 45. This is correct given current progress. However, this exposes a logic flaw: recommending SAGA (M, 10h) before the student has practiced ANY concurrency pattern is pedagogically questionable. The algorithm optimizes for "coverage of zero-point skills" but ignores learning progression common sense.

### 6.3 Spaced repetition intervals

| Points | Reinforcement interval | Implemented? |
|--------|----------------------|--------------|
| 0-3 | 1 challenge | ✅ |
| 4-6 | 2-3 challenges | ✅ |
| 7-8 | 4-6 challenges | ✅ |
| 9-10 | 7-10 challenges | ✅ |

Correctly implemented in `cmd/next-challenge/main.go:83-101`. No point decay (see M5).

---

## 7. RECOMMENDATIONS (PRIORITY ORDER)

### P0 — Must fix before use

1. **Rewrite `complete-challenge` using YAML struct marshaling, not regex.** This is the single point of failure for the entire data model. Use `gopkg.in/yaml.v3` to load, modify, and re-serialize both PROGRESS.yaml and CHALLENGE_REGISTRY.yaml. The regex approach is fundamentally unsalvageable.

2. **Add tests for all three CLI tools.** Use `testdata/` fixtures with known-good and known-bad YAML files. Verify:
   - `validate` passes on clean data, fails on corrupt data
   - `complete-challenge` correctly modifies only the target challenge
   - `complete-challenge` is idempotent (second run is no-op)
   - `next-challenge` respects blocked prereqs

### P1 — Should fix before Phase 2

3. **Add CHALLENGE_LOG.yaml** machine-readable format. Keep `.md` as generated view.

4. **Populate `challenges: []` arrays** in PROGRESS.yaml on completion.

5. **Add a `Makefile`** for centralized command dispatch.

6. **Fix AGENTS.md typo** `PROGRAMS.yaml` → `PROGRESS.yaml`.

### P2 — Strategic additions

7. **Design 1-2 system design challenges** for the ezCater target (ranking/recommendation). Format: `design.md` with context, constraints, trade-off questions, and evaluation rubric — no code.

8. **Design 1 database challenge** (Postgres transactions, query optimization, migration patterns).

---

## 8. VERDICT

The **tutor scaffolding is structurally sound and well-conceived.** The skill taxonomy, prerequisite chaining, and spaced repetition model are superior to mainstream alternatives. The curriculum covers the differentiating topics (runtime internals, Raft, SAGA) that most resources ignore.

However, the **core completion tool is broken in a way that corrupts the data it manages.** This is not a minor bug — it means every use of the system damages the record it depends on. Additionally, zero test coverage means regressions are invisible.

**The system is not reliable enough to run unattended.** It requires an agent to validate every output manually, which defeats the purpose of automation. With ~8h of focused repair (rewrite `complete-challenge` + tests + CHALLENGE_LOG.yaml), it becomes a solid foundation.

**The curriculum itself is worth preserving.** The skill selection and challenge design are better than 90% of Go interview prep resources. The gaps are in system design and databases, which are the exact filters for ezCater senior roles.

---

*Data sources: interviewkickstart.com, 1point3acres.com, sharpskill.dev, codeforgeek.com, nazarboyko.com, glassdoor.com (ezCater), jobright.ai.*
