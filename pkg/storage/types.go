package storage

import "fmt"

const (
	DefaultTablePagesCapacity   = 1024
	DefaultTableRecordsCapacity = 1024
)

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

// Table struct is the general structure representing our table.
// In general table implementation should support any types of records (should be interface)
// In our case it supports just ItemRecord
// Table can grow infinitely and just limited by size of memory (RAM)
// DB tables allocation problem very similar to what do language memory manager in runtime (heap allocation).
// Golang uses TCMalloc algorithm (very efficient) handling most of issues
// but it is good for limited size objects (biggest default object arena 64MB).
// If we use a simple slice object to store all ItemRecord we can have issues with
// performance related to allocation large size objects.
// Also using one large object is not efficient to solve fragmentation issues (and after defragmentation managing).
// We use slice instead of map because we need same access speed for all column (at least code and name).
// For access speed up we can use indexes.
// So we divide table structure for two level (it can be 3 and more levels
// based on expected size of table, actually need make performance test).
// Table has pages, pages stores records.
// Instead deleting record we mark it as deleted. It will delete later with clean process.
type table struct {
	pages []*tbPage
}

type tbPage struct {
	// size of records can not exceed max size (default value)
	// if need add more item should be created new page
	recs []ItemRecord
}

type ItemRecord struct {
	Code string
	Name string
	// Store $1000.99 as 100099
	Price int

	// flags
	isDeleted bool
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
