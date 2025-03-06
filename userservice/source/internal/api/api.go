package api

import "social/userservice/internal/storage"

type Api struct {
	storage storage.CommonRepository
}

func NewApi(storage storage.CommonRepository) *Api {
	return &Api{
		storage: storage,
	}
}
