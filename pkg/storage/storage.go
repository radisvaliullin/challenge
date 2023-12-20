package storage

import (
	"sync"
)

const (
	DefaultItemsTableCapacity = 1000_000

	// 16 alphnumeric characters + 3 dashes
	FullLenghtOfCode = 16 + 3
	// Code contains 4 alphanumeric parts devided by dash
	CodePartsNumber = 4
	// each part 4 characters
	CodePartLength = 4
)

type Config struct {
	Fixtures []ItemRecord
}

type Storage struct {
	config Config

	// Items Table mutex
	itemsMux   sync.Mutex
	itemsTable []ItemRecord
}

func New(conf Config) (*Storage, error) {
	s := &Storage{
		config:     conf,
		itemsTable: make([]ItemRecord, 0, DefaultItemsTableCapacity),
	}

	if err := s.applyFixtures(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) applyFixtures() error {
	areq := AddReq{
		Items: s.config.Fixtures,
	}
	_, err := s.Add(areq)
	if err != nil {
		return err
	}
	return nil
}
