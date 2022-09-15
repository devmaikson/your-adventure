package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"your-adventure/handler"
)

var (
	defaultPort     = 8080
	defaultFilename = "stories/gopher.json"
)

func main() {

	port := flag.Int("p", defaultPort, "set the port to be used for the website")
	filename := flag.String("f", defaultFilename, "set the filename with -f")
	flag.Parse()

	// read JSON to []byte data
	data, err := handler.ReadJsonFile(*filename)

	if err != nil {
		panic(err)
	}
	// retrieve Story from []byte data
	story, err := handler.ExtractStoryFromData(data)
	if err != nil {
		panic(err)
	}
	//fmt.Println(story)

	tpl := template.Must((template.New("").Parse(storyTmpl)))
	//tpl := template.Must(template.New("").Parse("Hello Sir!"))
	//storyHandler := handler.NewHandler(story, handler.WithTemplate(tpl))
	storyHandler := handler.NewHandler(story, handler.WithTemplate(tpl), handler.WithPathFunc(pathFn))
	mux := http.NewServeMux()
	mux.Handle("/story/", storyHandler)
	// http listen func
	fmt.Printf("Starting server on :%d\n", *port)
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
    <meta charset="UTF-8">
    <title>Adventure Story</title>
</head>

<body>
    <section class="page">
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
        <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
            <li><a href="/{{.StoryChapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </section>
    <style>
        body {
            font-family: helvetica, arial;
        }

        h1 {
            text-align: center;
            position: relative;
        }

        .page {
            width: 80%;
            max-width: 500px;
            margin: auto;
            margin-top: 40px;
            margin-bottom: 40px;
            padding: 80px;
            background: #FFFCF6;
            border: 1px solid #eee;
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

        p {
            text-indent: 1em;
        }
    </style>
</body>

</html>`
