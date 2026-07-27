package types

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	Cluster       string   `yaml:"cluster"`
	Description   string   `yaml:"description"`
	MaxPoints     int      `yaml:"max_points"`
	Prerequisites []string `yaml:"prerequisites"`
}

type SkillsFile struct {
	Version int               `yaml:"version"`
	Skills  map[string]*Skill `yaml:"skills"`
}

type ProgressEntry struct {
	Points       int      `yaml:"points"`
	LastReviewed *string  `yaml:"last_reviewed"`
	Challenges   []string `yaml:"challenges"`
}

type ProgressFile struct {
	Progress map[string]*ProgressEntry `yaml:"progress"`
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

type RegistryFile struct {
	Challenges map[string]*ChallengeEntry `yaml:"challenges"`
}

func FindRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "SKILLS.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Fprintln(os.Stderr, "fatal: SKILLS.yaml not found")
			os.Exit(1)
		}
		dir = parent
	}
}

func LoadProgress(path string) map[string]*ProgressEntry {
	var pf ProgressFile
	MustUnmarshal(path, &pf)
	return pf.Progress
}

func MustRead(path string) []byte {
	d, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: reading %s: %v\n", path, err)
		os.Exit(1)
	}
	return d
}

func MustWrite(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: writing %s: %v\n", path, err)
		os.Exit(1)
	}
}

func MustUnmarshal(path string, v interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: reading %s: %v\n", path, err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(data, v); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: parsing %s: %v\n", path, err)
		os.Exit(1)
	}
}
