package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type CacheEntry struct {
	RequestID       string `json:"request_id"`
	ModelVersion    string `json:"model_version"`
	SystemFingerprint string `json:"system_fingerprint"`
	Timestamp       string `json:"timestamp"`
	TokensUsed      int    `json:"tokens_used"`
	Response        any    `json:"response"`
}

type Cache struct {
	dir string
}

func NewCache(cacheDir string) *Cache {
	return &Cache{dir: cacheDir}
}

func (c *Cache) key(prompt, model string, seed int) string {
	h := sha256.New()
	h.Write([]byte(prompt))
	h.Write([]byte(model))
	h.Write([]byte{byte(seed)})
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Cache) Get(prompt, model string, seed int) (*CacheEntry, bool) {
	path := filepath.Join(c.dir, c.key(prompt, model, seed)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	return &entry, true
}

func (c *Cache) Put(prompt, model string, seed int, entry *CacheEntry) error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(c.dir, c.key(prompt, model, seed)+".json")
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
