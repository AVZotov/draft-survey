package storage

import (
	"database/sql"

	"github.com/AVZotov/draft-survey/internal/types"
)

var _ DictionariesRepository = (*DictionariesStore)(nil)

type DictionariesStore struct {
	db *sql.DB
}

func NewDictionariesStore(db *sql.DB) *DictionariesStore {
	return &DictionariesStore{db: db}
}

func (s *DictionariesStore) GetCountries() (*[]types.Country, error) {
	rows, err := s.db.Query(`SELECT code, name FROM countries ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []types.Country
	for rows.Next() {
		var c types.Country
		if err = rows.Scan(&c.CountryCode, &c.Name); err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return &countries, nil
}

func (s *DictionariesStore) GetPorts(countryCode string) (*[]types.Port, error) {
	rows, err := s.db.Query(`SELECT locode, name, country_code, coordinates FROM ports WHERE country_code = ? ORDER BY name`, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []types.Port
	for rows.Next() {
		var p types.Port
		if err = rows.Scan(&p.Locode, &p.Name, &p.CountryCode, &p.Coordinates); err != nil {
			return nil, err
		}
		ports = append(ports, p)
	}
	return &ports, nil
}

func (s *DictionariesStore) GetCargoTypes() ([]string, error) {
	rows, err := s.db.Query(`SELECT name FROM cargo_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		types = append(types, name)
	}
	return types, nil
}

func (s *DictionariesStore) GetPacking() ([]string, error) {
	rows, err := s.db.Query(`SELECT name FROM packing ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	return items, nil
}
