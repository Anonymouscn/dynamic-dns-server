package error

var RecordNotFoundErr = &RecordNotFound{}

type RecordNotFound struct{}

func (RecordNotFound) Error() string {
	return "record not found"
}
