package responses

// Pagination response
type Pagination[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}
