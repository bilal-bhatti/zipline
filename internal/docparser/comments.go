package docparser

import (
	"go/token"
	"os"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/util"
	"github.com/fatih/structtag"
	"golang.org/x/tools/go/packages"
)

type Comments struct {
	Raw      []string
	Tags     map[string]*structtag.Tags
	Comments []string
}

func GetDocComments(pkgs []*packages.Package, pos token.Position) (*DocData, error) {
	// just in case
	if pos.Line-2 <= 0 {
		return &DocData{
			// Tags: make(map[string]*structtag.Tags),
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

	return ParseDoc(pkgs, strings.Join(comments, "\n"))
}
