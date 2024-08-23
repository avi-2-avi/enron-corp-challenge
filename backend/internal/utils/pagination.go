package utils

import "backend/internal/models"

type Pagination struct {
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalElements int            `json:"total_elements"`
	TotalPages    int            `json:"total_pages"`
	Emails        []models.Email `json:"emails"`
}

func Paginate(emails []models.Email, totalElements, page, size int) Pagination {
	totalPages := (totalElements + size - 1) / size

	return Pagination{
		Page:          page,
		Size:          size,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Emails:        emails,
	}
}
