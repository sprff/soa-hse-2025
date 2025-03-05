package models

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](x T) *T {
	return &x
}
func TestUser_Vaildate(t *testing.T) {
	t.Run("All invalid", func(t *testing.T) {

		user := User{
			Name:        ptr(strings.Repeat("long", 10)),
			Surname:     ptr(strings.Repeat("long", 10)),
			DateOfBirth: ptr(time.Now()),
			Email:       "unknown",
			Phone:       ptr("unknown"),
		}
		err := user.Validate()
		reasons := []string{
			"Name should be 30 charecters or less",
			"Surname should be 30 charecters or less",
			"You have to be at least 14 years old",
			"Invalid email",
			"Invalid phone",
		}

		assert.Equal(t, ErrUserInvalid{reasons: reasons}, err)
		assert.True(t, errors.Is(err, ErrUserInvalid{}))
		assert.True(t, errors.Is(err, ErrUserInvalid{reasons: reasons}))
	})
	t.Run("Ok", func(t *testing.T) {

		user := User{
			Name:        ptr("Maxim"),
			Surname:     ptr("Saprykin"),
			DateOfBirth: ptr(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			Email:       "mail@sprff.ru",
			Phone:       ptr("+78005553535"),
		}
		err := user.Validate()
		assert.NoError(t, err)
	})
}

func TestUser_HidePassword(t *testing.T) {
	user := User{
		Login:    "sprff",
		Password: "aboba",
	}

	user.HidePassword()
	assert.Equal(t, "", user.Login)
}
