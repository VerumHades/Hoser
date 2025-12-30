package util

type PaginatedResponse[T any, CursorType any] struct {
	Items  []T    `json:"items"`
	Cursor string `json:"cursor,omitempty"` // opaque base64-encoded cursor
}
