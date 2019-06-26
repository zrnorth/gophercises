package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
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
		<style>
			body {
				margin: 20px;
				background-color: #d8e7ff;
				font-family: Georgia;
				color: #222;
			}
			h1 {
				font-family: Verdana;
			}
			ul {
				list-style: none;
			}
		</style>
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

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tmpl}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:] // remove the '/'

	if arc, ok := h.s[path]; ok {
		err := tmpl.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
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
