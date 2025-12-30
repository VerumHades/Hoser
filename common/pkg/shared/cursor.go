package shared

type BatchRequest[CursorType any] struct {
	Cursor       CursorType
	MaxBatchSize int
}
