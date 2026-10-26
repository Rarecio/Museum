package models

// PageResponse повторяет формат Spring Data Page<T>
type Page[T any] struct {
	Content          []T   `json:"content"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int64 `json:"totalPages"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int64 `json:"numberOfElements"`
	Number           int64 `json:"number"`
	//Size             int64 `json:"size"`
	//Sort             Sort  `json:"sort"`
	//Empty            bool  `json:"empty"`
}

func NewPage[T any](content []T, totalElements, totalPages, page int64) *Page[T] {
	return &Page[T]{
		Content:          content,
		TotalElements:    totalElements,
		TotalPages:       totalPages,
		First:            page == 0,
		Last:             page == totalPages,
		NumberOfElements: int64(len(content)),
		Number:           page,
	}
}

// Sort – состояние сортировки
type Sort struct {
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
	Empty    bool `json:"empty"`
}
