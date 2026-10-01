package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return fmt.Errorf(
				"request body exceeds %d bytes: %w",
				maxBytesErr.Limit,
				core_errors.ErrInvalidArgument,
			)
		}

		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf(
			"request body must contain a single JSON object: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	var (
		err error
	)

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
