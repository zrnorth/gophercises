package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
	cmdLineTmpl = template.Must(template.New("").Parse(defaultTextTemplate))
}

var tmpl, cmdLineTmpl *template.Template

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

var defaultTextTemplate = `
	{{.Title}}
	{{range .Paragraphs}}
		{{.}}
	{{end}}
	
	{{range .Options}}
		{{.Idx}}:  {{.Text}}
  {{end}}
`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func WithCommandLineMode() HandlerOption {
	return WithTemplate(cmdLineTmpl)
}

func NewHandler(s Story, opts ...HandlerOption) handler {
	h := handler{s, tmpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:] // remove the '/'
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if arc, ok := h.s[path]; ok {
		err := h.t.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func (h handler) ServeTextToConsole(path string) {
	if arc, ok := h.s[path]; ok {
		err := h.t.Execute(os.Stdout, arc)
		if err != nil {
			log.Printf("%v", err)
		}
		return
	}
}

// Given a path and a response, returns the next path in the sequence
func (h handler) GetNext(path string, resp int) string {
	if arc, ok := h.s[path]; ok {
		if len(arc.Options) == 0 {
			// We reached the end of a path. Finish the game
			// Kinda dumb but we treat empty string like "the end"
			return ""
		}
		for _, opt := range arc.Options {
			if opt.Idx == resp {
				return opt.Arc
			}
		}
	}
	return path
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
	Idx  int    `json:"idx"`
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
