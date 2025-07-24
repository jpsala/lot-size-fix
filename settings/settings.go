package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	AutoClose bool `json:"autoClose"`
}

func LoadSettings() (Settings, error) {
	var settings Settings
	exePath, err := os.Executable()
	if err != nil {
		return settings, err
	}

	dir := filepath.Dir(exePath)
	settingsPath := filepath.Join(dir, "settings.json")
	file, err := os.Open(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultSettings(settingsPath)
		}
		return settings, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	return settings, err
}

func createDefaultSettings(path string) (Settings, error) {
	settings := Settings{
		AutoClose: true,
	}

	file, err := os.Create(path)
	if err != nil {
		return settings, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(settings)
	return settings, err
}
