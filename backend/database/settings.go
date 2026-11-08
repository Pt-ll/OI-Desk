package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Setting struct {
	Key       string
	ValueJSON string
	ValueType string
	CreatedAt string
	UpdatedAt string
}

type SettingsStore struct {
	db *sql.DB
}

func NewSettingsStore(db *sql.DB) *SettingsStore { return &SettingsStore{db: db} }

func (s *SettingsStore) Set(ctx context.Context, key, valueJSON, valueType string) error {
	if s == nil || s.db == nil {
		return errors.New("settings store is not initialized")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("setting key is empty")
	}
	if !json.Valid([]byte(valueJSON)) {
		return errors.New("setting value is not valid JSON")
	}
	if !validSettingType(valueType) {
		return fmt.Errorf("unsupported setting value type %q", valueType)
	}
	now := utcNow()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO app_settings(key, value_json, value_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value_json = excluded.value_json, value_type = excluded.value_type, updated_at = excluded.updated_at
	`, key, valueJSON, valueType, now, now)
	if err != nil {
		return fmt.Errorf("set setting %q: %w", key, err)
	}
	return nil
}

func (s *SettingsStore) Get(ctx context.Context, key string) (Setting, error) {
	if s == nil || s.db == nil {
		return Setting{}, errors.New("settings store is not initialized")
	}
	var result Setting
	err := s.db.QueryRowContext(ctx, `SELECT key, value_json, value_type, created_at, updated_at FROM app_settings WHERE key = ?`, key).
		Scan(&result.Key, &result.ValueJSON, &result.ValueType, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Setting{}, fmt.Errorf("setting %q: %w", key, sql.ErrNoRows)
	}
	if err != nil {
		return Setting{}, fmt.Errorf("get setting %q: %w", key, err)
	}
	return result, nil
}

func (s *SettingsStore) Delete(ctx context.Context, key string) error {
	if s == nil || s.db == nil {
		return errors.New("settings store is not initialized")
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM app_settings WHERE key = ?`, key); err != nil {
		return fmt.Errorf("delete setting %q: %w", key, err)
	}
	return nil
}

func validSettingType(valueType string) bool {
	switch valueType {
	case "string", "number", "boolean", "object", "array", "null":
		return true
	default:
		return false
	}
}
