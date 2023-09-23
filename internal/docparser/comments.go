package docparser

import (
	"go/token"
	"os"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/util"
	"github.com/fatih/structtag"
)

type Comments struct {
	Raw      []string
	Tags     map[string]*structtag.Tags
	Comments []string
}

func GetComments(pos token.Position) (*Comments, error) {
	// just in case
	if pos.Line-2 <= 0 {
		return &Comments{
			Tags: make(map[string]*structtag.Tags),
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
	util.Reverse(comments)
	return ParsedComments(strings.Join(comments, "\n"))
}

func ParsedComments(docs string) (*Comments, error) {
	// fmt.Println("parsing", docs)
	lines := strings.Split(docs, "\n")

	comms := &Comments{
		Tags: make(map[string]*structtag.Tags),
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
					comms.Tags[strings.TrimSpace(split[0])] = tags
				}
			}
		} else {
			comms.Comments = append(comms.Comments, line)
		}
		comms.Raw = append(comms.Raw, line)
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
