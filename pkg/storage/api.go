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

	// check duplicates
	for _, newItem := range addReq.Items {
		if _, ok := s.getRecByCode(newItem.Code); ok {
			// if ok we have duplicate
			return AddRes{}, ErrStorageAddReqCodeDuplDB
		}
	}

	// set items to DB (add to latest page or add to new page)
	codes := []string{}
	for _, item := range addReq.Items {
		s.setRec(item)
		codes = append(codes, item.Code)
	}

	// response
	res := AddRes{
		ItemCodes: codes,
		ItemCount: len(codes),
	}
	return res, nil
}

func (s *Storage) Search(srchReq SearchReq) SearchRes {
	parts := strings.Fields(strings.ToLower(srchReq.Search))
	res := SearchRes{}
	items := []ItemRecord{}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	// find items
	for _, p := range parts {

		if recs := s.getRecByNamePref(p); len(recs) > 0 {
			items = append(items, recs...)
		}
		// optionally we can use contains search but it not supports by index
		// and will require full search
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

	if rec, ok := s.getRecByCode(strings.ToUpper(fReq.Code)); ok {
		res.Item = rec
		return res, nil
	}
	return res, ErrStorageItemNotFound
}

func (s *Storage) Delete(delReq DeleteReq) DeleteRes {
	res := DeleteRes{}

	s.itbMux.Lock()
	defer s.itbMux.Unlock()

	for _, dCode := range delReq.ItemCodes {
		dCode = strings.ToUpper(dCode)
		if ok := s.delRecByCode(dCode); ok {
			res.ItemCount++
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
