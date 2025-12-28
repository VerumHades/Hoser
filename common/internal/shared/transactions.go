package shared

type Transaction interface {
	Rollback() error
	Commit() error
}
