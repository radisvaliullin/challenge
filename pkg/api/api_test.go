package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/radisvaliullin/challenge/pkg/storage"
)

func TestAPIAddPositive(t *testing.T) {

	// init storage
	conf := storage.Config{
		Fixtures: storage.GetTestConfigFixtures(),
	}
	strg, err := storage.New(conf)
	if err != nil {
		t.Fatalf("api add positive: new storage: %v", err)
	}

	// init api
	apiOjb := New(strg)

	// test cases
	testCases := []struct {
		AddReqRaw string
		ExpResp   AddRes
	}{
		{
			AddReqRaw: `{"items":[{"code":"1111-2222-3333-4444","name":"apple","price":3.44}]}`,
			ExpResp: AddRes{
				ItemCodes: []string{"1111-2222-3333-4444"},
				ItemCount: 1,
			},
		},
	}

	// run test cases
	for i, tc := range testCases {
		tc := tc

		tn := fmt.Sprintf("%v", i)
		t.Run(tn, func(t *testing.T) {

			// make request
			req := httptest.NewRequest("POST", "/store/add", strings.NewReader(tc.AddReqRaw))
			w := httptest.NewRecorder()
			apiOjb.AddHandler(w, req)

			// check result
			reqResp := w.Result()
			resp := AddRes{}
			err := json.NewDecoder(reqResp.Body).Decode(&resp)
			if err != nil {
				t.Fatalf("api add postive: decode err - %v, req %+v, ", err, tc.AddReqRaw)
			}
			if !reflect.DeepEqual(tc.ExpResp, resp) {
				t.Fatalf("api add postive: req %+v, wrong response: expResp - %+v, resp - %+v", tc.AddReqRaw, tc.ExpResp, resp)
			}
			t.Logf("api add postive: req %+v, expResp - %+v, resp - %+v", tc.AddReqRaw, tc.ExpResp, resp)
		})
	}
}
