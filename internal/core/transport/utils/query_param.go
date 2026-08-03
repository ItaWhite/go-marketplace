package utils

import (
	"fmt"
	"go-marketplace/internal/core/errors"
	"net/http"
	"strconv"
)

func GetQueryParam(r *http.Request, key string) (*int, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, fmt.Errorf("param %s=%s: %v: %w", key, valueStr, err, core_errors.ErrInvalidQueryParam)
	}

	return &value, nil
}
