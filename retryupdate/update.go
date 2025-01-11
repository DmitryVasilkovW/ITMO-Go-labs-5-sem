//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var (
	authError     *kvapi.AuthError
	conflictError *kvapi.ConflictError
	oldVersion    uuid.UUID
	oldVal        *string = nil
	updateFnc     func(oldVal *string) (newVal string, err error)
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	resetFields()
	updateFnc = updateFn

	newVal, err := tryToUpDateValue(c, key)
	if err != nil {
		return err
	}

	newVersion := uuid.Must(uuid.NewV4())

	return tryToSetValue(c, key, newVal, newVersion)
}

func resetFields() {
	oldVal = nil
	oldVersion = uuid.UUID{}
}

func tryToSetValue(c kvapi.Client, key, newVal string, newVersion uuid.UUID) error {
	for {
		_, err := c.Set(&kvapi.SetRequest{Key: key, Value: newVal, OldVersion: oldVersion, NewVersion: newVersion})

		if err == nil || errors.As(err, &authError) {
			return err
		} else if errors.Is(err, kvapi.ErrKeyNotFound) {
			newVal, err = updateFnc(nil)
			oldVersion = uuid.UUID{}
			if err != nil {
				return err
			}
		} else if errors.As(err, &conflictError) && conflictError.ExpectedVersion == newVersion {
			return nil
		} else if errors.As(err, &conflictError) {
			return UpdateValue(c, key, updateFnc)
		}
	}
}

func tryToUpDateValue(c kvapi.Client, key string) (string, error) {
	err := tryToGetValue(c, key)
	if err != nil {
		return "", err
	}

	newVal, err := updateFnc(oldVal)
	if err != nil {
		return "", err
	}

	return newVal, nil
}

func tryToGetValue(c kvapi.Client, key string) error {
tryToGetValue:
	for {
		response, err := c.Get(&kvapi.GetRequest{Key: key})

		switch {
		case errors.Is(err, kvapi.ErrKeyNotFound):
			break tryToGetValue
		case err == nil:
			oldVal = &response.Value
			oldVersion = response.Version
			break tryToGetValue
		case errors.As(err, &authError):
			return err
		}
	}

	return nil
}
