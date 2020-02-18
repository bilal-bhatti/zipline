package internal

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type swagger struct {
	swag *spec.Swagger
}

func newSwagger() swagger {
	// init defaults
	swag := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Version:     "1.0.0",
					Title:       "OpenAPI Version 2 Specification",
					Description: "OpenAPI Version 2 Specification",
				},
			},
			Host:     "api.host.com",
			BasePath: "/api",
			Schemes:  []string{"http"},
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
			Parameters:  make(map[string]spec.Parameter),
			Definitions: make(map[string]spec.Schema),
		},
	}

	ert := skema("object")
	ert.Properties["code"] = *spec.Int32Property()
	ert.Properties["status"] = *spec.StringProperty()
	swag.Definitions["Error"] = ert

	return swagger{
		swag: swag,
	}
}

func (s swagger) generate(packets []*packet) {
	erref := &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/Error"),
		},
	}
	erresp := spec.NewResponse().WithDescription("unexpected error").WithSchema(erref)

	for _, packet := range packets {
		for _, b := range packet.bindings {
			var pi spec.PathItem
			pi, found := s.swag.Paths.Paths[b.path]
			if !found {
				pi = spec.PathItem{}
			}

			op := &spec.Operation{
				OperationProps: spec.OperationProps{
					Description: "Route description",
					Responses: &spec.Responses{
						ResponsesProps: spec.ResponsesProps{
							StatusCodeResponses: make(map[int]spec.Response),
						},
					},
				},
			}

			for i := 0; i < len(b.handler.params); i++ {
				param := b.handler.params[i]

				var template string
				if len(b.paramTemplates) > i {
					template = b.paramTemplates[i]
				}

				if template == "Resolve" {
					continue
				}

				if template == "Path" {
					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:     param.varName(),
							In:       "path",
							Required: true,
						},
						SimpleSchema: spec.SimpleSchema{
							Type: toJSONType(param.signature),
						},
					})
				} else if template == "Query" {
					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:     param.varName(),
							In:       "query",
							Required: false,
						},
						SimpleSchema: spec.SimpleSchema{
							Type: toJSONType(param.signature),
						},
					})

				} else {
					skema := field("--", param.varType.Type())

					ref := spec.Schema{
						SchemaProps: spec.SchemaProps{
							Ref: spec.MustCreateRef("#/definitions/" + param.shortSignature()),
						},
					}

					s.swag.Definitions[param.shortSignature()] = skema

					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:     "body",
							In:       "body",
							Required: true,
							Schema:   &ref,
						},
					})
				}
			}

			for _, ret := range b.handler.returns {
				if ret.varType.Type().String() == "error" {
					op.Responses.ResponsesProps.Default = erresp
					continue
				}

				sk := field("--", ret.varType.Type())
				ref := spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: spec.MustCreateRef("#/definitions/" + ret.shortSignature()),
					},
				}

				s.swag.Definitions[ret.shortSignature()] = sk

				var schema spec.Schema
				if sk.Items != nil {
					schema = skema("array")
					schema.Items = &spec.SchemaOrArray{Schema: &ref}
				} else {
					schema = ref
				}

				op.Responses.ResponsesProps.StatusCodeResponses[200] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "200 response",
						Schema:      &schema,
					},
				}
			}

			if len(b.handler.returns) == 1 && b.handler.returns[0].varType.Type().String() == "error" {
				op.Responses.ResponsesProps.StatusCodeResponses[204] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "no content",
					},
				}
			}

			switch strings.ToUpper(b.template) {
			case "GET":
				pi.Get = op
			case "POST":
				pi.Post = op
			case "DELETE":
				pi.Delete = op
			case "PATCH":
				pi.Patch = op
			case "PUT":
				pi.Put = op
			default:
				panic("Oh noos what's this: " + b.template)
			}

			s.swag.Paths.Paths[b.path] = pi
		}
	}

	s.write()
	s.markdown()
}

func (s swagger) write() {
	// merge some info if spec exists
	err := s.readAndMergeSchema()
	if err != nil {
		// nothing to do
	}

	bites, err := json.MarshalIndent(s.swag, "", "  ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	err = ioutil.WriteFile(OpenAPIFile, bites, 0644)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get working directory"))
	}

	log.Println("Wrote OpenAPI spec to", path.Join(wd, OpenAPIFile))
}

