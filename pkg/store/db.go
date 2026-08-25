package store

import (
	"errors"
	"fmt"
	"task_tracker/pkg/model"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // "_" - вызывает только init() модуля и регистрирует "postgres" драйвер
)

var schema = `
CREATE TABLE IF NOT EXISTS task (
		id SERIAL PRIMARY KEY,
		name text NOT NULL,
		status text NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP
	);
`

// database handler
type DBHandler struct {
	db *sqlx.DB
}

func NewDBHandler(connStr string) (DataHandler, error) {
	if connStr == "" {
		return nil, errors.New("DB connection string not set")
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.MustExec(schema)

	return &DBHandler{db}, nil
}

// returns id of a new task
func (h *DBHandler) Add(t *model.Task) (int, error) {
	var err error

	tx := h.db.MustBegin()
	defer func() {
		if pan := recover(); pan != nil || err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	var id int
	rows, err := tx.NamedQuery("INSERT INTO task (name, status, created_at) VALUES (:name, :status, :created_at) RETURNING id", t)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return id, nil
}

func (h *DBHandler) Update(id int, field string, val string) error {

	var err error

	tx := h.db.MustBegin()
	defer func() {
		if pan := recover(); pan != nil || err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	// whitelist allowed fields to avoid SQL injection
	allowed := map[string]bool{
		"name":   true,
		"status": true,
	}

	if !allowed[field] {
		return errors.New("invalid field for update")
	}

	query := fmt.Sprintf("UPDATE task SET %s = $1, updated_at = $2 WHERE id = $3", field)
	_, err = tx.Exec(query, val, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (h *DBHandler) Delete(id int) error {
	_, err := h.db.Exec("DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (h *DBHandler) Select(field string, val string) ([]model.Task, error) {
	allowed := map[string]bool{
		"name":   true,
		"status": true,
	}

	if !allowed[field] {
		return nil, errors.New("invalid field for update")
	}

	var tasks []model.Task
	var query string
	if val != "" {
		query = fmt.Sprintf("SELECT * FROM task WHERE %s = $1 ORDER BY id ASC", field)
		err := h.db.Select(&tasks, query, val)
		if err != nil {
			return nil, err
		}
	} else {
		query = "SELECT * FROM task ORDER BY id ASC"
		err := h.db.Select(&tasks, query)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

var _ DataHandler = (*DBHandler)(nil)
