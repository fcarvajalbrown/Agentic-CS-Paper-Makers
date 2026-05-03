package artifacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func Validate(schemaName SchemaName, data []byte) error {
	raw, err := LoadSchema(schemaName)
	if err != nil {
		return fmt.Errorf("load schema %s: %w", schemaName, err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(string(schemaName), bytes.NewReader(raw)); err != nil {
		return fmt.Errorf("compile schema %s: %w", schemaName, err)
	}

	schema, err := compiler.Compile(string(schemaName))
	if err != nil {
		return fmt.Errorf("compile schema %s: %w", schemaName, err)
	}

	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := schema.Validate(v); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

func ValidateFile(schemaName SchemaName, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return Validate(schemaName, data)
}

func DumpError(errorsDir string, name SchemaName, data []byte, valErr error) error {
	if err := os.MkdirAll(errorsDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(errorsDir, string(name)+".invalid.json")
	if err := os.WriteFile(dest, data, 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "schema error: %v\ninvalid output written to: %s\n", valErr, dest)
	return nil
}
