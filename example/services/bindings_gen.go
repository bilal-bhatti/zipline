// Code generated by Zipline. DO NOT EDIT.

// go:generate zipline
// +build !ziplinegen

package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bilal-bhatti/zipline/example/models"
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(Authentication)
	mux.Post("/contacts", ContactsServiceCreateHandlerFunc())
	mux.Get("/contacts/{id}", ContactsServiceGetOneHandlerFunc())
	mux.Get("/contacts/{month}-{day}-{year}", ContactsServiceGetByDateHandlerFunc())
	mux.Post("/contacts/{id}", ContactsServiceUpdateHandlerFunc())
	mux.Put("/contacts/{id}", ContactsServiceReplaceHandlerFunc())
	mux.Delete("/things/{id}", ThingsServiceDeleteHandlerFunc())
	mux.Post("/echo", EchoHandlerFunc())
	return mux
}

// ContactsServiceCreateHandlerFunc handles requests to:
// path  : /contacts
// method: Post
func ContactsServiceCreateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing post request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// extract json body and marshall contactRequest
		contactRequest := &models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(contactRequest)
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// initialize application handler
		contactsService := InitContactsService()
		// execute application handler
		response, err := contactsService.Create(ctx, contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// ContactsServiceGetOneHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: Get
func ContactsServiceGetOneHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing get request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// parse path parameter id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// initialize application handler
		contactsService := InitContactsService()
		// execute application handler
		response, err := contactsService.GetOne(ctx, id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// ContactsServiceGetByDateHandlerFunc handles requests to:
// path  : /contacts/{month}-{day}-{year}
// method: Get
func ContactsServiceGetByDateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing get request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// resolve df dependency through a provider function
		df := ProvideDateFilter(r)

		// initialize application handler
		contactsService := InitContactsService()
		// execute application handler
		response, err := contactsService.GetByDate(ctx, df)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// ContactsServiceUpdateHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: Post
func ContactsServiceUpdateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing post request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// parse path parameter id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// extract json body and marshall contactRequest
		contactRequest := models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(&contactRequest)
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// initialize application handler
		contactsService := InitContactsService()
		// execute application handler
		response, err := contactsService.Update(ctx, id, contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// ContactsServiceReplaceHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: Put
func ContactsServiceReplaceHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing put request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// parse path parameter id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// extract json body and marshall contactRequest
		contactRequest := models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(&contactRequest)
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// initialize application handler
		contactsService := InitContactsService()
		// execute application handler
		response, err := contactsService.Replace(ctx, id, contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

// ThingsServiceDeleteHandlerFunc handles requests to:
// path  : /things/{id}
// method: Delete
func ThingsServiceDeleteHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing delete request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// parse path parameter id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// initialize application handler
		thingsService := InitThingsService()
		// execute application handler
		err = thingsService.Delete(ctx, id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// EchoHandlerFunc handles requests to:
// path  : /echo
// method: Post
func EchoHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing post request")
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		// resolve ctx dependency through a provider function
		ctx := ProvideContext(r)

		// extract json body and marshall echoRequest
		echoRequest := EchoRequest{}
		err = json.NewDecoder(r.Body).Decode(&echoRequest)
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := Echo(ctx, echoRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}
