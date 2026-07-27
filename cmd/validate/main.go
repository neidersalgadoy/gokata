package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/neidersalgadoy/gokata/internal/types"
)

func main() {
	exit := 0

	repoRoot := types.FindRoot()
	skillsFile := filepath.Join(repoRoot, "SKILLS.yaml")
	progressFile := filepath.Join(repoRoot, "PROGRESS.yaml")
	registryFile := filepath.Join(repoRoot, "CHALLENGE_REGISTRY.yaml")
	kataDir := filepath.Join(repoRoot, "kata")

	var sf types.SkillsFile
	types.MustUnmarshal(skillsFile, &sf)
	skills := sf.Skills

	progress := types.LoadProgress(progressFile)

	var rf types.RegistryFile
	if err := func() error {
		// registry is optional
		if _, err := os.Stat(registryFile); os.IsNotExist(err) {
			return err
		}
		types.MustUnmarshal(registryFile, &rf)
		return nil
	}(); err != nil {
		rf = types.RegistryFile{}
	}

	errors := validateProgress(skills, progress, kataDir)
	warnings := checkProgressPrerequisites(skills, progress)

	if len(rf.Challenges) > 0 {
		rErrors, rWarns := validateRegistry(skills, kataDir, &rf)
		errors = append(errors, rErrors...)
		warnings = append(warnings, rWarns...)
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "=== ERRORS ===\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  %s\n", e)
		}
		exit = 1
	}

	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "=== WARNINGS ===\n")
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "  %s\n", w)
		}
	}

	if exit == 0 {
		fmt.Println("OK — all checks passed")
	}
	os.Exit(exit)
}

func validateProgress(skills map[string]*types.Skill, progress map[string]*types.ProgressEntry, kataDir string) []string {
	var errors []string

	for name := range progress {
		if _, ok := skills[name]; !ok {
			errors = append(errors, fmt.Sprintf("PROGRESS: skill %q not found in SKILLS.yaml", name))
		}
	}

	for name, entry := range progress {
		skill, ok := skills[name]
		if !ok {
			continue
		}
		if entry.Points > skill.MaxPoints {
			errors = append(errors, fmt.Sprintf("PROGRESS: skill %q points (%d) > max_points (%d)", name, entry.Points, skill.MaxPoints))
		}
		if entry.Points < 0 {
			errors = append(errors, fmt.Sprintf("PROGRESS: skill %q points (%d) < 0", name, entry.Points))
		}
		if entry.LastReviewed != nil {
			t, err := time.Parse("2006-01-02", *entry.LastReviewed)
			if err != nil {
				errors = append(errors, fmt.Sprintf("PROGRESS: skill %q invalid last_reviewed date %q: %v", name, *entry.LastReviewed, err))
			} else if t.After(time.Now().Add(24*time.Hour)) {
				errors = append(errors, fmt.Sprintf("PROGRESS: skill %q last_reviewed (%s) is in the future", name, *entry.LastReviewed))
			}
		}
	}

	for _, ch := range allChallengeRefs(progress) {
		info, err := os.Stat(filepath.Join(kataDir, ch))
		if err != nil || !info.IsDir() {
			errors = append(errors, fmt.Sprintf("PROGRESS: challenge %q referenced but not found in kata/", ch))
		}
	}

	return errors
}

func checkProgressPrerequisites(skills map[string]*types.Skill, progress map[string]*types.ProgressEntry) []string {
	var warns []string
	for name, entry := range progress {
		if entry.Points == 0 {
			continue
		}
		skill, ok := skills[name]
		if !ok {
			continue
		}
		for _, prereq := range skill.Prerequisites {
			prEntry, ok := progress[prereq]
			if !ok {
				warns = append(warns, fmt.Sprintf("PROGRESS: skill %q (pts=%d) has prerequisite %q but no progress entry", name, entry.Points, prereq))
				continue
			}
			if prEntry.Points == 0 {
				warns = append(warns, fmt.Sprintf("PROGRESS: skill %q (pts=%d) has prerequisite %q with 0 points", name, entry.Points, prereq))
			}
		}
	}
	return warns
}

func validateRegistry(skills map[string]*types.Skill, kataDir string, rf *types.RegistryFile) ([]string, []string) {
	var errors []string
	var warns []string

	validStatuses := map[string]bool{"pending": true, "active": true, "completed": true}

	for id, ch := range rf.Challenges {
		if _, ok := validStatuses[ch.Status]; !ok {
			errors = append(errors, fmt.Sprintf("REGISTRY: challenge %q has invalid status %q", id, ch.Status))
		}

		if ch.Talla == "" {
			errors = append(errors, fmt.Sprintf("REGISTRY: challenge %q missing talla", id))
		}

		validTallas := map[string]bool{"XS": true, "S": true, "M": true, "L": true, "XL": true}
		if !validTallas[ch.Talla] {
			warns = append(warns, fmt.Sprintf("REGISTRY: challenge %q has unusual talla %q", id, ch.Talla))
		}

		for sname := range ch.Skills {
			if _, ok := skills[sname]; !ok {
				errors = append(errors, fmt.Sprintf("REGISTRY: challenge %q references unknown skill %q", id, sname))
			}
		}

		for _, prereq := range ch.Prereqs {
			if _, ok := skills[prereq]; !ok {
				errors = append(errors, fmt.Sprintf("REGISTRY: challenge %q has unknown prerequisite %q (not a skill)", id, prereq))
			}
		}

		chDir := filepath.Join(kataDir, id)
		if ch.Status == "completed" || ch.Status == "active" {
			if _, err := os.Stat(chDir); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("REGISTRY: challenge %q is %s but kata/ dir does not exist", id, ch.Status))
			}
		}
	}

	return errors, warns
}

func allChallengeRefs(progress map[string]*types.ProgressEntry) []string {
	seen := make(map[string]bool)
	for _, entry := range progress {
		for _, ch := range entry.Challenges {
			seen[ch] = true
		}
	}
	var list []string
	for ch := range seen {
		list = append(list, ch)
	}
	sort.Strings(list)
	return list
}
