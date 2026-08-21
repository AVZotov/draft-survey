package storage

import (
	"database/sql"
	"encoding/json"
	"time"

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
	data, err := json.Marshal(survey)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO surveys (id, vessel_name, imo, status, operation, cargo, created_at, data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			vessel_name = excluded.vessel_name,
			imo         = excluded.imo,
			status      = excluded.status,
			operation   = excluded.operation,
			cargo       = excluded.cargo,
			data        = excluded.data`,
		survey.ID,
		survey.VesselData.Name,
		survey.VesselData.IMO,
		string(survey.Status),
		string(survey.CargoOperation.Operation),
		survey.CargoOperation.Cargo,
		survey.CreatedAt.Format(time.RFC3339),
		string(data),
	)
	return err
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

func (s *SQLiteSurveyStore) GetPage(limit, offset int) ([]*types.Survey, error) {
	rows, err := s.db.Query(`SELECT data FROM surveys ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []*types.Survey
	for rows.Next() {
		var data string
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		var survey types.Survey
		if err = json.Unmarshal([]byte(data), &survey); err != nil {
			return nil, err
		}
		surveys = append(surveys, &survey)
	}
	return surveys, nil
}

func (s *SQLiteSurveyStore) GetStats() (types.SurveyStats, error) {
	var stats types.SurveyStats
	err := s.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'complete' THEN 1 END) as complete,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress,
			COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft
		FROM surveys
	`).Scan(&stats.Total, &stats.Complete, &stats.InProgress, &stats.Draft)
	return stats, err
}

func (s *SQLiteSurveyStore) Delete(id string) error {
	const query = `DELETE FROM surveys WHERE id = ?;`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *SQLiteSurveyStore) Search(filter types.SurveyFilter) ([]*types.Survey, error) {
	query := `SELECT data FROM surveys WHERE 1=1`
	args := []any{}

	if filter.Query != "" {
		query += ` AND (vessel_name LIKE ? OR imo LIKE ?)`
		q := "%" + filter.Query + "%"
		args = append(args, q, q)
	}
	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	if filter.Operation != "" {
		query += ` AND operation = ?`
		args = append(args, filter.Operation)
	}
	if filter.Cargo != "" {
		query += ` AND cargo = ?`
		args = append(args, filter.Cargo)
	}
	if !filter.From.IsZero() {
		query += ` AND created_at >= ?`
		args = append(args, filter.From.Format(time.RFC3339))
	}
	if !filter.To.IsZero() {
		query += ` AND created_at <= ?`
		args = append(args, filter.To.Format(time.RFC3339))
	}

	query += ` ORDER BY created_at DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []*types.Survey
	for rows.Next() {
		var data string
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		var survey types.Survey
		if err = json.Unmarshal([]byte(data), &survey); err != nil {
			return nil, err
		}
		surveys = append(surveys, &survey)
	}
	return surveys, nil
}
