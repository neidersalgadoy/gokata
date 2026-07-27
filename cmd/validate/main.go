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
	Version int               `yaml:"version"`
	Skills  map[string]*Skill `yaml:"skills"`
}

type Skill struct {
	Cluster       string   `yaml:"cluster"`
	Description   string   `yaml:"description"`
	MaxPoints     int      `yaml:"max_points"`
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

type RegistryFile struct {
	Challenges map[string]*ChallengeEntry `yaml:"challenges"`
}

type ChallengeEntry struct {
	Title       string         `yaml:"title"`
	Talla       string         `yaml:"talla"`
	Horas       int            `yaml:"horas"`
	Cluster     string         `yaml:"cluster"`
	Skills      map[string]int `yaml:"skills"`
	Prereqs     []string       `yaml:"prereqs"`
	Status      string         `yaml:"status"`
}

func main() {
	exit := 0

	repoRoot := findRepoRoot()
	skillsFile := filepath.Join(repoRoot, "SKILLS.yaml")
	progressFile := filepath.Join(repoRoot, "PROGRESS.yaml")
	registryFile := filepath.Join(repoRoot, "CHALLENGE_REGISTRY.yaml")
	kataDir := filepath.Join(repoRoot, "kata")

	skills := loadSkills(skillsFile)
	progress := loadProgress(progressFile)
	registry := loadRegistry(registryFile)

	errors := validateProgress(skills, progress, kataDir)
	warnings := checkProgressPrerequisites(skills, progress)

	if registry != nil {
		rErrors, rWarns := validateRegistry(skills, kataDir, registry)
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

func loadRegistry(path string) *RegistryFile {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var rf RegistryFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: parsing %s: %v\n", path, err)
		os.Exit(1)
	}
	return &rf
}

func validateProgress(skills map[string]*Skill, progress map[string]*ProgressEntry, kataDir string) []string {
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

func checkProgressPrerequisites(skills map[string]*Skill, progress map[string]*ProgressEntry) []string {
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

func validateRegistry(skills map[string]*Skill, kataDir string, rf *RegistryFile) ([]string, []string) {
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

func allChallengeRefs(progress map[string]*ProgressEntry) []string {
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
