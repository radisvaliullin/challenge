package api

type RespError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ItemRecord struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type AddReq struct {
	Items []ItemRecord `json:"items"`
}

type AddRes struct {
	ItemCodes []string `json:"item_codes"`
	ItemCount int      `json:"item_count"`
}

type SearchReq struct {
	Search string `json:"search"`
}

type SearchRes struct {
	Items []ItemRecord `json:"items"`
}

type FetchRes struct {
	Item ItemRecord `json"item"`
}

type DeleteReq struct {
	ItemCodes []string `json:"item_codes"`
}

type DeleteRes struct {
	ItemCount int `json:"item_count"`
}
