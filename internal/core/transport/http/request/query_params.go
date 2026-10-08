package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val64, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	val := int(val64)

	return &val, nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}

func GetVersionQueryParam(r *http.Request) (*int, error) {
	const versionQueryParamKey = "version"

	version, err := GetIntQueryParam(r, versionQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'version' query param: %w", err)
	}

	if version != nil && *version < 1 {
		return nil, fmt.Errorf(
			"'version' must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return version, nil
}

func GetDateQueryParam(r *http.Request, key string, location *time.Location) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.ParseInLocation(layout, param, location)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}
