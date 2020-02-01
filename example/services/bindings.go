// +build ziplinegen

package services

import (
	zipline "github.com/bilal-bhatti/zipline/pkg"

	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()

	contacts := &ContactsService{}
	mux.Post("/contacts", zipline.Post(contacts.Create))
	mux.Get("/contacts/{id}", zipline.Get(contacts.GetOne))
	mux.Post("/contacts/{id}", zipline.Post(contacts.Update))
	mux.Post("/echo", zipline.Post(Echo))

	return mux
}
