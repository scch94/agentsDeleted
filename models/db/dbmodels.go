package modeldb

type ModelsDb interface {
	Condition() string
	CanDeleted() bool
}
