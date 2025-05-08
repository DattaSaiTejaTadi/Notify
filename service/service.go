package service

import "Notify/store"

type service struct {
	Store store.Store
}

func New(store store.Store) *service {
	return &service{
		Store: store,
	}
}
