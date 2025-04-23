//go:build !js

package env

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Getenv(name string) string {
	file, err := os.Open("wrangler.toml")
	if err != nil {
		log.Fatalf("Failed to open wrangler.toml: %v", err)
	}
	defer file.Close()

	// Read properties file instead of CSV
	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue // Skip empty lines and comments
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			properties[key] = value
		}
	}

	return properties[name]
}

// func GetBinding(name string) js.Value {
// 	return cfruntimecontext.MustGetRuntimeContextEnv().Get(name)
// }
