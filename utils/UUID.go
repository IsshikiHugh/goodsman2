package utils

import "github.com/gofrs/uuid"

func GenerateUID() (string, error) {
	ul, err := uuid.NewV4()
	return ul.String(), err
}
