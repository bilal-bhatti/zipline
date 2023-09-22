package internal

import (
	"encoding/json"
	"fmt"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type swagger struct {
	swag      *spec.Swagger
	typeSpecs map[string]*typeSpecWithPkg
}

func newSwagger(typeSpecs map[string]*typeSpecWithPkg) (*swagger, error) {
	// init defaults
	swag := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			// Swagger: "2.0",
			// Info: &spec.Info{
			// 	InfoProps: spec.InfoProps{
			// 		Version:     "1.0.0",
			// 		Title:       "OpenAPI Version 2 Specification",
			// 		Description: "OpenAPI Version 2 Specification",
			// 	},
			// },
			// Host:     "api.host.com",
			// BasePath: "/api",
			// Schemes:  []string{"http", "https"},
			// Consumes: []string{"application/json"},
			// Produces: []string{"application/json"},
			// Paths: &spec.Paths{
			// 	Paths: make(map[string]spec.PathItem),
			// },
			// Parameters:  make(map[string]spec.Parameter),
			// Definitions: make(map[string]spec.Schema),
		},
	}

	// ert, err := skema("object")
	// if err != nil {
	// 	return nil, err
	// }
	// ert.Description = "error response object"
	// ert.Properties["code"] = spec.Schema{
	// 	SchemaProps: spec.SchemaProps{
	// 		Type: spec.StringOrArray{"integer"},
	// 	},
	// }
	// ert.Properties["status"] = *spec.StringProperty()
	// swag.Definitions["Error"] = *ert

	return &swagger{
		swag:      swag,
		typeSpecs: typeSpecs,
	}, nil
}

