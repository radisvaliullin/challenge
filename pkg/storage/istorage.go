package storage

var _ IStorage = (*Storage)(nil)

// interface for mocking storage if need
type IStorage interface {
	Add(AddReq) (AddRes, error)
	Search(SearchReq) SearchRes
	Fetch(FetchReq) (FetchRes, error)
	Delete(DeleteReq) DeleteRes
}
