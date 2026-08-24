package security

import (
	"crypto/md5" //nolint:revive,nolintlint
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"

	"github.com/elsyahtech/go-rest/config"
)

func HashPassword(rawPassword string) (hashedPassword string, troubleShootMsg string, httpCode int, err error) {
	cfg := config.GlobalConfig

	switch cfg.Database.PasswordEncryptionType {
	case "MD5":
		hash := md5.Sum([]byte(rawPassword))
		hashedPassword = hex.EncodeToString(hash[:])

	case "BCRYPT":
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if err != nil {
			troubleShootMsg = "failed to generate bcrypt hash. " +
				"Check system resource limits or CPU availability for cryptographic operations."
			httpCode = http.StatusInternalServerError

			return "", troubleShootMsg, httpCode, fmt.Errorf("hashing password using bcrypt error: %w", err)
		}

		hashedPassword = string(hashBytes)

	case "SHA-512":
		hash := sha512.Sum512([]byte(rawPassword))
		hashedPassword = hex.EncodeToString(hash[:])

	case "ARGON":
		saltByte := []byte(cfg.Database.Salt)
		timeCost := uint32(1)
		memoryCost := uint32(64 * 1024) // 64 MB
		threads := uint8(4)
		keyLen := uint32(32)
		hashBytes := argon2.IDKey([]byte(rawPassword), saltByte, timeCost, memoryCost, threads, keyLen)

		hashedPassword = hex.EncodeToString(hashBytes)

	default:
		troubleShootMsg = "unsupported password encryption type configured in app config. " +
			"Check 'PasswordEncryptionType' in ./app/config/database value (supported: MD5, BCRYPT, SHA-512, ARGON)."
		httpCode = http.StatusInternalServerError

		return "", troubleShootMsg, httpCode, errors.New("unsupported password encryption type")
	}

	httpCode = http.StatusOK

	return hashedPassword, troubleShootMsg, httpCode, nil
}

func PasswordValidation(passwordFromDB string, plainTextPassword string) (success bool, troubleShootMsg string, httpCode int, err error) {
	cfg := config.GlobalConfig

	switch cfg.Database.PasswordEncryptionType {
	case "MD5":
		hash := md5.Sum([]byte(plainTextPassword))
		hashedInput := hex.EncodeToString(hash[:])

		return hashedInput == passwordFromDB, troubleShootMsg, http.StatusOK, nil
	case "BCRYPT":
		err := bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(plainTextPassword))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return false, troubleShootMsg, http.StatusOK, nil
			}

			troubleShootMsg = "unexpected error occurred during bcrypt password comparison. " +
				"Check if the stored hash format in the database is corrupted."

			return false, troubleShootMsg, http.StatusInternalServerError, fmt.Errorf("bcrypt comparison error: %w", err)
		}

		return true, troubleShootMsg, http.StatusOK, nil
	case "SHA-512":
		hash := sha512.Sum512([]byte(plainTextPassword))
		hashedInput := hex.EncodeToString(hash[:])

		return hashedInput == passwordFromDB, troubleShootMsg, http.StatusOK, nil
	case "ARGON":
		saltByte := []byte(cfg.Database.Salt)
		timeCost := uint32(1)
		memoryCost := uint32(64 * 1024) // 64 MB
		threads := uint8(4)
		keyLen := uint32(32)

		hashBytes := argon2.IDKey([]byte(plainTextPassword), saltByte, timeCost, memoryCost, threads, keyLen)
		hashedInput := hex.EncodeToString(hashBytes)

		return hashedInput == passwordFromDB, troubleShootMsg, http.StatusOK, nil
	default:
		troubleShootMsg = "unsupported password encryption type configured in app config. " +
			"Check 'PasswordEncryptionType' in ./app/config/database value (supported: MD5, BCRYPT, SHA-512, ARGON)."

		return false, troubleShootMsg, http.StatusInternalServerError, errors.New("unsupported password encryption type")
	}
}