func (s swagger) generate(packets []*packet) error {
	erref := &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/Error"),
		},
	}
	erresp := spec.NewResponse().WithDescription("unexpected error").WithSchema(erref)

	for _, packet := range packets {
		docs, err := parsedocs(packet.funcDecl.Doc.Text())
		if err != nil {
			return err
		}

		docsbytes, err := json.Marshal(docs)
		if err != nil {
			return err
		}

		err = json.Unmarshal(docsbytes, s.swag)
		if err != nil {
			return err
		}

		// populate some defaults
		s.swag.SwaggerProps.Paths = &spec.Paths{
			Paths: make(map[string]spec.PathItem),
		}
		s.swag.SwaggerProps.Parameters = make(map[string]spec.Parameter)
		s.swag.SwaggerProps.Definitions = make(map[string]spec.Schema)

		ert, err := skema("object")
		if err != nil {
			return err
		}
		ert.Description = "error response object"
		ert.Properties["code"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"integer"},
			},
		}
		ert.Properties["status"] = *spec.StringProperty()
		s.swag.Definitions["Error"] = *ert
		// defaults

		for _, b := range packet.bindings {
			var pi spec.PathItem
			pi, found := s.swag.Paths.Paths[b.path]
			if !found {
				pi = spec.PathItem{}
			}

			op := &spec.Operation{
				OperationProps: spec.OperationProps{
					Description: "Route description",
					ID:          b.id(),
					Tags:        []string{b.handler.pkg[strings.LastIndex(b.handler.pkg, "/")+1:]},
					Responses: &spec.Responses{
						ResponsesProps: spec.ResponsesProps{
							StatusCodeResponses: make(map[int]spec.Response),
						},
					},
				},
			}

			if len(b.handler.comments.raw) > 0 {
				op.Description = strings.Join(b.handler.comments.raw, "\n")
				op.Summary = op.Description
			}

			for i := 0; i < len(b.handler.params); i++ {
				param := b.handler.params[i]

				var template string
				if len(b.paramTemplates) > i {
					template = b.paramTemplates[i]
				}

				if template != "Body" && template != "Path" && template != "Query" {
					continue
				}

				if template == "Path" || template == "Query" {
					simpleSchema := paramSimpleSchema(param)
					tags, ok := b.handler.comments.tags[param.VarName()]
					if ok {
						fmtTag, err := tags.Get("format")
						if err == nil {
							if len(fmtTag.Options) > 0 {
								simpleSchema.Format = fmt.Sprintf("%s,%s", fmtTag.Name, strings.Join(fmtTag.Options, ","))
							} else {
								simpleSchema.Format = fmtTag.Name
							}
						}
					}

					var paramComment string = ""
					for _, c := range b.handler.comments.raw {
						ctag := "@" + strings.TrimSpace(param.VarName())
						if strings.HasPrefix(c, ctag) {
							paramComment = strings.TrimSpace(strings.TrimLeft(c, ctag))
							break
						}
					}

					required := true
					if template == "Query" {
						required = !param.IsPtr
					}

					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:        param.VarName(),
							In:          strings.ToLower(template),
							Required:    required,
							Description: paramComment,
						},
						SimpleSchema: simpleSchema,
					})
				} else {
					skema, err := field("--", param.VarType.Type())
					if err != nil {
						return err
					}

					ts := s.typeSpecs[param.Signature]
					if ts != nil {
						pos := ts.pkg.Fset.PositionFor(ts.typeSpec.Pos(), true)
						comms, err := getComments(pos)
						if err != nil {
							// let's not fail on comments but log the error
							log.Println("failed to extract comments", err.Error())
						}

						skema.Description = strings.Join(comms.raw, "\n")
					}

					ref := spec.Schema{
						SchemaProps: spec.SchemaProps{
							Ref: spec.MustCreateRef("#/definitions/" + param.SimpleSignature(packet.pkg.PkgPath)),
						},
					}

					s.swag.Definitions[param.SimpleSignature(packet.pkg.PkgPath)] = *skema

					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:        "body",
							In:          "body",
							Required:    true,
							Description: skema.Description,
							Schema:      &ref,
						},
					})
				}
			}

			for _, ret := range b.handler.returns {
				if ret.VarType.Type().String() == "error" {
					op.Responses.ResponsesProps.Default = erresp
					continue
				}

				sk, err := field("--", ret.VarType.Type())
				if err != nil {
					return err
				}

				ts := s.typeSpecs[ret.Signature]
				if ts != nil {
					pos := ts.pkg.Fset.PositionFor(ts.typeSpec.Pos(), true)
					comms, err := getComments(pos)
					if err != nil {
						// let's not fail on comments but log the error
						log.Println("failed to extract comments", err.Error())
					}

					sk.Description = strings.Join(comms.raw, "\n")
				}

				ref := spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: spec.MustCreateRef("#/definitions/" + ret.SimpleSignature(packet.pkg.PkgPath)),
					},
				}

				s.swag.Definitions[ret.SimpleSignature(packet.pkg.PkgPath)] = *sk

				var schema *spec.Schema
				if ret.IsSlice {
					schema, err = skema("array")
					if err != nil {
						return err
					}
					schema.Items = &spec.SchemaOrArray{Schema: &ref}
				} else {
					schema = &ref
				}

				op.Responses.ResponsesProps.StatusCodeResponses[200] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "200 success response",
						Schema:      schema,
					},
				}
			}

			if len(b.handler.returns) == 1 && b.handler.returns[0].VarType.Type().String() == "error" {
				op.Responses.ResponsesProps.StatusCodeResponses[204] = spec.Response{
					ResponseProps: spec.ResponseProps{
						Description: "no content",
					},
				}
			}

			switch strings.ToUpper(b.method) {
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
				return errors.New(fmt.Sprintf("oh noos what's this method: %s", b.template))
			}

			s.swag.Paths.Paths[b.path] = pi
		}
	}

	err := s.write()
	if err != nil {
		return errors.Wrap(err, "OpenAPI spec generation failed")
	}
	err = s.markdown()
	if err != nil {
		return errors.Wrap(err, "failed to generate API.md summary file")
	}

	return nil
}

