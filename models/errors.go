package models

const (
	ErrNotFound modelError = "models: resource not found."
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}
