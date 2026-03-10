package main

import (
	"database/sql"
	"fmt"
)

func ensureAdminSchema(db *sql.DB) error {
	// Ensure `role` and `status` exist on `user`.
	if err := ensureColumnExists(db, "user", "role", "ALTER TABLE user ADD COLUMN role VARCHAR(20) DEFAULT 'user' AFTER email"); err != nil {
		return err
	}
	if err := ensureColumnExists(db, "user", "status", "ALTER TABLE user ADD COLUMN status VARCHAR(20) DEFAULT 'active' AFTER role"); err != nil {
		return err
	}

	// Ensure admin_log exists.
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS admin_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  actor_user_id BIGINT NOT NULL,
  action TEXT NOT NULL,
  ip VARCHAR(64),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_admin_log_created_at (created_at),
  INDEX idx_admin_log_actor (actor_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		return err
	}

	// Best-effort: mark first user as admin for UI display consistency.
	_, _ = db.Exec("UPDATE user SET role = 'admin' WHERE id = 1")

	return nil
}

func ensureColumnExists(db *sql.DB, table string, column string, alterSQL string) error {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?",
		table,
		column,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check column %s.%s: %w", table, column, err)
	}
	if count > 0 {
		return nil
	}
	_, err = db.Exec(alterSQL)
	if err != nil {
		return fmt.Errorf("alter table for %s.%s: %w", table, column, err)
	}
	return nil
}
