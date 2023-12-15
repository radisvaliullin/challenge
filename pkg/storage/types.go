package storage

import "fmt"

var (
	ErrStorageCodePartNotAlphaNum = StorageError{Message: "code, one of part not alphanumeric"}
	ErrStorageNamePartNotAlphaNum = StorageError{Message: "name, one of part not alphanumeric"}
	ErrStorageCodePartWrongLen    = StorageError{Message: "code, one of part has wrong length"}
	ErrStorageCodePartsWrongNum   = StorageError{Message: "code, wrong parts number"}
	ErrStorageCodeWrongLen        = StorageError{Message: "code, wrong length"}
	ErrStorageAddReqDuplCode      = StorageError{Message: "add request: validate: duplicated codes"}
	ErrStorageAddReqCodeDuplDB    = StorageError{Message: "add request: code duplicated with db"}
	ErrStorageItemNotFound        = StorageError{Message: "fetch request: item not found"}
)

type ItemRecord struct {
	Code string
	Name string
	// Store $1000.99 as 100099
	Price int
}

// AddReq
// item price 99.99 as 9999
type AddReq struct {
	Items []ItemRecord
}

type AddRes struct {
	ItemCodes []string
	ItemCount int
}

type SearchReq struct {
	Search string
}

type SearchRes struct {
	Items []ItemRecord
}

type FetchReq struct {
	Code string
}

type FetchRes struct {
	Item ItemRecord
}

type DeleteReq struct {
	ItemCodes []string
}

type DeleteRes struct {
	ItemCount int
}

type StorageError struct {
	Message string
}

func (m StorageError) Error() string {
	return fmt.Sprintf("storage: type error: %v", m.Message)
}

var _ error = StorageError{}
