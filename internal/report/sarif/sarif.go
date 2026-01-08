package sarif

import (
	"encoding/json"
	"os"
)

func Write(data any, path string) error {
	b, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(path, b, 0644)
}
