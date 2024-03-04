package docparser

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

type DocData struct {
	Doc      string
	Raw      []string
	Comments []string
	Data     map[string]interface{}
}

func (d *DocData) Parameter(varName string) (map[string]interface{}, bool) {
	if params, ok := d.Data["parameters"]; ok {
		params := params.([]interface{})
		for _, p := range params {
			param := p.(map[string]interface{})
			if param["name"] == varName {
				return param, true
			}
		}
	}

	return nil, false
}

func ParseDoc(pkgs []*packages.Package, doc string) (*DocData, error) {
	// TODO: investigate using Go doc formatter to parse blocks of comments
	// var p comment.Parser
	// doc := p.Parse(doc)
	// fmt.Println()
	// for _, c := range doc.Content {
	// 	switch v := c.(type) {
	// 	case *comment.Code:
	// 		fmt.Printf("Code block: %v\n\n", v.Text)
	// 	case *comment.Heading:
	// 		fmt.Printf("Heading: %v\n\n", v.Text)
	// 	case *comment.Paragraph:
	// 		fmt.Printf("Paragraph: %v\n\n", v.Text)
	// 	}
	// }

	var doclines []string
	spec := make(map[string]interface{})
	lines := strings.Split(doc, "\n")

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if strings.HasPrefix(lines[i], "@") {
			kv := strings.SplitN(strings.TrimPrefix(lines[i], "@"), " ", 2)

			kv[0], kv[1] = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			for ; i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "@") && len(lines[i+1]) > 0; i++ {
				kv[1] = kv[1] + lines[i+1]
			}

			err := ParseLine(pkgs, strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]), spec)
			if err != nil {
				return nil, err
			}
		} else {
			doclines = append(doclines, strings.TrimSpace(lines[i]))
		}
	}

	return &DocData{Doc: doc, Raw: lines, Comments: doclines, Data: spec}, nil
}

func ParseLine(pkgs []*packages.Package, key, value string, data map[string]interface{}) error {
	switch key {
	case "schemes", "consumes", "produces", "tags":
		p, err := parseStringValue(value)
		if err != nil {
			return err
		}

		if pl, ok := p.([]string); ok {
			data[key] = pl
		} else {
			if ov, ok := data[key]; ok {
				data[key] = append(ov.([]string), p.(string))
			} else {
				data[key] = []string{p.(string)}
			}
		}
	case "parameters":
		p, err := parseValue(pkgs, value)
		if err != nil {
			return err
		}

		if pl, ok := p.([]interface{}); ok {
			data[key] = pl
		} else {
			if ov, ok := data[key]; ok {
				data[key] = append(ov.([]interface{}), p)
			} else {
				data[key] = []interface{}{p}
			}
		}
	default:
		keys := strings.Split(key, ".")
		for i := 0; i < len(keys); i++ {
			if i+1 == len(keys) {
				pv, err := parseValue(pkgs, value)
				if err != nil {
					return err
				}
				if pl, ok := pv.([]interface{}); ok {
					data[keys[i]] = pl
				} else {
					data[keys[i]] = pv
				}
			} else {
				if ov, ok := data[keys[i]]; ok {
					if ov, ok := ov.(map[string]interface{}); ok {
						data = ov
					} else {
						return errors.New("expecting a map for key: " + key)
					}
				} else {
					temp := make(map[string]interface{})
					data[keys[i]] = temp
					data = temp
				}
			}
		}
	}
	return nil
}

func parseValue(pkgs []*packages.Package, value string) (interface{}, error) {
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		obj := make(map[string]interface{})
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			// attempt to parse struct in format `{packagename.structname}`
			data, err := FindStruct(pkgs, strings.Trim(value, "{}"))
			if err == nil {
				return data, err
			}

			log.Printf("failed to parse annotation `%s`, error: %v", value, err)
		}
		return obj, nil
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		obj := make([]interface{}, 0)
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			log.Printf("failed to parse annotation `%s`, error: %v", value, err)
		}
		return obj, nil
	}

	return value, nil
}

func parseStringValue(value string) (interface{}, error) {
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		obj := make([]string, 0)
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			log.Printf("failed to parse annotation `%s`, error: %v", value, err)
		}
		return obj, nil
	}

	return value, nil
}
