package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/neidersalgadoy/gokata/internal/types"
)

type candidate struct {
	ID    string
	Entry *types.ChallengeEntry
	Score int
}

func main() {
	root := types.FindRoot()
	progress := types.LoadProgress(filepath.Join(root, "PROGRESS.yaml"))
	reg := types.RegistryFile{}
	types.MustUnmarshal(filepath.Join(root, "CHALLENGE_REGISTRY.yaml"), &reg)

	var cs []candidate
	for id, ch := range reg.Challenges {
		if ch.Status != "pending" {
			continue
		}
		ok := true
		for _, p := range ch.Prereqs {
			if e, exists := progress[p]; !exists || e.Points == 0 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		score := 0
		for s := range ch.Skills {
			if e, exists := progress[s]; !exists || e.Points == 0 {
				score += 10
			} else if needsReinforcement(e) {
				score += 5
			}
		}
		cs = append(cs, candidate{id, ch, score})
	}

	if len(cs) == 0 {
		fmt.Println("No eligible challenges.")
		return
	}

	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Score != cs[j].Score {
			return cs[i].Score > cs[j].Score
		}
		return cs[i].ID < cs[j].ID
	})

	top := cs[0]
	fmt.Printf("→ %s — %s\n", top.ID, top.Entry.Title)
	fmt.Printf("  Talla: %s | Horas: %d\n", top.Entry.Talla, top.Entry.Horas)
	fmt.Println("  Skills:")
	zero, reinforce := []string{}, []string{}
	for s, pts := range top.Entry.Skills {
		e := progress[s]
		if e == nil || e.Points == 0 {
			zero = append(zero, fmt.Sprintf("    + %s (%d) — NEW", s, pts))
		} else if needsReinforcement(e) {
			reinforce = append(reinforce, fmt.Sprintf("    + %s (%d) — REINFORCE (cur: %d)", s, pts, e.Points))
		}
	}
	for _, l := range zero {
		fmt.Println(l)
	}
	for _, l := range reinforce {
		fmt.Println(l)
	}
}

func needsReinforcement(e *types.ProgressEntry) bool {
	if e.LastReviewed == nil {
		return e.Points > 0
	}
	t, err := time.Parse("2006-01-02", *e.LastReviewed)
	if err != nil {
		return false
	}
	d := int(time.Since(t).Hours() / 24)
	switch {
	case e.Points <= 3:
		return d >= 1
	case e.Points <= 6:
		return d >= 3
	case e.Points <= 8:
		return d >= 5
	default:
		return d >= 8
	}
}


