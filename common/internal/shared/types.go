package shared

type Cursor interface {
}

type BatchRequest struct {
	Cursor       Cursor
	MaxBatchSize int
}
