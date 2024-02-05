//go:build ziplinegen
// +build ziplinegen

package web

import (
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"encoding/json"
	"net/http"

	"github.com/bilal-bhatti/zipline/example/connectors"
	"github.com/bilal-bhatti/zipline/example/render"
	"github.com/bilal-bhatti/zipline/example/services"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// NewRouter returns a router configured with endpoints and handlers.
//
// @swagger                   2.0
// @info.title                Example OpenAPI Version 2 Specification
// @info.version              1.0.0
// @info.description          Example OpenAPI Version 2 Specification
// info.termsOfService       http://swagger.io/terms/
// info.contact.name         API Support
// info.contact.url          http://www.swagger.io/support
// info.contact.email        support@swagger.io
// info.license              (name: Apache 2.0, url: http://www.apache.org/licenses/LICENSE-2.0.html)
// @schemes                   [http, https]
// @host                      api.example.com
// @basePath                  /api
// @consumes                  application/json
// @produces                  application/json
// produces                  application/text
// securityDefinitions.basic  BasicAuth
// externalDocs.description  OpenAPI
// externalDocs.url          https://swagger.io/resources/open-api/
//
// summary          Get a list of contacts by ids
// description      Get a list of contacts by ids
// tags              contacts
// produces         application/json
// parameters        (name:ids, description: list of contact ids, required:true)
// parameters        (name:foo, description: foo description, required:false)
// responses.400     {models.ErrorResponse}
// responses.404     {models.ErrorResponse}
// responses.default {models.ErrorResponse}

// @info.title                Zipline Example Swagger API
// @info.version              1.1
// @info.description          This is a sample Zipline generated server.
// @info.termsOfService       http://swagger.io/terms/
// @info.contact.name         API Support
// @info.contact.url          http://www.swagger.io/support
// @info.contact.email        support@swagger.io
// @info.license              (name: Apache 2.0, url: http://www.apache.org/licenses/LICENSE-2.0.html)
// @schemes                   [http, https, "s,ftp"]
// @host                      zipline.example.com
// @basePath                  /api
// @consumes                  application/json
// @produces                  application/json
// @produces                  application/text
// @scurityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// @summary           Get a list of contacts by ids
// @description       Get a list of contacts by ids
// @tags              contacts
// @produces          application/json
// @parameters        (name:ids, description: list of contact ids, required:true)
// @parameters        (name:foo, description: foo description, required:false)
// @responses.400     {models.ErrorResponse}
// @responses.404     {models.ErrorResponse}
// @responses.default {models.ErrorResponse}
func NewRouter(env *connectors.Env) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(services.Authentication)

	mux.Post("/contacts", z.Post(services.ContactsService.Create, env, z.Resolve, z.Body))

	mux.Get("/contacts/{id}", z.Get(services.ContactsService.GetOne, z.Resolve, z.Path))
	mux.Get("/contacts", z.Get(services.ContactsService.GetBunch, z.Resolve, z.Query))
	mux.Get("/contacts/{month}-{day}-{year}", z.Get(services.ContactsService.GetByDate, z.Resolve, z.Path, z.Path, z.Path))
	mux.Post("/contacts/{id}", z.Post(services.ContactsService.Update, env, z.Resolve, z.Path, z.Body))
	mux.Put("/contacts/{id}", z.Put(new(services.ContactsService).Replace, z.Resolve, z.Path, z.Body))
	mux.Delete("/contacts", z.Delete(services.ContactsService.DeleteBulk, z.Resolve, z.Query))
	//mux.Get("/contacts/redirect", z.Redirect(services.ContactsService.Redirect, z.Resolve, z.Query))

	mux.Post("/things", z.Post(ThingsService.Create, z.Resolve, z.Body))
	mux.Get("/things/{category}", z.Get(ThingsService.GetByCategoryAndQuery, z.Resolve, z.Path, z.Query))
	mux.Get("/things", z.Get(ThingsService.GetByDateRange, z.Resolve, z.Query, z.Query))
	mux.Delete("/things/{id}", z.Delete(new(ThingsService).Delete, z.Path))

	mux.Get("/echo/{input}", z.Get(Echo, env, z.Resolve, z.Path))

	mux.Post("/doodads", z.Post(services.DoodadsService.Create, env, z.Resolve, z.Resolve, z.Resolve, z.Body))

	mux.Post("/ping", z.Post(services.Ping, env, z.Resolve, z.Resolve, z.Body))

	return mux
}

var z ZiplineTemplate

// ZiplineTemplate is the code generation template for zipline cli
// It is required, without it the tool does nothing
type ZiplineTemplate struct {
	// marker that the func returns a type and an error, so we have var handles in the template
	ReturnResponseAndError func() (interface{}, error)

	// marker that the func returns a string and an error, so we have var handles in the template
	ReturnResponseStringAndError func() (string, error)

	// marker that the func has a single return of type error, so we have a error handle in the template
	ReturnError func() error

	// directive for zipline to go find a way to resolve the type
	Resolve func() (ZiplineTemplate, error)

	// used in the templates to stop go from complaining about unused vars omitted from output
	DevNull func(i ...interface{})
}

