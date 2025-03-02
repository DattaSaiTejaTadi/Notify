package handler

import "Notify/service"

type handler struct {
	service service.Service
}

func New(service service.Service) *handler {
	return &handler{
		service: service,
	}
}
