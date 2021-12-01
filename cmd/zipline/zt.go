package main

var zTemplate = `
//go:build ziplinegen
// +build ziplinegen

package {{.PackageName}}

import (
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"encoding/json"
	"net/http"

	"github.com/bilal-bhatti/zipline/example/render"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// RegisterEndpoints accepts a github.com/go-chi/v5 router instance and
// returns a router configured with endpoints and handlers
func RegisterEndpoints(mux *chi.Mux) *chi.Mux {
	// register your routes here

	return mux
}

var z ZiplineTemplate

// ZiplineTemplate is the code generation template for zipline cli
// It is required, without it the tool does nothing
type ZiplineTemplate struct {
	// marker that the func returns a type and an error, so we have var handles in the template
	ReturnResponseAndError func() (interface{}, error)

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
		//render.NewBadRequestError("")
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
		var err error // why not

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
		var err error // why not

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
		var err error // why not

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
		var err error // why not

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

`