// Path is a template applied to path parameters
// Must only contain a single switch statement
// Each expected path parameter type needs a case clause
// name is a reserved keyword used for AST rewriting
// NOTE: best to have a single expression to find and convert the parameter
func (z ZiplineTemplate) Path(kind string, w http.ResponseWriter, r *http.Request) {
	switch kind {
	case "string":
		name := chi.URLParam(r, "name")
		z.DevNull(name)
	case "int":
		name, err := strconv.Atoi(chi.URLParam(r, "name"))
		if err != nil {
			render.Error(w, render.WrapBadRequestError(err, "failed to parse parameter"))
			return
		}
		z.DevNull(name)
	}
}

// Query is a template applied to query parameters
// Must only contain a single switch statement
// Each expected path parameter type needs a case clause
// name is a reserved keyword used for AST rewriting
// NOTE: best to have a single expression to find and convert the parameter
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
			render.Error(w, render.WrapBadRequestError(err, "failed to parse parameter"))
			return
		}
		z.DevNull(name)
	case "[]int64":
		name, err := ParseInt64(r.URL.Query()["name"])
		if err != nil {
			render.Error(w, render.WrapBadRequestError(err, "failed to parse parameter"))
			return
		}
		z.DevNull(name)
	case "*time.Time":
		name, err := ParseTime(r.URL.Query().Get("name"), "format")
		if err != nil {
			render.Error(w, render.WrapBadRequestError(err, "failed to parse parameter"))
			return
		}
		z.DevNull(name)
	}
}

// Body is a template applied to request body
// name := ZiplineTemplate{} line is replaced with the an instance of the desired struct
// Remaining code is rewritten with name substituted with the desired type and name
func (z ZiplineTemplate) Body(w http.ResponseWriter, r *http.Request) {
	var err error

	defer io.Copy(ioutil.Discard, r.Body)

	name := ZiplineTemplate{}
	err = json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		render.Error(w, render.WrapBadRequestError(err, "failed to parse request body"))
		return
	}
}

// Post template is expected to be applied to HTTP POST requests.
// It resolves HTTP parameters, request body, required application service handler
// and invokes specified method
// NOTE: All HTTP related marshalling/unmarshalling should take place here
// All business related input validation and processing logic should be in the service
func (z ZiplineTemplate) Post(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			render.Error(w, errors.Wrap(err, "failed to resolve application handler"))
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			render.Error(w, errors.Wrap(err, "application handler failed"))
			return
		}

		render.Response(w, response)
	}
}

// Get template is expected to be applied to HTTP GET requests.
// It resolves HTTP parameters, required application service handler
// and invokes specified method
// NOTE: All HTTP related marshalling/unmarshalling should take place here
// All business related input validation and processing logic should be in the service
func (z ZiplineTemplate) Get(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			render.Error(w, errors.Wrap(err, "failed to resolve application handler"))
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			render.Error(w, errors.Wrap(err, "application handler failed"))
			return
		}

		render.Response(w, response)
	}
}

// Delete template is expected to be applied to HTTP DELETE requests.
// It resolves HTTP parameters, required application service handler
// and invokes specified method
// NOTE: All HTTP related marshalling/unmarshalling should take place here
// All business related input validation and processing logic should be in the service
func (z ZiplineTemplate) Delete(i interface{}, params ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			render.Error(w, errors.Wrap(err, "failed to resolve application handler"))
			return
		}

		err = handler.ReturnError()
		if err != nil {
			render.Error(w, errors.Wrap(err, "application handler failed"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Put template is expected to be applied to HTTP PUT requests.
// It resolves HTTP parameters, request body, required application service handler
// and invokes specified method
// NOTE: All HTTP related marshalling/unmarshalling should take place here
// All business related input validation and processing logic should be in the service
func (z ZiplineTemplate) Put(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			render.Error(w, errors.Wrap(err, "failed to resolve application handler"))
			return
		}

		response, err := handler.ReturnResponseAndError()
		if err != nil {
			render.Error(w, errors.Wrap(err, "application handler failed"))
			return
		}

		render.Response(w, response)
	}
}

// Redirect template is expected to be applied to HTTP GET requests that respond
// with a redirect directive.
// It resolves HTTP parameters, required application service handler
// and invokes specified method
// NOTE: All HTTP related marshalling/unmarshalling should take place here
// All business related input validation and processing logic should be in the service
func (z ZiplineTemplate) Redirect(i interface{}, p ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime)
			log.Printf("It took %s to process request\n", duration.String())
		}()

		handler, err := z.Resolve()
		if err != nil {
			render.Error(w, errors.Wrap(err, "failed to resolve application handler"))
			return
		}

		response, err := handler.ReturnResponseStringAndError()
		if err != nil {
			render.Error(w, errors.Wrap(err, "application handler failed"))
			return
		}

		http.Redirect(w, r, response, http.StatusFound)
	}
}
