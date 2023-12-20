package storage

import (
	"strings"
)

/*
Storage public api
*/

// Add
// item price 99.99 as 9999
func (s *Storage) Add(addReq AddReq) (AddRes, error) {
	// validate input
	err := validateAddRequest(addReq)
	if err != nil {
		return AddRes{}, err
	}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	// check dublicates
	for _, newItem := range addReq.Items {
		for _, page := range s.itemTable.pages {
			for _, item := range page.recs {
				// skip deleted item
				if item.isDeleted {
					continue
				}
				// check code matching
				if newItem.Code == item.Code {
					return AddRes{}, ErrStorageAddReqCodeDuplDB
				}
			}
		}
	}

	// set items to DB (add to latest page or add to new page)
	codes := []string{}
	for _, item := range addReq.Items {
		// check latest page records capacity (we have at least one page)
		lastPage := s.itemTable.pages[len(s.itemTable.pages)-1]
		if len(lastPage.recs) >= DefaultTableRecordsCapacity {
			newPage := &tbPage{
				recs: make([]ItemRecord, 0, DefaultTableRecordsCapacity),
			}
			s.itemTable.pages = append(s.itemTable.pages, newPage)
			lastPage = newPage
		}
		lastPage.recs = append(lastPage.recs, item)
		codes = append(codes, item.Code)
	}
	res := AddRes{
		ItemCodes: codes,
		ItemCount: len(codes),
	}
	return res, nil
}

func (s *Storage) Search(srchReq SearchReq) SearchRes {
	phrases := strings.Fields(srchReq.Search)
	res := SearchRes{}
	items := []ItemRecord{}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	// find items
	for _, p := range phrases {
		for _, page := range s.itemTable.pages {
			for _, item := range page.recs {
				// skip deleted item
				if item.isDeleted {
					continue
				}
				// compare
				if strings.Contains(strings.ToLower(item.Name), strings.ToLower(p)) {
					items = append(items, item)
				}
			}
		}
	}

	// clean duplicates
	for _, item := range items {
		isAdded := false
		for _, ritem := range res.Items {
			if item.Code == ritem.Code {
				isAdded = true
				break
			}
		}
		if !isAdded {
			res.Items = append(res.Items, item)
		}
	}

	return res
}

func (s *Storage) Fetch(fReq FetchReq) (FetchRes, error) {
	res := FetchRes{}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	isFind := false
	for _, page := range s.itemTable.pages {
		for _, item := range page.recs {
			// skip deleted item
			if item.isDeleted {
				continue
			}
			// compare
			if item.Code == fReq.Code {
				res.Item = item
				isFind = true
				break
			}
		}
	}
	if isFind {
		return res, nil
	}
	return res, ErrStorageItemNotFound
}

func (s *Storage) Delete(delReq DeleteReq) DeleteRes {
	res := DeleteRes{}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	for _, dCode := range delReq.ItemCodes {
		for _, page := range s.itemTable.pages {
			for idx, item := range page.recs {
				// skip deleted item
				if item.isDeleted {
					continue
				}
				// compare
				if item.Code == dCode {
					// set as deleted
					page.recs[idx].isDeleted = true
					res.ItemCount++
				}
			}
		}
	}

	return res
}

func validateAddRequest(addReq AddReq) error {
	for i, item := range addReq.Items {
		// validate code
		err := validateItemCode(item.Code)
		if err != nil {
			return err
		}
		// validate name
		err = validateItemName(item.Name)
		if err != nil {
			return err
		}
		// exclude duplicated codes
		for _, item2 := range addReq.Items[i+1:] {
			if item.Code == item2.Code {
				return ErrStorageAddReqDuplCode
			}
		}
		// clean data
		item.Code = strings.ToUpper(item.Code)
		item.Name = strings.TrimSpace(item.Name)
	}
	return nil
}

func validateItemCode(code string) error {
	if len(code) != FullLenghtOfCode {
		return ErrStorageCodeWrongLen
	}
	codeParts := strings.Split(code, "-")
	if len(codeParts) != CodePartsNumber {
		return ErrStorageCodePartsWrongNum
	}
	for _, p := range codeParts {
		if len(p) != CodePartLength {
			return ErrStorageCodePartWrongLen
		}
		if !IsAlphaNumeric(p) {
			return ErrStorageCodePartNotAlphaNum
		}
	}
	return nil
}

func validateItemName(name string) error {
	nameParts := strings.Fields(name)
	for _, p := range nameParts {
		if !IsAlphaNumeric(p) {
			return ErrStorageNamePartNotAlphaNum
		}
	}
	return nil
}
