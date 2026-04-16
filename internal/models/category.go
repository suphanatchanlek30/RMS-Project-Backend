package models

import "time"

type CreateCategoryRequest struct {
	CategoryName string `json:"categoryName"`
	Description  string `json:"description"`
}

type CategoryResponse struct {
	CategoryID   int       `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Description  *string   `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
}
