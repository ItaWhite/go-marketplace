package utils

import (
	"fmt"
	"go-marketplace/internal/core/transport/errors"
	"net/http"
	"strconv"
)

func GetPathValue(r *http.Request, key string) (int, error) {
	valStr := r.PathValue(key)

	if valStr == "" {
		return 0, fmt.Errorf("no value by key %v: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value by key %v: %w: %w", key, err, core_errors.ErrInvalidArgument)
	}

	return val, nil
}
