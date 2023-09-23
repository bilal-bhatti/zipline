package internal

import (
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
)

type comments struct {
	raw      []string
	tags     map[string]*structtag.Tags
	comments []string
}

func getComments(pos token.Position) (*comments, error) {
	// just in case
	if pos.Line-2 <= 0 {
		return &comments{
			tags: make(map[string]*structtag.Tags),
		}, nil
	}

	fileBytes, err := os.ReadFile(pos.Filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileBytes), "\n")
	var comments []string

	// start from func declaration and go backwards
	// stop when non comment line found
	for i := pos.Line - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
			if len(line) == 0 {
				continue
			}

			comments = append(comments, line)
		} else {
			break
		}
	}
	reverse(comments)
	return getParsedComments(strings.Join(comments, "\n"))
}

func getCommentss(pos token.Position) (*comments, error) {
	fileBytes, err := os.ReadFile(pos.Filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileBytes), "\n")

	comms := &comments{
		tags: make(map[string]*structtag.Tags),
	}

	// just in case
	if pos.Line-2 <= 0 {
		return comms, nil
	}

	// start from func declaration and go backwards
	// stop when non comment line found
	for i := pos.Line - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
			if len(line) == 0 {
				continue
			}

			if strings.HasPrefix(line, "@") {
				tagline := strings.TrimPrefix(line, "@")
				split := strings.SplitN(tagline, " ", 2)

				if len(split) == 2 {
					tags, err := structtag.Parse(parseTagIfAny(split[1]))
					if err == nil {
						comms.tags[strings.TrimSpace(split[0])] = tags
					}
				}
			} else {
				comms.comments = append(comms.comments, line)
			}
			comms.raw = append(comms.raw, line)
		} else {
			break
		}
	}

	reverse(comms.raw)
	return comms, nil
}

func getParsedComments(docs string) (*comments, error) {
	// fmt.Println("parsing", docs)
	lines := strings.Split(docs, "\n")

	comms := &comments{
		tags: make(map[string]*structtag.Tags),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "@") {
			tagline := strings.TrimPrefix(line, "@")
			split := strings.SplitN(tagline, " ", 2)

			if len(split) == 2 {
				tags, err := structtag.Parse(parseTagIfAny(split[1]))
				if err == nil {
					comms.tags[strings.TrimSpace(split[0])] = tags
				}
			}
		} else {
			comms.comments = append(comms.comments, line)
		}
		comms.raw = append(comms.raw, line)
	}

	return comms, nil
}

// parse tags, with the same syntax as standard Go tags.
// `json:"blah,x,y,z" foo:"bar"` ignore otherwise
func parseTagIfAny(str string) (result string) {
	const sep string = "`"
	start := strings.Index(str, sep)
	if start == -1 {
		return
	}
	start += len(sep)
	end := strings.Index(str[start:], sep)
	if end == -1 {
		return
	}
	return str[start : start+end]
}
