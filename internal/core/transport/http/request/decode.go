package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = newRequestValidator()

type validatable interface {
	Validate() error
}

func newRequestValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			return field.Name
		}

		return name
	})

	return v
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
			"decode json: %s: %w",
			jsonDecodeErrorText(err),
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
			"request validation: %s: %w",
			validationErrorText(err),
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func jsonDecodeErrorText(err error) string {
	if errors.Is(err, io.EOF) {
		return "request body is empty"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		if typeErr.Field == "" {
			return fmt.Sprintf("unexpected JSON %s", typeErr.Value)
		}

		return fmt.Sprintf("field '%s' has invalid type: got JSON %s", typeErr.Field, typeErr.Value)
	}

	return strings.TrimPrefix(err.Error(), "json: ")
}

func validationErrorText(err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err.Error()
	}

	texts := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		text := fmt.Sprintf("field '%s' failed on '%s'", fieldErr.Field(), fieldErr.Tag())
		if fieldErr.Param() != "" {
			text += "=" + fieldErr.Param()
		}

		texts = append(texts, text)
	}

	return strings.Join(texts, "; ")
}
