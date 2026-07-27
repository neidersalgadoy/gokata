package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

type SkillsFile struct {
	Version int                `yaml:"version"`
	Skills  map[string]*Skill  `yaml:"skills"`
}

type Skill struct {
	Cluster      string   `yaml:"cluster"`
	Description  string   `yaml:"description"`
	MaxPoints    int      `yaml:"max_points"`
	Prerequisites []string `yaml:"prerequisites"`
}

type ProgressFile struct {
	Progress map[string]*ProgressEntry `yaml:"progress"`
}

type ProgressEntry struct {
	Points       int      `yaml:"points"`
	LastReviewed *string  `yaml:"last_reviewed"`
	Challenges   []string `yaml:"challenges"`
}

func main() {
	exit := 0

	repoRoot := findRepoRoot()
	skillsFile := filepath.Join(repoRoot, "SKILLS.yaml")
	progressFile := filepath.Join(repoRoot, "PROGRESS.yaml")
	kataDir := filepath.Join(repoRoot, "kata")

	skills := loadSkills(skillsFile)
	progress := loadProgress(progressFile)

	invalidSkills := findInvalidSkills(skills, progress)
	for _, name := range invalidSkills {
		fmt.Fprintf(os.Stderr, "ERROR: skill %q in PROGRESS.yaml not found in SKILLS.yaml\n", name)
		exit = 1
	}

	var errors []string
	var warnings []string

	for name, entry := range progress {
		skill, ok := skills[name]
		if !ok {
			continue
		}
		if entry.Points > skill.MaxPoints {
			errors = append(errors, fmt.Sprintf("skill %q: points (%d) > max_points (%d)", name, entry.Points, skill.MaxPoints))
		}
		if entry.Points < 0 {
			errors = append(errors, fmt.Sprintf("skill %q: points (%d) < 0", name, entry.Points))
		}
		if entry.LastReviewed != nil {
			t, err := time.Parse("2006-01-02", *entry.LastReviewed)
			if err != nil {
				errors = append(errors, fmt.Sprintf("skill %q: invalid last_reviewed date %q: %v", name, *entry.LastReviewed, err))
			} else if t.After(time.Now().Add(24 * time.Hour)) {
				warnings = append(warnings, fmt.Sprintf("skill %q: last_reviewed (%s) is in the future", name, *entry.LastReviewed))
			}
		}
	}

	for _, ch := range findAllReferencedChallenges(progress) {
		dir := filepath.Join(kataDir, ch)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			warnings = append(warnings, fmt.Sprintf("challenge %q referenced but not found in kata/", ch))
		}
	}

	prereqErrs := checkPrerequisites(skills, progress)
	warnings = append(warnings, prereqErrs...)

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

func findRepoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "SKILLS.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Fprintln(os.Stderr, "fatal: SKILLS.yaml not found in any parent directory")
			os.Exit(1)
		}
		dir = parent
	}
}

func loadSkills(path string) map[string]*Skill {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: reading %s: %v\n", path, err)
		os.Exit(1)
	}
	var sf SkillsFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: parsing %s: %v\n", path, err)
		os.Exit(1)
	}
	return sf.Skills
}

func loadProgress(path string) map[string]*ProgressEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: reading %s: %v\n", path, err)
		os.Exit(1)
	}
	var pf ProgressFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: parsing %s: %v\n", path, err)
		os.Exit(1)
	}
	return pf.Progress
}

func findInvalidSkills(skills map[string]*Skill, progress map[string]*ProgressEntry) []string {
	var invalid []string
	for name := range progress {
		if _, ok := skills[name]; !ok {
			invalid = append(invalid, name)
		}
	}
	sort.Strings(invalid)
	return invalid
}

func findAllReferencedChallenges(progress map[string]*ProgressEntry) []string {
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

func checkPrerequisites(skills map[string]*Skill, progress map[string]*ProgressEntry) []string {
	var warnings []string
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
				warnings = append(warnings, fmt.Sprintf("skill %q (points=%d) has prerequisite %q but no progress entry found", name, entry.Points, prereq))
				continue
			}
			if prEntry.Points == 0 {
				warnings = append(warnings, fmt.Sprintf("skill %q (points=%d) has prerequisite %q with 0 points", name, entry.Points, prereq))
			}
		}
	}
	return warnings
}
