package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	nameRegex  = regexp.MustCompile(`^[a-z][a-z0-9-]+$`)
	ownerRegex = regexp.MustCompile(`^group:account/.+`)
	refRegex   = regexp.MustCompile(`^(system|resource|component|api):account/.+`)
	validKinds = map[string]bool{"System": true, "Component": true, "Resource": true, "API": true}
)

type Entity struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Identifier string                 `yaml:"identifier"`
	Name       string                 `yaml:"name"`
	Owner      string                 `yaml:"owner"`
	Spec       map[string]interface{} `yaml:"spec"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: catalog-validator <path-to-catalog-info.yaml>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read file: %v\n", err)
		os.Exit(1)
	}

	var errors []string
	seenIdentifiers := map[string]bool{}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	docIndex := 0
	for {
		var entity Entity
		err := decoder.Decode(&entity)
		if err != nil {
			break
		}
		if entity.APIVersion == "" && entity.Kind == "" {
			docIndex++
			continue
		}
		docIndex++
		prefix := fmt.Sprintf("doc[%d]", docIndex)

		if entity.APIVersion != "harness.io/v1" {
			errors = append(errors, fmt.Sprintf("%s: apiVersion must be harness.io/v1, got %q", prefix, entity.APIVersion))
		}
		if !validKinds[entity.Kind] {
			errors = append(errors, fmt.Sprintf("%s: kind must be System|Component|Resource|API, got %q", prefix, entity.Kind))
		}
		if entity.Identifier == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level identifier", prefix))
		} else if seenIdentifiers[entity.Identifier] {
			errors = append(errors, fmt.Sprintf("%s: duplicate identifier %q", prefix, entity.Identifier))
		} else {
			seenIdentifiers[entity.Identifier] = true
		}
		if entity.Name == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level name", prefix))
		} else if !nameRegex.MatchString(entity.Name) {
			errors = append(errors, fmt.Sprintf("%s: name %q must match ^[a-z][a-z0-9-]+$", prefix, entity.Name))
		}
		if entity.Owner == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level owner", prefix))
		} else if !ownerRegex.MatchString(entity.Owner) {
			errors = append(errors, fmt.Sprintf("%s: owner %q must be group:account/<id>", prefix, entity.Owner))
		}
		if entity.Spec != nil {
			if sys, ok := entity.Spec["system"]; ok {
				refs, ok := sys.([]interface{})
				if !ok {
					errors = append(errors, fmt.Sprintf("%s: spec.system must be a list", prefix))
				} else {
					for _, r := range refs {
						s, _ := r.(string)
						if !refRegex.MatchString(s) {
							errors = append(errors, fmt.Sprintf("%s: spec.system entry %q must be system:account/<id>", prefix, s))
						}
					}
				}
			}
			if deps, ok := entity.Spec["dependsOn"]; ok {
				refs, ok := deps.([]interface{})
				if !ok {
					errors = append(errors, fmt.Sprintf("%s: spec.dependsOn must be a list", prefix))
				} else {
					for _, r := range refs {
						s, _ := r.(string)
						if !refRegex.MatchString(s) {
							errors = append(errors, fmt.Sprintf("%s: spec.dependsOn entry %q must be resource:account/<id>", prefix, s))
						}
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
	fmt.Println("catalog-info.yaml is valid")
}
