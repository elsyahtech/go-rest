package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/go-rest/security"
)

func (*Helper) HashingPassword(rawPassword string) (hashedPassword, troubleShootMsg string, httpCode int, err error) {
	if rawPassword == "" {
		return "", troubleShootMsg, http.StatusOK, nil
	}

	hashedPassword, troubleShootMsg, httpCode, err = security.HashPassword(rawPassword)
	if err != nil {
		return hashedPassword, troubleShootMsg, httpCode, fmt.Errorf("%w", err)
	}

	return hashedPassword, troubleShootMsg, httpCode, nil
}

func (*Helper) ValidatingPassword(passwordFromDB string, plainTextPassword string) (success bool, troubleShootMsg string, httpCode int, err error) {
	if passwordFromDB == "" || plainTextPassword == "" {
		troubleShootMsg = "either the database password hash or the input plain text password is empty during validation. " +
			"Check data integrity from DB or request payload."

		return false, troubleShootMsg, http.StatusBadRequest, errors.New("password cannot be empty")
	}

	success, troubleShootMsg, httpCode, err = security.PasswordValidation(passwordFromDB, plainTextPassword)
	if err != nil {
		return success, troubleShootMsg, httpCode, fmt.Errorf("%w", err)
	}

	return success, troubleShootMsg, httpCode, nil
}
