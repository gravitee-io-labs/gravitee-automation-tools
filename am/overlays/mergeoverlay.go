package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type overlay struct {
	Overlay string      `yaml:"overlay"`
	Info    yaml.Node   `yaml:"info"`
	Actions []yaml.Node `yaml:"actions"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: mergeoverlay in1.yaml [in2.yaml ...] out.yaml")
		os.Exit(1)
	}
	outPath := os.Args[len(os.Args)-1]
	var merged overlay
	for i, path := range os.Args[1 : len(os.Args)-1] {
		raw, err := os.ReadFile(path)
		if err != nil {
			fatal(err)
		}
		var next overlay
		if err := yaml.Unmarshal(raw, &next); err != nil {
			fatal(fmt.Errorf("%s: %w", path, err))
		}
		if i == 0 {
			merged.Overlay = next.Overlay
			merged.Info = next.Info
		}
		merged.Actions = append(merged.Actions, next.Actions...)
	}
	raw, err := yaml.Marshal(&merged)
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(outPath, raw, 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
