package utils

import (
	"fmt"
	productfeat "go-marketplace/internal/product"
	"net/http"
	"strconv"
)

func GetQueryParam(r *http.Request, key string) (int, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("param %s is not integer: %v: %w", valueStr, err, productfeat.ErrInvalidArgument)
	}

	return value, nil
}
