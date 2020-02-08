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

	mux.Post("/echo", zipline.Post(Echo))

	return mux
}

type ZiplineTemplate struct{}

func (z ZiplineTemplate) ReturnResponseAndError() (interface{}, error) {
	panic("Hi there! Whatcha doin?")
}

func (z ZiplineTemplate) NoReturn() {
	panic("Hi there! Whatcha doin?")
}

func (z ZiplineTemplate) Post() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var err error // temporary
		log.Println("Processing post request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()
		
		response, err := z.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


		err = json.NewEncoder(responseWriter).Encode(response)
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func (z ZiplineTemplate) Get() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var err error // temporary
		log.Println("Processing get request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		response, err := z.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(responseWriter).Encode(response)
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func (z ZiplineTemplate) Delete() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Println("Processing delete request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		z.NoReturn()

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func (z ZiplineTemplate) Put() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Println("Processing put request")
		startTime := time.Now()

		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %d to process request\n", duration)
		}()

		response, err := z.ReturnResponseAndError()
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(responseWriter).Encode(response)
		if err != nil {
			// write error response
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}