func (s swagger) readAndMergeSchema() error {
	bites, err := ioutil.ReadFile(OpenAPIFile)
	if err != nil {
		return err
	}

	old := &spec.Swagger{}
	err = old.UnmarshalJSON(bites)
	if err != nil {
		return err
	}

	s.swag.Swagger = old.Swagger
	s.swag.Info = old.Info
	s.swag.Host = old.Host
	s.swag.BasePath = old.BasePath
	s.swag.Schemes = old.Schemes
	s.swag.Consumes = old.Consumes
	s.swag.Produces = old.Produces

	return nil
}

func field(name string, t types.Type) spec.Schema {

	switch tt := t.(type) {
	case *types.Pointer:
		return field(name, tt.Elem())
	case *types.Struct:
		skema := skema("object")
		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)

			var tn string
			if tt.Tag(i) != "" {
				tags, err := structtag.Parse(tt.Tag(i))
				if err != nil {
					panic(err)
				}

				jsonTag, err := tags.Get("json")
				if err != nil {
					panic(err)
				}

				trace("name", jsonTag.Name)
				tn = jsonTag.Name
			}

			// TODO: handle recursive refs
			// what to do when Account struct has an Account field?
			if f.Name() != name && tn != "-" {
				fs := field(f.Name(), f.Type())
				if f.Embedded() && tn == "" {
					for k, v := range fs.Properties {
						skema.Properties[k] = v
					}
				} else {
					var named = tn
					if tn == "" {
						named = f.Name()
					}

					skema.Properties[named] = fs
				}
			}
		}
		return skema
	case *types.Named:
		return field(name, tt.Underlying())
	case *types.Slice:
		arr := skema("array")
		sk := field(name, tt.Elem())
		arr.Items = &spec.SchemaOrArray{
			Schema: &sk,
		}
		return arr
	case *types.Basic:
		return skema(tt.String())
	case *types.Map:
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
			},
		}
	default:
		panic(fmt.Sprint("Oh noos unknown type: ", reflect.TypeOf(tt)))
	}
}

func skema(t string) spec.Schema {
	switch t {
	case "object":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:       spec.StringOrArray{"object"},
				Properties: make(map[string]spec.Schema),
			},
		}
	case "int", "int64", "uint", "uint8", "uint64":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{toJSONType(t)},
			},
		}
	case "float64":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{toJSONType(t)},
			},
		}
	case "string", "array":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{t},
			},
		}
	case "bool":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{toJSONType(t)},
			},
		}
	case "byte":
		// TODO: Swagger 2.0 doesn't support binary data type
		// map byte to string, but it should really be removed
		// fail at this point? maybe a strict mode vs lax mode?
		log.Print("Swagger 2.0 doesn't support byte type, generated spec will not be valid")
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{t},
			},
		}
	default:
		panic("Oh noos, what's this: " + t)
	}
}

func toJSONType(t string) string {
	var typeMap = map[string]string{
		"int":     "integer",
		"int8":    "integer",
		"int32":   "integer",
		"int64":   "integer",
		"uint":    "integer",
		"uint8":   "integer",
		"uint32":  "integer",
		"uint64":  "integer",
		"float64": "number",
		"bool":    "boolean",
	}

	jt := typeMap[t]
	if jt == "" {
		return t
	}
	return jt
}

func obj(lvl int, s spec.Schema, buf buffer) {
	if s.Items != nil {
		buf.ws("%s- items\n", strings.Repeat("\t", lvl))
		obj(lvl+1, *s.Items.Schema, buf)
	}

	for k, v := range s.Properties {

		if v.Type.Contains("object") {
			buf.ws("%s- name: `%s`, type: `%s`\n", strings.Repeat("\t", lvl), k, v.Type[0])
			obj(lvl+1, v, buf)
		} else if v.Type.Contains("array") {
			log.Println("array")
			buf.ws("%s- name: `%s`, type: `[]%s`\n", strings.Repeat("\t", lvl), k, v.Type[0])
			obj(lvl+1, *v.Items.Schema, buf)
		} else {
			buf.ws("%s- name: `%s`, type: `%s`\n", strings.Repeat("\t", lvl), k, v.Type[0])
		}
	}
}

func (s swagger) markdown() {
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

	for k, route := range s.swag.Paths.Paths {
		if route.Get != nil {
			md("get", k, route.Get)
		}
		if route.Post != nil {
			md("post", k, route.Post)
		}
		if route.Put != nil {
			md("put", k, route.Put)
		}
		if route.Delete != nil {
			md("delete", k, route.Delete)
		}
	}

	err := ioutil.WriteFile("API.md", buf.buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get working directory"))
	}

	log.Println("Wrote API.md to", path.Join(wd, "API.md"))
}
