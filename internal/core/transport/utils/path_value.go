package utils

import (
	"fmt"
	productfeat "go-marketplace/internal/product"
	"net/http"
	"strconv"
)

func GetPathValue(r *http.Request, key string) (int, error) {
	valStr := r.PathValue(key)

	if valStr == "" {
		return 0, fmt.Errorf("no value by key %v: %w", key, productfeat.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value by key %v: %w: %w", key, err, productfeat.ErrInvalidArgument)
	}

	return val, nil
}
