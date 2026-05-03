package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ArtifactPath(artifactsDir, name string, version int) string {
	return filepath.Join(artifactsDir, fmt.Sprintf("%s_v%d.json", name, version))
}

func Write(artifactsDir, name string, version int, schema SchemaName, data []byte) error {
	if err := Validate(schema, data); err != nil {
		return fmt.Errorf("artifact %s v%d failed validation: %w", name, version, err)
	}

	path := ArtifactPath(artifactsDir, name, version)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func Read(artifactsDir, name string, version int) ([]byte, error) {
	path := ArtifactPath(artifactsDir, name, version)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read artifact %s v%d: %w", name, version, err)
	}
	return data, nil
}

func ReadInto(artifactsDir, name string, version int, dest any) error {
	data, err := Read(artifactsDir, name, version)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func LatestVersion(artifactsDir, name string) (int, error) {
	version := 0
	for {
		path := ArtifactPath(artifactsDir, name, version+1)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		} else if err != nil {
			return 0, err
		}
		version++
	}
	if version == 0 {
		return 0, fmt.Errorf("no artifact found for %s", name)
	}
	return version, nil
}
