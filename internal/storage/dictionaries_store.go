package storage

import (
	"embed"
	"encoding/json"

	"github.com/AVZotov/draft-survey/internal/types"
)

var _ DictionariesRepository = (*DictionariesStore)(nil)

type DictionariesStore struct {
	fs embed.FS
}

func NewDictionariesStore(fs embed.FS) *DictionariesStore {
	return &DictionariesStore{fs: fs}
}

func (j *DictionariesStore) GetPorts() (*[]types.Port, error) {
	data, err := j.fs.ReadFile("data/dictionaries/ports.json")
	if err != nil {
		return nil, err
	}

	ports := &[]types.Port{}
	if err = json.Unmarshal(data, ports); err != nil {
		return nil, err
	}

	return ports, nil
}

func (j *DictionariesStore) GetCountries() (*[]types.Country, error) {
	data, err := j.fs.ReadFile("data/dictionaries/countries.json")
	if err != nil {
		return nil, err
	}

	countries := &[]types.Country{}
	if err = json.Unmarshal(data, countries); err != nil {
		return nil, err
	}

	return countries, nil
}
