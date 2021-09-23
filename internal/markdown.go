package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/util"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

func (s swagger) markdown() error {
	buf := util.NewBuffer()

	buf.Sprintf("# API Summary\n\n")

	buf.Sprintf("```\n")
	buf.Sprintf("Version:     %s\n", s.swag.SwaggerProps.Info.Version)
	buf.Sprintf("Title:       %s\n", s.swag.SwaggerProps.Info.Title)
	buf.Sprintf("Description: %s\n", s.swag.SwaggerProps.Info.Description)
	buf.Sprintf("Host:        %s\n", s.swag.Host)
	buf.Sprintf("BasePath:    %s\n", s.swag.BasePath)
	buf.Sprintf("Consumes:    %s\n", s.swag.Consumes)
	buf.Sprintf("Produces:    %s\n", s.swag.Produces)
	buf.Sprintf("```\n\n")

	md := func(m, path string, op *spec.Operation) {
		buf.Sprintf("<details>\n")
		buf.Sprintf("<summary>%s: %s</summary>\n", path, m)
		buf.Sprintf("\n\n```\n%s\n```\n\n", op.Summary)

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
				buf.Sprintf("`path parameters`\n")
				for _, p := range path {
					if p.In == "path" {
						buf.Sprintf("- %s: `%s`", p.Name, p.Type)
						if p.Format != "" {
							buf.Sprintf("format: `%s`", p.Format)
						}
						buf.Sprintf("\n")
					}
				}
				buf.Sprintf("\n")
			}

			if len(query) > 0 {
				buf.Sprintf("`query parameters`\n")
				for _, p := range query {
					if p.In == "query" {
						buf.Sprintf("- %s: `%s`", p.Name, p.Type)
						if p.Format != "" {
							buf.Sprintf(", format: `%s`", p.Format)
						}
						buf.Sprintf("\n")
					}
				}
				buf.Sprintf("\n")
			}

			if len(body) > 0 {
				buf.Sprintf("`body parameter`\n")

				for _, p := range body {
					if p.In == "body" {
						if p.Schema != nil {
							path := strings.Split(p.Schema.Ref.GetPointer().String(), "/")[2]
							def := s.swag.Definitions[path]
							buf.Sprintf("- %s: `%s`\n", p.Name, path)
							obj(1, def, buf)
						}
					}
				}
			}
		}
		buf.Sprintf("\n")

		buf.Sprintf("`responses`\n")
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
						buf.Sprintf("- code: `%d`, type: `[]%s`\n", code, ref)
						obj(1, def, buf)
					} else {
						buf.Sprintf("- code: `%d`, type: `%s`\n", code, ref)
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
				buf.Sprintf("- `default`, type: `%s`\n", ref)
				obj(1, def, buf)
			}
		}
		buf.Sprintf("</details>\n\n")
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
		if route.Patch != nil {
			md("patch", key, route.Patch)
		}
	}

	err := ioutil.WriteFile(Markdown, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write file %s", Markdown))
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Printf("wrote API summary to  %s\n", path.Join(cwd, Markdown))

	return nil
}

func obj(lvl int, s spec.Schema, buf *util.Buffer) {
	// sort keys so we can have predictable output
	keys := make([]string, len(s.Properties))
	i := 0
	for key := range s.Properties {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		v := s.Properties[key]
		if v.Type.Contains("object") {
			buf.Sprintf("%s- %s: `%s`\n", strings.Repeat("\t", lvl), key, v.Type[0])
			obj(lvl+1, v, buf)
		} else if v.Type.Contains("array") {
			buf.Sprintf("%s- %s: `[]%s`\n", strings.Repeat("\t", lvl), key, v.Type[0])
			obj(lvl+1, *v.Items.Schema, buf)
		} else {
			buf.Sprintf("%s- %s: `%s`", strings.Repeat("\t", lvl), key, v.Type[0])
			if v.Format != "" {
				buf.Sprintf(", format: `%s`", v.Format)
			}
			buf.Sprintf("\n")
		}
	}
}
