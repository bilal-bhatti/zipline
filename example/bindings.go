package example

// +build ziplinegen

import (
	"github.com/bilal-bhatti/zipline/example/services"
	zipline "github.com/bilal-bhatti/zipline/pkg"

	"github.com/go-chi/chi"
)

// NewRouter ....
func NewRouter() *chi.Mux {
	mux := chi.NewRouter()

	contacts := &services.ContactsService{}
	mux.Post("/contacts", zipline.Post(contacts.Create))
	mux.Get("/contacts/{id}", zipline.Get(contacts.GetOne))
	mux.Post("/contacts/{id}", zipline.Post(contacts.Update))
	mux.Post("/echo", zipline.Post(services.Echo))

	return mux
}
