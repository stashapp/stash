package plugins

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/envy"
)

type cachedPlugin struct {
	Commands Commands `json:"commands"`
	CheckSum string   `json:"check_sum"`
}

type cachedPlugins map[string]cachedPlugin

var cachePath = func() string {
	home := "."
	if usr, err := user.Current(); err == nil {
		home = usr.HomeDir
	}
	return filepath.Join(home, ".buffalo", "plugin.cache")
}()

var cacheMoot sync.RWMutex

var cacheOn = envy.Get("BUFFALO_PLUGIN_CACHE", "on")

var cache = func() cachedPlugins {
	m := cachedPlugins{}
	if cacheOn != "on" {
		return m
	}
	f, err := os.Open(cachePath)
	if err != nil {
		return m
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		f.Close()
		os.Remove(f.Name())
	}
	return m
}()

func findInCache(path string) (cachedPlugin, bool) {
	cacheMoot.RLock()
	defer cacheMoot.RUnlock()
	cp, ok := cache[path]
	return cp, ok
}

func saveCache() error {
	if cacheOn != "on" {
		return nil
	}
	cacheMoot.Lock()
	defer cacheMoot.Unlock()
	os.MkdirAll(filepath.Dir(cachePath), 0744)
	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(cache)
}

func sum(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return ""
	}
	sum := hash.Sum(nil)

	s := fmt.Sprintf("%x", sum)
	return s
}

func addToCache(path string, cp cachedPlugin) {
	if cp.CheckSum == "" {
		cp.CheckSum = sum(path)
	}
	cacheMoot.Lock()
	defer cacheMoot.Unlock()
	cache[path] = cp
}
