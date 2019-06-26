package cyoa

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"
)

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var tmpl *template.Template

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
      <p>{{.}}</p>
    {{end}}
    <ul>
      {{range .Options}}
        <li>
          <a href="/{{.Arc}}">{{.Text}}</a>
        </li>
      {{end}}
    </ul>
  </body>
</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Start with the intro story
	err := tmpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}

// Helper to read in stories
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story : self-explanatory
type Story map[string]Arc

// Arc : A path through the choose-your-own-adventure story
type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option for each Arc
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
