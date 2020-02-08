// +build ziplinegen

package services

import (
	"log"
	"time"

	"encoding/json"
	"net/http"

	"github.com/bilal-bhatti/zipline"
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(Authentication)

	mux.Post("/contacts", zipline.Post(ContactsService.Create))

	mux.Get("/contacts/{id}", zipline.Get(ContactsService.GetOne))
	mux.Get("/contacts/{month}-{day}-{year}", zipline.Get(ContactsService.GetByDate))
	mux.Post("/contacts/{id}", zipline.Post(ContactsService.Update))
	mux.Put("/contacts/{id}", zipline.Put(ContactsService.Replace))

	mux.Delete("/things/{id}", zipline.Delete(ThingsService.Delete))
	mux.Post("/echo", zipline.Post(Echo))

	return mux
}

type ZiplineTemplate struct{
	ReturnResponseAndError func() (interface{}, error)
	ReturnError func() error
}

func (z ZiplineTemplate) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing post request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()
		
		response, err := z.ReturnResponseAndError()
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func (z ZiplineTemplate) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing get request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		response, err := z.ReturnResponseAndError()
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func (z ZiplineTemplate) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not

		log.Println("Processing delete request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		err = z.ReturnError()

		if err != nil {
			// write error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (z ZiplineTemplate) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error // why not
		
		log.Println("Processing put request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		response, err := z.ReturnResponseAndError()
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}
