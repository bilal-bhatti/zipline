// Code generated by Zipline. DO NOT EDIT.

// go:generate zipline
// +build !ziplinegen

package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bilal-bhatti/zipline/example/models"
	"github.com/go-chi/chi"
)

type (
	ContactsServiceCreateType func(context.Context, *models.ContactRequest) (*models.ContactResponse, error)
	ContactsServiceGetOneType func(int) (*models.ContactResponse, error)
	ContactsServiceUpdateType func(context.Context, int, models.ContactRequest) (*models.ContactResponse, error)
	EchoType                  func(EchoRequest) (EchoResponse, error)
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Post("/contacts", ContactsServiceCreateHandlerFunc(InitContactsService().Create))
	mux.Get("/contacts/{id}", ContactsServiceGetOneHandlerFunc(InitContactsService().GetOne))
	mux.Post("/contacts/{id}", ContactsServiceUpdateHandlerFunc(InitContactsService().Update))
	mux.Post("/echo", EchoHandlerFunc(Echo))
	return mux
}

func ContactsServiceCreateHandlerFunc(funk ContactsServiceCreateType) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error // tempory fix

		ctx := ProvideContext(req)

		contactRequest := &models.ContactRequest{}
		err = json.NewDecoder(req.Body).Decode(contactRequest)
		if err != nil {
			// write error response
			// invalid request error
			panic(err)
		}

		result, err := funk(ctx, contactRequest)
		if err != nil {
			// write error response
			// internal error
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			// write error response
			panic(err)
		}
	}
}

func ContactsServiceGetOneHandlerFunc(funk ContactsServiceGetOneType) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error // tempory fix

		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			// invalid request error
			panic(err)
		}

		result, err := funk(id)
		if err != nil {
			// write error response
			// internal error
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			// write error response
			panic(err)
		}
	}
}

func ContactsServiceUpdateHandlerFunc(funk ContactsServiceUpdateType) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error // tempory fix

		ctx := ProvideContext(req)

		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			// invalid request error
			panic(err)
		}

		contactRequest := models.ContactRequest{}
		err = json.NewDecoder(req.Body).Decode(&contactRequest)
		if err != nil {
			// write error response
			// invalid request error
			panic(err)
		}

		result, err := funk(ctx, id, contactRequest)
		if err != nil {
			// write error response
			// internal error
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			// write error response
			panic(err)
		}
	}
}

func EchoHandlerFunc(funk EchoType) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error // tempory fix

		echoRequest := EchoRequest{}
		err = json.NewDecoder(req.Body).Decode(&echoRequest)
		if err != nil {
			// write error response
			// invalid request error
			panic(err)
		}

		result, err := funk(echoRequest)
		if err != nil {
			// write error response
			// internal error
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			// write error response
			panic(err)
		}
	}
}
