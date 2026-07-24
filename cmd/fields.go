package cmd

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"
)

// resolveFields merges field values for record create/update from three
// optional sources, in increasing precedence: --file, --data, --field.
// Later sources overwrite earlier ones on key conflicts.
func resolveFields(dataJSON, filePath string, fieldFlags []string) (map[string]any, error) {
	fields := map[string]any{}

	if filePath != "" {
		b, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", filePath, err)
		}
		if err := json.Unmarshal(b, &fields); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", filePath, err)
		}
	}

	if dataJSON != "" {
		var data map[string]any
		if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
			return nil, fmt.Errorf("parsing --data: %w", err)
		}
		maps.Copy(fields, data)
	}

	for _, entry := range fieldFlags {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			return nil, fmt.Errorf("invalid --field %q (want key=value)", entry)
		}
		fields[key] = value
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields provided: use --field, --data, or --file")
	}

	return fields, nil
}
