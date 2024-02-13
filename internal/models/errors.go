package models

type ErrorDefinition struct {
	Message   string
	ErrorCode string
}

func (e *ErrorDefinition) Error() string {
	return e.Message
}
