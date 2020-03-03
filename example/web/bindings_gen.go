// Code generated by Zipline. DO NOT EDIT.

//+build !ziplinegen

package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bilal-bhatti/zipline/example/models"
	"github.com/bilal-bhatti/zipline/example/services"
	"github.com/go-chi/chi"
)

// NewRouter returns a router configured with endpoints and handlers.
func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(services.Authentication)

	mux.Post("/contacts", ContactsServiceCreateHandlerFunc())

	mux.Get("/contacts/{id}", ContactsServiceGetOneHandlerFunc())
	mux.Get("/contacts/{month}-{day}-{year}", ContactsServiceGetByDateHandlerFunc())
	mux.Post("/contacts/{id}", ContactsServiceUpdateHandlerFunc())
	mux.Put("/contacts/{id}", ContactsServiceReplaceHandlerFunc())

	mux.Get("/things/{category}", ThingsServiceGetByCategoryAndQueryHandlerFunc())
	mux.Get("/things", ThingsServiceGetByDateRangeHandlerFunc())
	mux.Delete("/things/{id}", ThingsServiceDeleteHandlerFunc())

	mux.Post("/echo", EchoHandlerFunc())

	return mux
}

// ContactsServiceCreateHandlerFunc handles requests to:
// path  : /contacts
// method: post
// Create a new contact request entity.
func ContactsServiceCreateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := services.InitContactsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [contactRequest] with [Body] template
		contactRequest := &models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := handler.Create(ctx, contactRequest)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ContactsServiceGetOneHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: get
// GetOne contact by id
func ContactsServiceGetOneHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := services.InitContactsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [id] with [Path] template
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := handler.GetOne(ctx, id)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ContactsServiceGetByDateHandlerFunc handles requests to:
// path  : /contacts/{month}-{day}-{year}
// method: get
// Get contacts list by date
func ContactsServiceGetByDateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := services.InitContactsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [month] with [Path] template
		month := chi.URLParam(r, "month")

		// resolve parameter [day] with [Path] template
		day := chi.URLParam(r, "day")

		// resolve parameter [year] with [Path] template
		year := chi.URLParam(r, "year")

		// execute application handler
		response, err := handler.GetByDate(ctx, month, day, year)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ContactsServiceUpdateHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: post
// Update a contact entity with provided data.
func ContactsServiceUpdateHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := services.InitContactsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [id] with [Path] template
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// resolve parameter [contactRequest] with [Body] template
		contactRequest := models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(&contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := handler.Update(ctx, id, contactRequest)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ContactsServiceReplaceHandlerFunc handles requests to:
// path  : /contacts/{id}
// method: put
// Replace a contact entity completely.
func ContactsServiceReplaceHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := services.InitContactsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [id] with [Path] template
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// resolve parameter [contactRequest] with [Body] template
		contactRequest := models.ContactRequest{}
		err = json.NewDecoder(r.Body).Decode(&contactRequest)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := handler.Replace(ctx, id, contactRequest)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ThingsServiceGetByCategoryAndQueryHandlerFunc handles requests to:
// path  : /things/{category}
// method: get
// Get things by category and search query
func ThingsServiceGetByCategoryAndQueryHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := InitThingsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [category] with [Path] template
		category := chi.URLParam(r, "category")

		// resolve parameter [q] with [Query] template
		q := r.URL.Query().Get("q")

		// execute application handler
		response, err := handler.GetByCategoryAndQuery(ctx, category, q)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ThingsServiceGetByDateRangeHandlerFunc handles requests to:
// path  : /things
// method: get
// Get things by date range
//
// @from `format:"date-time,2006-01-02"` date should be in Go time format
// @to   `format:"date-time,2006-01-02"` date should be in Go time format
func ThingsServiceGetByDateRangeHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := InitThingsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [from] with [Query] template
		from, err := ParseTime(r.URL.Query().Get("from"), "date-time,2006-01-02")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// resolve parameter [to] with [Query] template
		to, err := ParseTime(r.URL.Query().Get("to"), "date-time,2006-01-02")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		response, err := handler.GetByDateRange(ctx, from, to)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

// ThingsServiceDeleteHandlerFunc handles requests to:
// path  : /things/{id}
// method: delete
// Delete thing by id
func ThingsServiceDeleteHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		// initialize application handler
		handler, err := InitThingsService()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [id] with [Path] template
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// execute application handler
		err = handler.Delete(id)
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// EchoHandlerFunc handles requests to:
// path  : /echo
// method: post
// Echo returns body with 'i's replaced with 'o's
func EchoHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// resolve parameter [ctx] through a provider
		ctx := services.ProvideContext(r)

		// resolve parameter [echoRequest] with [Body] template
		echoRequest := EchoRequest{}
		err = json.NewDecoder(r.Body).Decode(&echoRequest)
		if err != nil {
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
