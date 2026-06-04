package storage

import (
	"database/sql"
	"encoding/json"

	"github.com/AVZotov/draft-survey/internal/types"
)

var _ SurveyRepository = (*SQLiteSurveyStore)(nil)

type SQLiteSurveyStore struct {
	db *sql.DB
}

func NewSQLiteSurveyStore(db *sql.DB) *SQLiteSurveyStore {
	return &SQLiteSurveyStore{
		db: db,
	}
}

func (s *SQLiteSurveyStore) Save(survey *types.Survey) error {
	const query = `INSERT INTO surveys (id, imo, data) VALUES (?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET imo = excluded.imo, data = excluded.data;`
	data, err := json.Marshal(survey)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, survey.ID, survey.VesselData.IMO, data)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLiteSurveyStore) Get(id string) (*types.Survey, error) {
	const query = `SELECT data FROM surveys WHERE id = ?`
	var data []byte
	if err := s.db.QueryRow(query, id).Scan(&data); err != nil {
		return nil, err
	}

	survey := new(types.Survey)
	if err := json.Unmarshal(data, survey); err != nil {
		return nil, err
	}

	return survey, nil
}

func (s *SQLiteSurveyStore) GetAll() ([]*types.Survey, error) {
	return get(s.db, `SELECT data FROM surveys ORDER BY created_at DESC;`)
}

func (s *SQLiteSurveyStore) Delete(id string) error {
	const query = `DELETE FROM surveys WHERE id = ?;`
	_, err := s.db.Exec(query, id)
	return err
}

func get(db *sql.DB, query string, args ...any) ([]*types.Survey, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	surveys := make([]*types.Survey, 0)
	for rows.Next() {
		survey := new(types.Survey)
		var data []byte
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, survey); err != nil {
			return nil, err
		}
		surveys = append(surveys, survey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return surveys, nil
}
