package main

import (
	adventure "tanachtextadventure"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 3000, "port to start web app")
	filename := flag.String("file", "storyData.json", "JSON file with choose your own adventure structured story")
	flag.Parse()
	fmt.Printf("using the story in %s. \n.", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	story, err := adventure.JsonStory(f)
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Parse(storyTmpl))
	h := adventure.NewHandler(story,
		adventure.WithTemplate(tpl),
		adventure.WithPathFunc(pathFn),
	)
	mux := http.NewServeMux()
	mux.Handle("/story", h)

	mux.Handle("/", adventure.NewHandler(story))
	// handle static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("starting the server at port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Allegorical Adventure Game - Jacob And Esav: Gangster Vegas</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      <ul>
      {{range .Options}}
        <li>
          <a href="/{{.Chapter}}" onclick="onClick({{.Sound}}, {{.Image}})">{{.Text}}</a>
        </li>
      {{end}}
      </ul>
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #797;
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
        text-decoration: underline;
        color: #555;
      }
      a:active,
      a:hover {
        color: #222;
      }
      p {
        text-indent: 1em;
      }
    </style>
    <script>
      function onClick(audioPath, imagePath) {
        playSound(audioPath);
        flashImage(imagePath);
      }

      const playSound = (audioPath) => {
        if (audioPath != "") {
          let audio = new Audio(audioPath);
          // add audio play to event queue
          setTimeout(() => {
            audio.play();
          }, 500);

        }
        else {
          return
        }
      }
      const flashImage = (imagePath) => {
        if (imagePath != "") {
          let flash = document.createElement('img');
          flash.src = imagePath;
          flash.style.position = 'absolute';
          flash.style.left = '0px';
          flash.style.top = '0px';
          flash.style.width = '100%';
          flash.style.height = '100%';
          flash.style.zIndex = '10';
          document.body.appendChild(flash);
          setTimeout(function() {
            document.body.removeChild(flash);
          }, 2000);
        }
        else {
          return
        }
      }

    </script>

  </body>
</html>`
