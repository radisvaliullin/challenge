package storage

import (
	"sync"
)

const (
	DefaultItemTableCapacity = 1000_000

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

	// Storage's single table and table mutex
	// In our case we have only one table (Item Table)
	// In general case we can have list of tables (or something like that)
	// Tables number usually do not have huge number so there is no scalability issue with that
	// Table records opposite it can infinitely grow until allowed by memory size (see Table struct details).
	itbMux    sync.Mutex
	itemTable []ItemRecord
}

func New(conf Config) (*Storage, error) {
	s := &Storage{
		config:    conf,
		itemTable: make([]ItemRecord, 0, DefaultItemTableCapacity),
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
