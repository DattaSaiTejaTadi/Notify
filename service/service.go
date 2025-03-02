package service

import "Notify/store"

type service struct {
	store store.Store
}

func New(store store.Store) *service {
	return &service{
		store: store,
	}
}
