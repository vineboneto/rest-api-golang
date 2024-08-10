package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTableTenant, downCreateTableTenant)
}

func upCreateTableTenant(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS tbl_tenant (
			id serial primary key,
			nome text
		)
	`)
	return err
}

func downCreateTableTenant(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE IF EXISTS tbl_tenant
	`)
	return err
}
