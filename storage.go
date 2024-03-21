package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type entry struct {
	name  string
	score int64
}

type storage interface {
	read(name string) (entry, error)
	write(e entry) error
	all() ([]entry, error)
	persistAndClose() error
}

type fileStorage struct {
	fname string
	data  map[string]int64
}

func openFileStorage(fname string) (storage, error) {
	s := &fileStorage{fname: fname}

	data, err := os.ReadFile(fname)
	if errors.Is(err, os.ErrNotExist) {
		s.data = map[string]int64{}
		return s, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not read storage file: %w", err)
	}

	err = json.Unmarshal(data, &s.data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal storage data: %w", err)
	}

	return s, nil
}

func (s *fileStorage) read(name string) (entry, error) {
	score := s.data[name]

	return entry{name: name, score: score}, nil
}

func (s *fileStorage) write(e entry) error {
	s.data[e.name] = e.score

	return nil
}

func (s *fileStorage) all() ([]entry, error) {
	ee := make([]entry, 0, len(s.data))

	for k, v := range s.data {
		ee = append(ee, entry{name: k, score: v})
	}

	return ee, nil
}

func (s *fileStorage) persistAndClose() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("could not marshal storage data: %w", err)
	}

	err = os.WriteFile(s.fname, data, 0666)
	if err != nil {
		return fmt.Errorf("could not write storage file: %w", err)
	}

	return nil
}
