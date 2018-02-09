package client

import "errors"


type Config struct {
	AutoApply       bool   `json:"auto_apply"`
	IsTransactional bool   `json:"is_transactional"`
	MigrationPath   string `json:"migration_path"`
	DatabaseName    string `json:"database_name"`
	TableName       string `json:"table_name"`
}

func (c *Config) Prepare() error {
	if c.MigrationPath == "" {
		c.MigrationPath = "migrations"
	}

	if c.DatabaseName == "" {
		return errors.New("config: database name is required")
	}

	if c.TableName == "" {
		c.TableName = "migration_history"
	}

	return nil
}
