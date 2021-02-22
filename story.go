package cyoa

import(
	"encoding/json"
	"io"
	"net/http"
	"html/template"
	"strings"
	"log"
)

func init() {
	tpl = template.Must(template.New("").Parse(htmlTemplate))
}

var tpl *template.Template

var htmlTemplate = `
<!doctype html>
<html class="no-js" lang="">

<head>
  <meta charset="utf-8">
  <title>Golang - Create your own adventure</title>
  <meta name="description" content="">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}

        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>
`

// Decodes a JSON file and returns contents 
func JsonStory(r io.Reader) (Story, error) {
	d:= json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// For handling a web request
func NewHandler(s Story) http.Handler {
  return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// trim the / from what ever the path value is so we can easily pass this to our map to get the correct chapter based on path
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			panic(err)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

type Story map[string]Chapter

// Chapter represents the format a chapter should be in
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option represents an option a user has regarding which chapter to proceed to
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
