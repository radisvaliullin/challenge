package storage

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Test positive cases without error
func TestStorageFetchPositive(t *testing.T) {

	// init storage
	conf := Config{
		Fixtures: getTestConfigFixtures(),
	}
	storage, err := New(conf)
	if err != nil {
		t.Fatalf("storage fetch postive: new storage: %v", err)
	}

	// test cases
	testCases := []struct {
		Code    string
		ExpResp FetchRes
	}{
		{
			Code: "A12T-4GH7-QPL9-3N4M",
			ExpResp: FetchRes{Item: ItemRecord{
				Code:  "A12T-4GH7-QPL9-3N4M",
				Name:  "Lettuce",
				Price: 3_46,
			}},
		},
		{
			Code: "E5T6-9UI3-TH15-QR88",
			ExpResp: FetchRes{Item: ItemRecord{
				Code:  "E5T6-9UI3-TH15-QR88",
				Name:  "Peach",
				Price: 2_99,
			}},
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			fReq := FetchReq{
				Code: tc.Code,
			}
			resp, err := storage.Fetch(fReq)
			if err != nil {
				t.Fatalf("storage fetch postive: fetch code %v, error: expResp - %+v, err - %+v", tc.Code, tc.ExpResp, err)
			}
			if !reflect.DeepEqual(tc.ExpResp, resp) {
				t.Fatalf("storage fetch postive: fetch code %v, wrong response: expResp - %+v, resp - %+v", tc.Code, tc.ExpResp, resp)
			}
			t.Logf("storage fetch postive: fetch code %v: expResp - %+v, resp - %+v", tc.Code, tc.ExpResp, resp)
		})
	}
}

func TestStorageAddPositive(t *testing.T) {

	// init storage
	conf := Config{
		Fixtures: getTestConfigFixtures(),
	}
	storage, err := New(conf)
	if err != nil {
		t.Fatalf("storage add postive: new storage: %v", err)
	}

	// test cases
	testCases := []struct {
		AddReq  AddReq
		ExpResp AddRes
	}{
		{
			AddReq: AddReq{Items: []ItemRecord{
				{
					Code:  "AA11-BB22-CC33-DD44",
					Name:  "Pomegranate",
					Price: 5_00,
				},
			}},
			ExpResp: AddRes{
				ItemCodes: []string{"AA11-BB22-CC33-DD44"},
				ItemCount: 1,
			},
		},
		{
			AddReq: AddReq{Items: []ItemRecord{
				{
					Code:  "A111-B222-C333-D444",
					Name:  "GMO Pomegranate",
					Price: 6_00,
				},
				{
					Code:  "A112-B223-C334-D445",
					Name:  "GMO Grapefruit",
					Price: 6_99,
				},
			}},
			ExpResp: AddRes{
				ItemCodes: []string{"A111-B222-C333-D444", "A112-B223-C334-D445"},
				ItemCount: 2,
			},
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			resp, err := storage.Add(tc.AddReq)
			if err != nil {
				t.Fatalf("storage add postive: req %+v, fail: %v", tc.AddReq, err)
			}
			if !reflect.DeepEqual(tc.ExpResp, resp) {
				t.Fatalf("storage add postive: req %+v, wrong response: expResp - %+v, resp - %+v", tc.AddReq, tc.ExpResp, resp)
			}
			t.Logf("storage add postive: req %+v: expResp - %+v, resp - %+v", tc.AddReq, tc.ExpResp, resp)
		})
	}
}

func TestStorageAddNegative(t *testing.T) {

	// init storage
	conf := Config{
		Fixtures: getTestConfigFixtures(),
	}
	storage, err := New(conf)
	if err != nil {
		t.Fatalf("storage add negative: new storage: %v", err)
	}

	// test cases
	testCases := []struct {
		AddReq AddReq
		ExpErr error
	}{
		{
			AddReq: AddReq{Items: []ItemRecord{
				{
					Code:  "AA11-BB22-CC33-DD4@",
					Name:  "Pomegranate",
					Price: 5_00,
				},
			}},
			ExpErr: ErrStorageCodePartNotAlphaNum,
		},
		{
			AddReq: AddReq{Items: []ItemRecord{
				{
					Code:  "A12T-4GH7-QPL9-3N4M",
					Name:  "GMO Pomegranate",
					Price: 6_00,
				},
			}},
			ExpErr: ErrStorageAddReqCodeDuplDB,
		},
		{
			AddReq: AddReq{Items: []ItemRecord{
				{
					Code:  "A12T-4GH7-QPL9-3N4M",
					Name:  "GMO Pomegranate",
					Price: 6_00,
				},
				{
					Code:  "A12T-4GH7-QPL9-3N4M",
					Name:  "GMO Pomegranate",
					Price: 6_00,
				},
			}},
			ExpErr: ErrStorageAddReqDuplCode,
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			_, err := storage.Add(tc.AddReq)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("storage add negative: unexpected error: req %+v, expErr - %v, err - %v", tc.AddReq, tc.ExpErr, err)
			}
			t.Logf("storage add negative: req %+v: expErr - %+v, err - %+v", tc.AddReq, tc.ExpErr, err)
		})
	}
}

