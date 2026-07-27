# gokata — AGENTS.md

## Project

Go interview-prep katas. The student implements production-quality patterns from tests + hints alone.

## Structure

```
gokata/
  GOAL.md          — curriculum & student profile (Spanish)
  SURVEY.md        — knowledge tree
  kata/N-name/
    hints.md       — context, signatures, hints, alternatives
    *_test.go      — tests only (no impl files)
```

- **No implementation files exist** — student writes `.go` files in each kata dir.
- Types are defined by the student; tests reference them from the same `package kata`.
- Tests **will fail** until the student implements the required types and functions.

## Working with katas

- Each kata is self-contained under `kata/N-name/`.
- Run a single kata: `go test ./kata/N-name/ -v`
- Run a specific test: `go test ./kata/N-name/ -run TestRanker_ReturnsExactN -v`
- Benchmark: `go test ./kata/N-name/ -bench=. -benchmem`
- Lint/vet: `go vet ./kata/N-name/`
- Check compilation without running tests: `go build ./kata/N-name/`

## Commands

```bash
go test ./kata/01-workerpool/ -v          # single kata, verbose
go test ./kata/01-workerpool/ -run TestX  # single test
go test ./...                              # all katas
```

## Dependencies

- Go 1.22
- `golang.org/x/sync` (errgroup, semaphore.Weighted)

## Notable

- Tests in `ranker_test.go` use haversine formula for distance.
- The heap approach is recommended for top-N sorting over full sort.
- Kata 01 expects types: `Restaurant`, `ScoredRestaurant`, `Config`, `Ranker`, and function `New`.