func (s swagger) write() error {
	// merge some info if spec exists
	// ignore if any errors
	_ = s.readAndMergeSchema()

	bites, err := json.MarshalIndent(s.swag, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to generate write OpenAPI json")
	}
	bites = append(bites, []byte("\n")...)

	err = os.WriteFile(OpenAPIFile, bites, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write OpenAPI json to file")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Printf("wrote OpenAPI v2 spec to  %s\n", path.Join(cwd, OpenAPIFile))
	return nil
}

func (s swagger) readAndMergeSchema() error {
	bites, err := os.ReadFile(OpenAPIFile)
	if err != nil {
		return err
	}

	old := &spec.Swagger{}
	err = old.UnmarshalJSON(bites)
	if err != nil {
		return err
	}

	if s.swag.Swagger != "" {
		s.swag.Swagger = old.Swagger
	}

	if s.swag.Info != nil {
		s.swag.Info = old.Info
	}

	if s.swag.Host != "" {
		s.swag.Host = old.Host
	}

	if s.swag.BasePath != "" {
		s.swag.BasePath = old.BasePath
	}

	if len(s.swag.Schemes) == 0 {
		s.swag.Schemes = old.Schemes
	}

	if len(s.swag.Consumes) == 0 {
		s.swag.Consumes = old.Consumes
	}

	if len(s.swag.Produces) == 0 {
		s.swag.Produces = old.Produces
	}

	if len(s.swag.SecurityDefinitions) == 0 {
		s.swag.SecurityDefinitions = old.SecurityDefinitions
	}

	if len(s.swag.Security) == 0 {
		s.swag.Security = old.Security
	}

	return nil
}

func field(name string, t types.Type) (*spec.Schema, error) {
	switch tt := t.(type) {
	case *types.Pointer:
		return field(name, tt.Elem())
	case *types.Struct:
		skema, err := skema("object")
		if err != nil {
			return nil, err
		}

		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)
			if !f.Exported() {
				continue
			}

			fieldTypeToken := tokens.NewTypeToken(f.Type().String(), f.Name())

			var jsonName string
			var tag = tt.Tag(i)
			var tags *structtag.Tags
			if tag != "" {
				tags, err = structtag.Parse(tag)
				if err != nil {
					return nil, err
				}

				jsonTag, err := tags.Get("json")
				if err != nil {
					return nil, err
				}
				jsonName = jsonTag.Name
			}

			// TODO: handle recursive refs
			// what to do when Account struct has an Account field?
			if f.Name() != name && jsonName != "-" {
				var fs *spec.Schema

				if fieldTypeToken.Signature == "time.Time" {
					fs = spec.DateTimeProperty()
				} else {
					fs, err = field(f.Name(), f.Type())
					if err != nil {
						return nil, err
					}
				}

				if tags != nil {
					timeFmt, err := tags.Get("format")
					if err == nil {
						fs.Format = timeFmt.Value()
					}
				}

				if f.Embedded() && jsonName == "" {
					for k, v := range fs.Properties {
						skema.Properties[k] = v
					}
				} else {
					var named = jsonName
					if jsonName == "" {
						named = f.Name()
					}

					skema.Properties[named] = *fs
				}
			}
		}
		return skema, nil
	case *types.Named:
		return field(name, tt.Underlying())
	case *types.Slice:
		obj, err := field(name, tt.Elem())
		if err != nil {
			return nil, err
		}
		return spec.ArrayProperty(obj), nil
	case *types.Basic:
		return skema(tt.String())
	case *types.Map:
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 spec.StringOrArray{"object"},
				AdditionalProperties: &spec.SchemaOrBool{Allows: true},
			},
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("oh noos unknown type: %v", reflect.TypeOf(tt)))
	}
}

func skema(t string) (*spec.Schema, error) {
	switch t {
	case "object":
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 spec.StringOrArray{"object"},
				Properties:           make(map[string]spec.Schema),
				AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			},
		}, nil
	case "int", "uint":
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"integer"},
			},
		}, nil
	case "int8", "uint8":
		return spec.Int8Property(), nil
	case "int16", "uint16":
		return spec.Int16Property(), nil
	case "int32", "uint32":
		return spec.Int32Property(), nil
	case "int64", "uint64":
		return spec.Int64Property(), nil
	case "float32":
		return spec.Float32Property(), nil
	case "float64":
		return spec.Float64Property(), nil
	case "string":
		return spec.StringProperty(), nil
	case "bool":
		return spec.BoolProperty(), nil
	case "time.Time":
		return spec.DateTimeProperty(), nil
	case "byte":
		// TODO: Swagger 2.0 doesn't support binary data type
		// map byte to string, but it should really be removed
		// fail at this point? maybe a strict mode vs lax mode?
		log.Println("Swagger 2.0 doesn't support byte type, generated spec will not be valid")
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"byte"},
			},
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("oh noos, what's this: %s", t))
	}
}

func paramSimpleSchema(tkn *tokens.TypeToken) spec.SimpleSchema {
	typeSchema, err := skema(tkn.Signature)
	if err != nil {
		return spec.SimpleSchema{
			Type: "string",
		}
	}

	if tkn.IsSlice {
		return spec.SimpleSchema{
			Type: "array",
			Items: &spec.Items{
				SimpleSchema: spec.SimpleSchema{
					Type:   typeSchema.Type[0],
					Format: typeSchema.Format,
				},
			},
		}
	}

	return spec.SimpleSchema{
		Type:   typeSchema.Type[0],
		Format: typeSchema.Format,
	}
}
