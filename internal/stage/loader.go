package stage

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed data/*.json
var stageData embed.FS

// LoadAllStages loads all stages from the embedded JSON files.
func LoadAllStages() ([]Stage, error) {
	entries, err := stageData.ReadDir("data")
	if err != nil {
		return nil, fmt.Errorf("failed to read stage data directory: %w", err)
	}

	var allStages []Stage

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := stageData.ReadFile("data/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read stage file %s: %w", entry.Name(), err)
		}

		var sf StageFile
		if err := json.Unmarshal(data, &sf); err != nil {
			return nil, fmt.Errorf("failed to parse stage file %s: %w", entry.Name(), err)
		}

		allStages = append(allStages, sf.Stages...)
	}

	sort.Slice(allStages, func(i, j int) bool {
		return allStages[i].ID < allStages[j].ID
	})

	return allStages, nil
}
