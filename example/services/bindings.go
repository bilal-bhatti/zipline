// +build ziplinegen

package services

import (
	"github.com/bilal-bhatti/zipline"

	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(Authentication)
	
	mux.Post("/contacts", zipline.Post(ContactsService.Create))

	mux.Get("/contacts/{id}", zipline.Get(ContactsService.GetOne))
	mux.Get("/{month}-{day}-{year}", zipline.Get(ContactsService.GetByDate))
	mux.Post("/contacts/{id}", zipline.Post(ContactsService.Update))

	mux.Post("/echo", zipline.Post(Echo))

	return mux
}
