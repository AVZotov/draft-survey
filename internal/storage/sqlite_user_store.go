package storage

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/AVZotov/draft-survey/internal/types"
)

var _ UserRepository = (*SQLiteUserStore)(nil)

type SQLiteUserStore struct {
	db *sql.DB
}

func NewSQLiteUserStore(db *sql.DB) *SQLiteUserStore {
	return &SQLiteUserStore{db: db}
}

func (s *SQLiteUserStore) Save(user *types.User) error {
	const query = `INSERT INTO users (id, data) VALUES (?, ?)
		ON CONFLICT (id) DO UPDATE SET data = excluded.data;`

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(query, 1, data)

	return err
}

func (s *SQLiteUserStore) SaveSignature(data []byte) error {
	_, err := s.db.Exec(`UPDATE users SET signature = ? WHERE id = 1`, data)
	return err
}

func (s *SQLiteUserStore) Get() (*types.User, error) {
	user := new(types.User)
	var data []byte
	err := s.db.QueryRow(`SELECT data FROM users WHERE id = 1`).Scan(&data)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SQLiteUserStore) Delete() error {
	const query = `DELETE FROM users WHERE id = 1;`
	_, err := s.db.Exec(query)

	return err
}