func TestStorageSearchPositive(t *testing.T) {

	// init storage
	conf := Config{
		Fixtures: getTestConfigFixtures(),
	}
	storage, err := New(conf)
	if err != nil {
		t.Fatalf("storage search postive: new storage: %v", err)
	}

	// test cases
	testCases := []struct {
		Search  string
		ExpResp SearchRes
	}{
		{
			Search: "Pepp Pepp",
			ExpResp: SearchRes{Items: []ItemRecord{
				{
					Code:  "YRT6-72AS-K736-L4AR",
					Name:  "Green Pepper",
					Price: 79,
				},
			}},
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			sReq := SearchReq{
				Search: tc.Search,
			}
			resp := storage.Search(sReq)
			if !reflect.DeepEqual(tc.ExpResp, resp) {
				t.Fatalf("storage search postive: search %v, wrong response: expResp - %+v, resp - %+v", tc.Search, tc.ExpResp, resp)
			}
			t.Logf("storage search postive: search %v: expResp - %+v, resp - %+v", tc.Search, tc.ExpResp, resp)
		})
	}
}

func TestStorageDeletePositive(t *testing.T) {

	// init storage
	conf := Config{
		Fixtures: getTestConfigFixtures(),
	}
	storage, err := New(conf)
	if err != nil {
		t.Fatalf("storage delete postive: new storage: %v", err)
	}

	// test cases
	testCases := []struct {
		Codes       []string
		ExpResp     DeleteRes
		ExpFetchErr error
	}{
		{
			Codes:       []string{"A12T-4GH7-QPL9-3N4M"},
			ExpResp:     DeleteRes{ItemCount: 1},
			ExpFetchErr: ErrStorageItemNotFound,
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			dReq := DeleteReq{
				ItemCodes: tc.Codes,
			}
			resp := storage.Delete(dReq)
			if !reflect.DeepEqual(tc.ExpResp, resp) {
				t.Fatalf("storage delete postive: codes %v, wrong response: expResp - %+v, resp - %+v", tc.Codes, tc.ExpResp, resp)
			}
			_, err := storage.Fetch(FetchReq{Code: tc.Codes[0]})
			if !errors.Is(err, tc.ExpFetchErr) {
				t.Fatalf("storage delete postive: codes %v, wrong fetch error: expErr - %+v, err - %+v", tc.Codes, tc.ExpFetchErr, err)
			}
			t.Logf("storage delete postive: codes %v, expResp - %+v, resp - %+v", tc.Codes, tc.ExpResp, resp)
		})
	}
}

func TestIsAlphaNumeric(t *testing.T) {

	testCases := []struct {
		str        string
		isAlphaNum bool
	}{
		{str: "A12T", isAlphaNum: true},
		{str: "4GH7", isAlphaNum: true},
		{str: "QPL9", isAlphaNum: true},
		{str: "3N4M", isAlphaNum: true},
		{str: "3N4-", isAlphaNum: false},
		{str: "3N4@", isAlphaNum: false},
	}

	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {
			ok := IsAlphaNumeric(tc.str)
			if ok != tc.isAlphaNum {
				t.Fatalf("is alphanumeric: fail: str - %v, res - %v, expRes - %v", tc.str, ok, tc.isAlphaNum)
			}
			t.Logf("is alphanumeric: ok: str - %v, res - %v, expRes - %v", tc.str, ok, tc.isAlphaNum)
		})
	}
}

func getTestConfigFixtures() []ItemRecord {
	return []ItemRecord{
		{
			Code:  "A12T-4GH7-QPL9-3N4M",
			Name:  "Lettuce",
			Price: 3_46,
		},
		{
			Code:  "E5T6-9UI3-TH15-QR88",
			Name:  "Peach",
			Price: 2_99,
		},
		{
			Code:  "YRT6-72AS-K736-L4AR",
			Name:  "Green Pepper",
			Price: 79,
		},
		{
			Code:  "TQ4C-VV6T-75ZX-1RMR",
			Name:  "Gala Apple",
			Price: 3_59,
		},
	}
}
