package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Preset struct {
	Name string
	Src  string
	Dst  string
}

func configPath() string {
	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".go_photo")
}

func LoadPresets() ([]Preset, error) {
	path := configPath()

	file, err := os.Open(path)

	if os.IsNotExist(err) {
		return []Preset{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var presets []Preset

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		preset := parseLine(line)

		if preset.Name != "" {
			presets = append(presets, preset)
		}
	}

	return presets, scanner.Err()
}

func parseLine(line string) Preset {
	parts := strings.Split(line, "|")

	if len(parts) < 3 {
		return Preset{}
	}

	preset := Preset{Name: parts[0]}

	for _, p := range parts[1:] {
		kv := strings.Split(p, "::")

		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "src":
			preset.Src = kv[1]
		case "dst":
			preset.Dst = kv[1]
		}
	}

	return preset
}

func SavePresets(presets []Preset) error {
	path := configPath()

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	for _, p := range presets {
		line := p.Name + "|src::" + p.Src + "|dst::" + p.Dst + "\n"

		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

func AddPreset(preset Preset) error {
	presets, err := LoadPresets()

	if err != nil {
		return err
	}

	presets = append(presets, preset)

	return SavePresets(presets)
}

func DeletePreset(name string) error {
	presets, err := LoadPresets()

	if err != nil {
		return err
	}

	var filtered []Preset

	for _, p := range presets {
		if p.Name != name {
			filtered = append(filtered, p)
		}
	}

	return SavePresets(filtered)
}
