package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	gen, err := uuid.NewV7()
	if err != nil {
		panic("Can't generate uuid7")
	}
	return strings.Replace(gen.String(), "-", "", -1)
}
