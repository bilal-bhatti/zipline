//+build ziplinegen

package web

import (
	"log"
	"strconv"
	"time"

	"encoding/json"
	"net/http"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/services"
	"github.com/go-chi/chi"
)

// NewRouter returns a router configured with endpoints and handlers.
func NewRouter(env *connectors.Env) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(services.Authentication)

	mux.Post("/contacts", z.Post(services.ContactsService.Create, env, z.Resolve, z.Body))

	mux.Get("/contacts/{id}", z.Get(services.ContactsService.GetOne, z.Resolve, z.Path))
	mux.Get("/contacts/{month}-{day}-{year}", z.Get(services.ContactsService.GetByDate, z.Resolve, z.Path, z.Path, z.Path))
	mux.Post("/contacts/{id}", z.Post(services.ContactsService.Update, env, z.Resolve, z.Path, z.Body))
	mux.Put("/contacts/{id}", z.Put(new(services.ContactsService).Replace, z.Resolve, z.Path, z.Body))
	mux.Delete("/contacts", z.Delete(services.ContactsService.DeleteBulk, z.Resolve, z.Query))

	mux.Post("/things", z.Post(ThingsService.Create, z.Resolve, z.Body))
	mux.Get("/things/{category}", z.Get(ThingsService.GetByCategoryAndQuery, z.Resolve, z.Path, z.Query))
	mux.Get("/things", z.Get(ThingsService.GetByDateRange, z.Resolve, z.Query, z.Query))
	mux.Delete("/things/{id}", z.Delete(new(ThingsService).Delete, z.Path))

	mux.Get("/echo/{input}", z.Get(Echo, env, z.Resolve, z.Path))

	mux.Post("/doodads", z.Post(services.DoodadsService.Create, env, z.Resolve, z.Resolve, z.Body))

	mux.Post("/ping", z.Post(services.Ping, env, z.Resolve, z.Resolve, z.Body))

	return mux
}

var z ZiplineTemplate

type ZiplineTemplate struct {
	ReturnResponseAndError func() (interface{}, error)
	ReturnError            func() error
	Resolve                func() (ZiplineTemplate, error)
	DevNull                func(i ...interface{})
}

func (z ZiplineTemplate) Path(kind string, w http.ResponseWriter, r *http.Request) {
	switch kind {
	case "string":
		name := chi.URLParam(r, "name")
		z.DevNull(name)
	case "int":
		name, err := strconv.Atoi(chi.URLParam(r, "name"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		z.DevNull(name)
	}
}

func (z ZiplineTemplate) Query(kind string, w http.ResponseWriter, r *http.Request) {
	switch kind {
	case "string":
		name := r.URL.Query().Get("name")
		z.DevNull(name)
	case "[]string":
		name := r.URL.Query()["name"]
		z.DevNull(name)
	case "int":
		name, err := strconv.Atoi(r.URL.Query().Get("name"))
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		z.DevNull(name)
	case "*time.Time":
		name, err := ParseTime(r.URL.Query().Get("name"), "format")
		if err != nil {
			// invalid request error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		z.DevNull(name)
	}
}

func (z ZiplineTemplate) Body(w http.ResponseWriter, r *http.Request) {
	var err error
	name := &ZiplineTemplate{}
	err = json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		// invalid request error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (z ZiplineTemplate) Post(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

func (z ZiplineTemplate) Get(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

func (z ZiplineTemplate) Delete(i interface{}, params ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = handler.ReturnError()

		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (z ZiplineTemplate) Put(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
