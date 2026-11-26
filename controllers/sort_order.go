package controllers

import (
	"errors"
	"strconv"

	"api/portodata/config"
)

func nextSortOrder(model any) (int, error) {
	var maxSortOrder int
	if err := config.DB.Model(model).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error; err != nil {
		return 0, err
	}

	return maxSortOrder + 1, nil
}

func parseSortOrder(value string) (*int, error) {
	if value == "" {
		return nil, nil
	}

	sortOrder, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	if sortOrder < 1 {
		return nil, errors.New("sort_order harus lebih besar dari 0")
	}
	return &sortOrder, nil
}
