package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AVZotov/draft-survey/internal/types"
)

var _ DictionariesRepository = (*DictionariesStore)(nil)

type DictionariesStore struct {
	Path string
}

func NewDictionariesStore(path string) (*DictionariesStore, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &DictionariesStore{
		Path: path,
	}, nil
}

func (j *DictionariesStore) GetPorts() (*[]types.Port, error) {
	const filename = "ports.json"
	path := filepath.Join(j.Path, filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}(file)

	decoder := json.NewDecoder(file)
	ports := &[]types.Port{}
	if err = decoder.Decode(ports); err != nil {
		return nil, err
	}

	return ports, nil
}

func (j *DictionariesStore) GetCountries() (*[]types.Country, error) {
	const filename = "countries.json"
	path := filepath.Join(j.Path, filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}(file)

	decoder := json.NewDecoder(file)
	countries := &[]types.Country{}
	if err = decoder.Decode(countries); err != nil {
		return nil, err
	}

	return countries, nil
}
