package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// loadYAML reads a YAML file and returns a map[string]interface{}
func loadYAML(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read %s: %v", filePath, err)
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, errors.New("failed to parse %s: %v", filePath, err)
	}

	return content, nil
}

// flattenKeys recursively flattens a YAML structure into a set of keys
func flattenKeys(data map[string]interface{}, prefix string) map[string]struct{} {
	keys := make(map[string]struct{})
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		keys[fullKey] = struct{}{}

		if nested, ok := value.(map[string]interface{}); ok {
			nestedKeys := flattenKeys(nested, fullKey)
			for k := range nestedKeys {
				keys[k] = struct{}{}
			}
		}
	}
	return keys
}

func main() {
	prodConfigPath := "config/prod/prod.yaml"
	stagConfigDir := "config/stag"

	// Load the production YAML file
	prodConfig, err := loadYAML(prodConfigPath)
	if err != nil {
		log.Fatalf("❌ Error loading production config: %v", err)
	}

	prodKeys := flattenKeys(prodConfig, "")

	// Read and validate each file in stag directory
	missingKeysReport := make(map[string][]string)

	err = filepath.Walk(stagConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only YAML files
		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			stagConfig, err := loadYAML(path)
			if err != nil {
				log.Printf("⚠️ Warning: Skipping %s due to error: %v\n", path, err)
				return nil
			}

			stagKeys := flattenKeys(stagConfig, "")

			// Identify missing keys
			var missingKeys []string
			for key := range prodKeys {
				if _, exists := stagKeys[key]; !exists {
					missingKeys = append(missingKeys, key)
				}
			}

			if len(missingKeys) > 0 {
				missingKeysReport[path] = missingKeys
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("❌ Error reading staging config files: %v", err)
	}

	// If missing keys are found, print and fail
	if len(missingKeysReport) > 0 {
		fmt.Println("❌ Missing keys detected in staging config files:\n")
		for file, keys := range missingKeysReport {
			fmt.Printf("🔹 %s is missing:\n", file)
			for _, key := range keys {
				fmt.Printf("   - %s\n", key)
			}
		}
		os.Exit(1) // Fail the GitHub Action
	} else {
		fmt.Println("✅ All staging config files contain the required keys.")
	}
}
