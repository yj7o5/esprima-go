package esprimago

type ErrorHandler interface {
	getSource() *string
	isTolerant() bool
	RecordError(err interface{})
	Tolerate(err interface{})
	CreateError(index, line, column int, message string) error
	TolerateError(index, line, column int, message string)
}
