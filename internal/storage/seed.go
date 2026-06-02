package storage

import (
	"database/sql"
	"embed"
	"encoding/json"
)

type countryJSON struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type portJSON struct {
	Locode      string `json:"locode"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Coordinates string `json:"coordinates"`
}

func seedDictionaries(db *sql.DB, fs embed.FS) error {
	if err := seedCountries(db, fs); err != nil {
		return err
	}
	return seedPorts(db, fs)
}

func seedCountries(db *sql.DB, fs embed.FS) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM countries`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data, err := fs.ReadFile("data/dictionaries/countries.json")
	if err != nil {
		return err
	}

	var countries []countryJSON
	if err = json.Unmarshal(data, &countries); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO countries (code, name) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range countries {
		if _, err = stmt.Exec(c.Code, c.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func seedPorts(db *sql.DB, fs embed.FS) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ports`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data, err := fs.ReadFile("data/dictionaries/ports.json")
	if err != nil {
		return err
	}

	var ports []portJSON
	if err = json.Unmarshal(data, &ports); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO ports (locode, name, country_code, coordinates) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range ports {
		if _, err = stmt.Exec(p.Locode, p.Name, p.CountryCode, p.Coordinates); err != nil {
			return err
		}
	}

	return tx.Commit()
}
