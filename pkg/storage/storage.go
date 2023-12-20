package storage

import (
	"strings"
	"sync"
)

const (
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
	itemTable *table
	// code column index
	itbCodeIndex index
	// name column index
	itbNameIndex index
}

func New(conf Config) (*Storage, error) {
	s := &Storage{
		config:       conf,
		itemTable:    &table{},
		itbCodeIndex: make(index, DefaultIndexCapacity),
		itbNameIndex: make(index, DefaultIndexCapacity),
	}
	// init table
	s.itemTable.pages = make([]*tbPage, 0, DefaultTablePagesCapacity)
	page0 := &tbPage{
		recs: make([]ItemRecord, 0, DefaultTableRecordsCapacity),
	}
	s.itemTable.pages = append(s.itemTable.pages, page0)

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

// low level not thread-safe
func (s *Storage) getRecByCode(code string) (ItemRecord, bool) {
	if lrid, ok := s.itbCodeIndex[code]; ok {
		// code uniq, only last item can be undeleted
		rid := lrid[len(lrid)-1]
		rec := s.itemTable.pages[rid.page].recs[rid.rec]
		if rec.isDeleted {
			return ItemRecord{}, false
		}
		return rec, true
	}
	return ItemRecord{}, false
}

// low level not thread-safe
func (s *Storage) setRec(rec ItemRecord) {
	// set record to DB (add to latest page or add to new page)
	// check latest page records capacity (we have at least one page)
	lastPageIdx := len(s.itemTable.pages) - 1
	lastPage := s.itemTable.pages[lastPageIdx]
	if len(lastPage.recs) >= DefaultTableRecordsCapacity {
		newPage := &tbPage{
			recs: make([]ItemRecord, 0, DefaultTableRecordsCapacity),
		}
		s.itemTable.pages = append(s.itemTable.pages, newPage)
		lastPage = newPage
		lastPageIdx++
	}
	lastPage.recs = append(lastPage.recs, rec)
	recIdx := len(lastPage.recs) - 1

	// update index
	s.setIndex(rec, lastPageIdx, recIdx)
}

func (s *Storage) setIndex(rec ItemRecord, pageIdx, recIdx int) {

	id := rid{page: pageIdx, rec: recIdx}

	// code uniq, last item is actual
	s.itbCodeIndex[rec.Code] = append(s.itbCodeIndex[rec.Code], id)

	lName := strings.ToLower(rec.Name)
	lNameParts := strings.Fields(lName)
	// index whole name
	s.itbNameIndex[lName] = append(s.itbNameIndex[lName], id)
	// index each part and each part prefix
	// instead that ugly approach we should use b-tree index
	for _, part := range lNameParts {
		for i := 0; i < len(part); i++ {
			p := part[:len(part)-i]
			s.itbNameIndex[p] = append(s.itbNameIndex[p], id)
		}
	}
}

// low level not thread-safe
func (s *Storage) getRecByNamePref(namePref string) []ItemRecord {
	recs := []ItemRecord{}
	if lrid, ok := s.itbNameIndex[namePref]; ok {
		for _, id := range lrid {
			rec := s.itemTable.pages[id.page].recs[id.rec]
			// ignore deleted
			if rec.isDeleted {
				continue
			}
			recs = append(recs, rec)
		}
		return recs
	}
	return nil
}

// low level not thread-safe
func (s *Storage) delRecByCode(code string) bool {
	if lrid, ok := s.itbCodeIndex[code]; ok {
		// code uniq
		id := lrid[0]
		s.itemTable.pages[id.page].recs[id.rec].isDeleted = true

		return true
	}

	return false
}
