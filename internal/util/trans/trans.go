package trans

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed fa.json
var fa string

var translations map[string]string

func Get(key string) string {
	if err := json.Unmarshal([]byte(fa), &translations); err != nil {
		log.Fatalf("failed to parse translation file: %w", err)
	}

	return translations[key]
}
