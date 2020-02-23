package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

func (s swagger) markdown() error {
	buf := newBuffer()

	buf.ws("# API Summary\n\n")

	buf.ws("```\n")
	buf.ws("Version:     %s\n", s.swag.SwaggerProps.Info.Version)
	buf.ws("Title:       %s\n", s.swag.SwaggerProps.Info.Title)
	buf.ws("Description: %s\n", s.swag.SwaggerProps.Info.Description)
	buf.ws("Host:        %s\n", s.swag.Host)
	buf.ws("BasePath:    %s\n", s.swag.BasePath)
	buf.ws("Consumes:    %s\n", s.swag.Consumes)
	buf.ws("Produces:    %s\n", s.swag.Produces)
	buf.ws("```\n\n")

	md := func(m, path string, op *spec.Operation) {
		buf.ws("<details>\n")
		buf.ws("<summary>%s: %s</summary>\n\n", path, m)

		params := func(op *spec.Operation) ([]spec.Parameter, []spec.Parameter, []spec.Parameter) {
			path, query, body := []spec.Parameter{}, []spec.Parameter{}, []spec.Parameter{}

			for _, p := range op.Parameters {
				if p.In == "path" {
					path = append(path, p)
				}
				if p.In == "query" {
					query = append(query, p)
				}
				if p.In == "body" {
					body = append(body, p)
				}
			}

			return path, query, body
		}

		if op.Parameters != nil {
			path, query, body := params(op)

			if len(path) > 0 {
				buf.ws("`path parameters`\n")
				for _, p := range path {
					if p.In == "path" {
						buf.ws("- name: `%s`, type: `%s`\n", p.Name, p.Type)
					}
				}
				buf.ws("\n")
			}

			if len(query) > 0 {
				buf.ws("`query parameters`\n")
				for _, p := range query {
					if p.In == "query" {
						buf.ws("- name: `%s`, type: `%s`\n", p.Name, p.Type)
					}
				}
				buf.ws("\n")
			}

			if len(body) > 0 {
				buf.ws("`body parameter`\n")
				for _, p := range body {
					if p.In == "body" {
						buf.ws("- name: `%s`, type: `%s`\n", p.Name, strings.Split(p.Schema.Ref.GetPointer().String(), "/")[2])
						if p.Schema != nil {
							def := s.swag.Definitions[strings.Split(p.Schema.Ref.GetPointer().String(), "/")[2]]

							obj(1, def, buf)
						}
					}
				}
			}
		}
		buf.ws("\n")

		buf.ws("`responses`\n")
		if op.Responses != nil {
			for code, r := range op.Responses.StatusCodeResponses {
				if r.ResponseProps.Schema != nil {
					var ref string
					if r.ResponseProps.Schema.Items != nil {
						ref = r.ResponseProps.Schema.Items.Schema.Ref.GetPointer().String()
					} else {
						ref = r.ResponseProps.Schema.Ref.GetPointer().String()
					}

					idx := strings.LastIndex(ref, "/")

					if idx > 0 {
						ref = strings.TrimPrefix(ref[idx:], "/")
					}

					def := s.swag.Definitions[ref]

					if r.ResponseProps.Schema.Items != nil {
						buf.ws("- code: `%d`, type: `[]%s`\n", code, ref)
						obj(1, *def.Items.Schema, buf)
					} else {
						buf.ws("- code: `%d`, type: `%s`\n", code, ref)
						obj(1, def, buf)
					}
				}
			}

			if op.Responses.Default != nil {
				ref := op.Responses.Default.ResponseProps.Schema.Ref.GetPointer().String()

				idx := strings.LastIndex(ref, "/")

				if idx > 0 {
					ref = strings.TrimPrefix(ref[idx:], "/")
				}

				def := s.swag.Definitions[ref]
				buf.ws("- `default`, type: `%s`\n", ref)
				obj(1, def, buf)
			}
		}
		buf.ws("</details>\n\n")
	}

	// sort keys so we can have predictable output
	keys := make([]string, len(s.swag.Paths.Paths))
	i := 0
	for key := range s.swag.Paths.Paths {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		route := s.swag.Paths.Paths[key]

		if route.Get != nil {
			md("get", key, route.Get)
		}
		if route.Post != nil {
			md("post", key, route.Post)
		}
		if route.Put != nil {
			md("put", key, route.Put)
		}
		if route.Delete != nil {
			md("delete", key, route.Delete)
		}
	}

	err := ioutil.WriteFile(Markdown, buf.buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write file %s", Markdown))
	}

	log.Printf("wrote API summary to ./%s\n", Markdown)

	return nil
}
