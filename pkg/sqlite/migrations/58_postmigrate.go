package migrations

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type schema58Migrator struct {
	migrator
}

func post58(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 58")

	m := schema58Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate()
}

func (m *schema58Migrator) migrate() error {
	if err := m.migrateConfig(); err != nil {
		return fmt.Errorf("failed to migrate config: %w", err)
	}

	return nil
}

// fromSnakeCase converts a string from snake_case to camelCase
func (m *schema58Migrator) fromSnakeCase(v string) string {
	var buf bytes.Buffer
	leadingUnderscore := true
	capvar := false
	for i, c := range v {
		switch {
		case c == '_' && !leadingUnderscore && i > 0:
			capvar = true
		case c == '_' && leadingUnderscore:
			buf.WriteRune(c)
		case capvar:
			buf.WriteRune(unicode.ToUpper(c))
			capvar = false
		default:
			leadingUnderscore = false
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

// fromSnakeCaseMap recursively converts a map using snake_case keys to camelCase keys
func (m *schema58Migrator) fromSnakeCaseMap(mm map[string]interface{}) map[string]interface{} {
	return m.fromSnakeCaseValue(mm).(map[string]interface{})
}

func (m *schema58Migrator) fromSnakeCaseValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[interface{}]interface{}:
		ret := cast.ToStringMap(v)
		for k, vv := range ret {
			adjKey := m.fromSnakeCase(k)
			ret[adjKey] = m.fromSnakeCaseValue(vv)
		}
		return ret
	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, vv := range v {
			adjKey := m.fromSnakeCase(k)
			ret[adjKey] = m.fromSnakeCaseValue(vv)
		}
		return ret
	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, vv := range v {
			ret[i] = m.fromSnakeCaseValue(vv)
		}
		return ret
	default:
		return v
	}
}

// renameKey renames a fully qualified key name in a map
func (m *schema58Migrator) renameKey(mm map[string]interface{}, from, to string) {
	nm := utils.NestedMap(mm)
	v, found := nm.Get(from)
	if !found {
		return
	}

	nm.Delete(from)
	nm.Set(to, v)
}

func (m *schema58Migrator) renameFrontPageContentKeys(ui map[string]interface{}) {
	frontPageContent, found := ui["frontPageContent"].([]interface{})
	if !found {
		return
	}

	for _, v := range frontPageContent {
		vm := v.(map[string]interface{})
		m.renameKey(vm, "savedfilterid", "savedFilterId")
		m.renameKey(vm, "sortby", "sortBy")
	}
}

func (m *schema58Migrator) migrateConfig() error {
	c := config.GetInstance()

	orgPath := c.GetConfigFile()

	if orgPath == "" {
		// no config file to migrate (usually in a test)
		return nil
	}

	ui := c.GetUIConfiguration()
	if len(ui) == 0 {
		// no UI config to migrate
		return nil
	}

	// save a backup of the original config file
	backupPath := fmt.Sprintf("%s.57.%s", orgPath, time.Now().Format("20060102_150405"))

	data, err := c.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal backup config: %w", err)
	}

	logger.Infof("Backing up config to %s", backupPath)
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup config: %w", err)
	}

	// migrate the plugin and UI configs from snake_case to camelCase
	if ui != nil {
		ui = m.fromSnakeCaseMap(ui)

		// find and rename specific frontEndPage keys
		m.renameFrontPageContentKeys(ui)

		c.SetUIConfiguration(ui)
	}

	plugins := c.GetAllPluginConfiguration()
	newPlugins := make(map[string]interface{})
	for key, value := range plugins {
		key = m.fromSnakeCase(key)
		newPlugins[key] = m.fromSnakeCaseMap(value)
	}

	c.SetInterface(config.PluginsSetting, newPlugins)
	if err := c.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(58, post58)
}
