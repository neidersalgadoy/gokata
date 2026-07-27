package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/neidersalgadoy/gokata/internal/types"
)

func main() {
	root := types.FindRoot()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: complete-challenge <id> [horas]")
		os.Exit(1)
	}
	id := os.Args[1]

	reg := types.RegistryFile{}
	types.MustUnmarshal(filepath.Join(root, "CHALLENGE_REGISTRY.yaml"), &reg)
	ch, ok := reg.Challenges[id]
	if !ok {
		fmt.Fprintf(os.Stderr, "fatal: challenge %q not found\n", id)
		os.Exit(1)
	}

	today := time.Now().Format("2006-01-02")

	// Update PROGRESS.yaml (text-based, preserves comments)
	progPath := filepath.Join(root, "PROGRESS.yaml")
	data := string(types.MustRead(progPath))
	for s, pts := range ch.Skills {
		re := regexp.MustCompile(
			fmt.Sprintf(`(?m)^  %s:\n    points: (\d+)\n    last_reviewed: .*`, regexp.QuoteMeta(s)))
		m := re.FindStringSubmatch(data)
		if m == nil {
			fmt.Fprintf(os.Stderr, "warning: skill %q not found\n", s)
			continue
		}
		cur, _ := strconv.Atoi(m[1])
		repl := fmt.Sprintf("  %s:\n    points: %d\n    last_reviewed: \"%s\"", s, cur+pts, today)
		data = re.ReplaceAllString(data, repl)
	}
	types.MustWrite(progPath, []byte(data))

	// Update CHALLENGE_REGISTRY.yaml
	regPath := filepath.Join(root, "CHALLENGE_REGISTRY.yaml")
	regData := string(types.MustRead(regPath))
	pat := fmt.Sprintf("(?m)^  %s:\n(?:.|\\n)*?^    status: pending", regexp.QuoteMeta(id))
	statusRe := regexp.MustCompile(pat)
	regData = statusRe.ReplaceAllStringFunc(regData, func(b string) string {
		return regexp.MustCompile(`status: pending`).ReplaceAllString(b, "status: completed")
	})
	types.MustWrite(regPath, []byte(regData))

	// Append to CHALLENGE_LOG.md
	logPath := filepath.Join(root, "CHALLENGE_LOG.md")
	log := string(types.MustRead(logPath))
	numRe := regexp.MustCompile(`(?m)^\| (\d+) \|`)
	nums := numRe.FindAllStringSubmatch(log, -1)
	next := 1
	for _, n := range nums {
		if v, _ := strconv.Atoi(n[1]); v >= next {
			next = v + 1
		}
	}
	skills := ""
	total := 0
	for s, p := range ch.Skills {
		if skills != "" {
			skills += ", "
		}
		skills += s
		total += p
	}
	line := fmt.Sprintf("| %d | %s | %s | %s | %d | %s |\n", next, id, ch.Talla, skills, total, today)
	types.MustWrite(logPath, []byte(log+line))

	fmt.Printf("✓ %q completed. %d skills updated.\n", id, len(ch.Skills))
}
