// +build ziplinegen

package services

import (
	"github.com/bilal-bhatti/zipline"

	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()

	mux.Post("/contacts", zipline.Post(InitContactsService().Create))

	mux.Get("/contacts/{id}", zipline.Get(InitContactsService().GetOne))
	mux.Post("/contacts/{id}", zipline.Post(InitContactsService().Update))

	mux.Post("/echo", zipline.Post(Echo))

	return mux
}
