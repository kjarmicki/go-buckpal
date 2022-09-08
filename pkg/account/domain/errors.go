package account_domain

type IllegalStateError struct {
	message string
}

func NewIllegalStateError(message string) *IllegalStateError {
	return &IllegalStateError{
		message: message,
	}
}

func (ise *IllegalStateError) Error() string {
	return ise.message
}
