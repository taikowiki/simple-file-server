package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadEnv() (map[string]string, error) {
	env := make(map[string]string)

	envPath := filepath.Join(MainDir(), ".env.json")
	envJson, err := os.ReadFile(envPath)
	if err != nil {
		return env, err
	}

	err = json.Unmarshal(envJson, &env)
	if err != nil {
		return env, err
	}

	return env, nil
}
