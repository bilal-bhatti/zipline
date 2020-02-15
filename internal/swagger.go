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

			for _, param := range b.handler.params {
				if param.varType.Type().String() == "context.Context" {
					continue
				}

				// if strings.HasSuffix(param.varType.Type().String(), "util.Claims") {
				// 	// log.Println("skipping claims type")
				// 	continue
				// }

				// TODO: do better!
				if isPathParam(b, param) {
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

				skema := field("--", ret.varType.Type())
				ref := spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: spec.MustCreateRef("#/definitions/" + ret.shortSignature()),
					},
				}

				s.swag.Definitions[ret.shortSignature()] = skema

				op.Responses.ResponsesProps.StatusCodeResponses[200] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "200 response",
						Schema:      &ref,
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
