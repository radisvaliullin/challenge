package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/radisvaliullin/challenge/pkg/storage"
)

const (
	PathStoreFetch = "/store/fetch/"
)

type API struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) *API {
	api := &API{
		storage: storage,
	}
	return api
}

func (a *API) Muxer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", a.PingHandler)
	mux.HandleFunc("/store/add", a.AddHandler)
	mux.HandleFunc("/store/search", a.SearchHandler)
	mux.HandleFunc(PathStoreFetch, a.FetchHandler)
	mux.HandleFunc("/store/delete", a.DeleteHandler)
	return mux
}

func (a *API) PingHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Printf("api: ping: write error: %v", err)
	}
}

func (a *API) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeErrorCode(http.StatusNotFound, w)
		return
	}

	// decode payload
	addReq := AddReq{}
	if err := json.NewDecoder(r.Body).Decode(&addReq); err != nil {
		writeErrorCodeErr(http.StatusBadRequest, err, w)
		return
	}

	// request storage
	strgAddReq := storage.AddReq{}
	for _, item := range addReq.Items {
		strgItem := storage.ItemRecord{
			Code:  item.Code,
			Name:  item.Name,
			Price: int(item.Price * 100),
		}
		strgAddReq.Items = append(strgAddReq.Items, strgItem)
	}
	strgRes, err := a.storage.Add(strgAddReq)
	if err != nil {
		writeErrorCodeErr(http.StatusInternalServerError, err, w)
		return
	}

	// build resp
	res := AddRes{
		ItemCodes: strgRes.ItemCodes,
		ItemCount: strgRes.ItemCount,
	}

	writeResponse(res, w)
}

func (a *API) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeErrorCode(http.StatusNotFound, w)
		return
	}

	// decode payload
	srchReq := SearchReq{}
	if err := json.NewDecoder(r.Body).Decode(&srchReq); err != nil {
		writeErrorCodeErr(http.StatusBadRequest, err, w)
		return
	}

	// request storage
	strgRes := a.storage.Search(storage.SearchReq{Search: srchReq.Search})

	// build resp
	res := SearchRes{}
	for _, strgItem := range strgRes.Items {
		item := ItemRecord{
			Code:  strgItem.Code,
			Name:  strgItem.Name,
			Price: float64(strgItem.Price) / 100.0,
		}
		res.Items = append(res.Items, item)
	}

	writeResponse(res, w)
}

func (a *API) FetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeErrorCode(http.StatusNotFound, w)
		return
	}
	// parse path
	tailPath := strings.TrimPrefix(r.URL.Path, PathStoreFetch)
	if len(tailPath) == 0 || len(strings.Split(tailPath, "/")) != 1 {
		writeErrorCode(http.StatusNotFound, w)
		return
	}
	code := tailPath

	// request storage
	strgRes, err := a.storage.Fetch(storage.FetchReq{Code: code})
	if err != nil {
		writeErrorCodeErr(http.StatusInternalServerError, err, w)
		return
	}

	// build resp
	res := FetchRes{
		Item: ItemRecord{
			Code:  strgRes.Item.Code,
			Name:  strgRes.Item.Name,
			Price: float64(strgRes.Item.Price) / 100.0,
		},
	}

	writeResponse(res, w)

}

func (a *API) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		writeErrorCode(http.StatusNotFound, w)
		return
	}

	// decode payload
	delReq := DeleteReq{}
	if err := json.NewDecoder(r.Body).Decode(&delReq); err != nil {
		writeErrorCodeErr(http.StatusBadRequest, err, w)
		return
	}

	// request storage
	strgRes := a.storage.Delete(storage.DeleteReq{ItemCodes: delReq.ItemCodes})

	// build resp
	res := DeleteRes{
		ItemCount: strgRes.ItemCount,
	}

	writeResponse(res, w)
}

func writeErrorCode(code int, w http.ResponseWriter) {
	writeError(code, http.StatusText(code), "", w)
}

func writeErrorCodeErr(code int, err error, w http.ResponseWriter) {
	writeError(code, http.StatusText(code), fmt.Sprintf("%v", err), w)
}

func writeError(code int, errStr, msg string, w http.ResponseWriter) {
	// payload
	respErr := RespError{
		Code:    code,
		Error:   errStr,
		Message: msg,
	}
	respBytes, err := json.Marshal(&respErr)
	if err != nil {
		log.Printf("api: write error: response marshalling: %v", err)
	}
	// write
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(respBytes); err != nil {
		log.Printf("api: write error: %v", err)
	}
}

func writeResponse(res interface{}, w http.ResponseWriter) {
	// payload
	respBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("api: write response: response marshalling: %v", err)
	}
	// write
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		log.Printf("api: write response, error: %v", err)
	}
}
