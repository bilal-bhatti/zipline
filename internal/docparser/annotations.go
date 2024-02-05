package docparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseDocs(docs string) (map[string]interface{}, error) {
	// var p comment.Parser
	// doc := p.Parse(docs)
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

	// var doclines []string
	nested := make(map[string]interface{})
	lines := strings.Split(docs, "\n")

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if strings.HasPrefix(lines[i], "@") {
			kv := strings.SplitN(strings.TrimPrefix(lines[i], "@"), " ", 2)

			kv[0], kv[1] = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			for ; i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "@") && len(lines[i+1]) > 0; i++ {
				kv[1] = kv[1] + lines[i+1]
			}

			fmt.Println("kv", kv[0], kv[1])
			err := ParseLine(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]), nested)
			if err != nil {
				return nil, err
			}
		}
	}
	// for _, line := range lines {
	// 	if strings.HasPrefix(line, "@") {
	// 		kv := strings.SplitN(strings.TrimPrefix(line, "@"), " ", 2)

	// 		kv[0], kv[1] = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
	// 		fmt.Println("parsing, k, v:", kv)

	// 		err := ParseLine(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]), nested)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// 	// else {
	// 	// 	doclines = append(doclines, strings.TrimSpace(line))
	// 	// }
	// }

	yaml.NewEncoder(os.Stdout).Encode(nested)

	return nested, nil
}

func ParseLine(key, value string, data map[string]interface{}) error {
	switch key {
	case "schemes", "consumes", "produces", "tags", "parameters":
		p, err := parseValue(value)
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
				pv, err := parseValue(value)
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

func parseValue(value string) (interface{}, error) {
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		obj := make(map[string]interface{})
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			log.Printf("failed to parse annotation `%s`, error: %v", value, err)
		}
		json.NewEncoder(os.Stdout).Encode(obj)
		return obj, nil
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		obj := make([]interface{}, 0)
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			log.Printf("failed to parse annotation `%s`, error: %v", value, err)
		}
		// json.NewEncoder(os.Stdout).Encode(obj)
		return obj, nil
	}

	return value, nil
}

// func ParseLine_x(key, value string, data map[string]interface{}) error {
// 	switch key {
// 	case "schemes", "consumes", "produces", "tags", "parameters":
// 		p, err := parseValue(value)
// 		if err != nil {
// 			return err
// 		}

// 		if pl, ok := p.([]interface{}); ok {
// 			data[key] = pl
// 		} else {
// 			if ov, ok := data[key]; ok {
// 				data[key] = append(ov.([]interface{}), p)
// 			} else {
// 				data[key] = []interface{}{p}
// 			}
// 		}
// 	default:
// 		keys := strings.Split(key, ".")
// 		for i := 0; i < len(keys); i++ {
// 			if i+1 == len(keys) {
// 				pv, err := parseValue(value)
// 				if err != nil {
// 					return err
// 				}
// 				if pl, ok := pv.([]interface{}); ok {
// 					data[keys[i]] = pl
// 				} else {
// 					data[keys[i]] = pv
// 				}
// 			} else {
// 				if ov, ok := data[keys[i]]; ok {
// 					if ov, ok := ov.(map[string]interface{}); ok {
// 						data = ov
// 					} else {
// 						return errors.New("expecting a map for key: " + key)
// 					}
// 				} else {
// 					temp := make(map[string]interface{})
// 					data[keys[i]] = temp
// 					data = temp
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// func parseValue(value string) (interface{}, error) {
// 	if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
// 		p := make(map[string]interface{})
// 		value = strings.Trim(value, "()")
// 		r := csv.NewReader(strings.NewReader(value))
// 		r.TrimLeadingSpace = true
// 		fields, err := r.Read()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, f := range fields {
// 			kv := strings.Split(f, ":")
// 			p[strings.TrimSpace(kv[0])] = guess(strings.TrimSpace(kv[1]))
// 		}
// 		return p, nil
// 	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
// 		var p []interface{}
// 		value = strings.TrimSpace(strings.Trim(value, "[]"))
// 		if value == "" {
// 			return make([]interface{}, 0), nil
// 		}

// 		// fmt.Println("value", value)
// 		r := csv.NewReader(strings.NewReader(value))
// 		r.TrimLeadingSpace = true

// 		lines, err := r.Read()
// 		if err != nil {
// 			return err, nil
// 		}
// 		for _, f := range lines {
// 			p = append(p, guess(strings.TrimSpace(f)))
// 		}
// 		return p, nil
// 	} else {
// 		return guess(value), nil
// 	}
// }

// func guess(v interface{}) interface{} {
// 	if v, ok := v.(string); ok {
// 		if v == "yes" || v == "true" || v == "on" {
// 			return true
// 		}
// 		if v == "no" || v == "false" || v == "off" {
// 			return false
// 		}
// 	}

// 	return v
// }
