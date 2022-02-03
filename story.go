package tanachtextadventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Allegorical Adventure Game - Jacob And Esav: Gangster Vegas</title>
  </head>

  <body>

  <!-- img element without src attribute -->
  <img id="flash" src="" />
  <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li>
            <button class="option" onclick="clickFunction({{.Chapter}},{{.Sound}}, {{.Image}})">
              {{.Text}}
            </button>
          </li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
        background-repeat: no-repeat;
        background-size: cover;
        background-position: center;
        background-image:url({{.Background}});
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        color: white;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: rgba(0,0,0,0.75);
        border: 1px solid #eee;
        border-radius: 10px;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      button {
        border: 2px solid white;
        background: rgba(0,0,0,0.85);
        color: white;
        font-weight: bold;
        cursor: pointer;
        padding: 10px;
        border-radius: 10px;
      }
      button:hover {
        background: rgba(0,0,0,1);
        border-color: #7792a2;
        color: #7792a2;
      }
      button:focus {
        outline: none;
      }
      p {
        text-indent: 1em;
      }
      #flash {
        display:none;
        position: absolute;
        left: 0px;
        top: 0px;
        width: 100vw;
        height: 100vh;
        z-index: 2;

      }
    </style>
    <script>

      function clickFunction(redirectPath, audioPath, imagePath) {
        flashImage(imagePath);
        playSound(audioPath);

        setTimeout(function(){
          window.location.href = redirectPath;
        }, 1500);

      }

      const playSound = (audioPath) => {
        if (audioPath != "") {
          let audio = new Audio(audioPath);
          audio.play();
        }
        else {
          return
        }
      }
      const flashImage = (imagePath) => {
        if (imagePath != "") {
          let flash = document.getElementById('flash');
          flash.src = imagePath;
          flash.style.display = "block";
          setTimeout(() => {
            flash.style.display = "none";
          }, 1500);
        }
        else {
          return
        }

      }
    </script>
  </body>
</html>`

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

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
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
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
	Background string   `json:"background"`
}

type Option struct {
	Text      string   `json:"text"`
	Chapter   string   `json:"action"`
	Inventory []string `json:"inventory_addition"`
	Sound     string   `json:"sound"`
	Image     string   `json:"image"`
}
