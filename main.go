package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/Jeffail/tunny"
	"github.com/cbroglie/mustache"
	"github.com/stretchr/stew/objects"
)

// default demo content
const (
	paylaod = `{
		"name": "Alice",
		"age":333,
		"users":[
			{
				"name":"dalong",
				"age":"333"
			}
		]
	}`
)

func main() {
	pool := tunny.NewFunc(1000, func(payload interface{}) interface{} {
		log.Println(string(payload.([]byte)))
		jsonmap, err := objects.NewMapFromJSON(string(payload.([]byte)))
		if err != nil {
			log.Println(err)
			return nil
		}
		data, err := mustache.RenderFileInLayout("app.mustache", "layout.mustache", jsonmap.MSI())
		if err != nil {
			log.Println(err)
			return nil
		}
		log.Println(data)
		return data
	})
	defer pool.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := pool.Process([]byte(paylaod))
		if result != nil {
			msg := result.(string)
			w.Write([]byte(msg))
			return
		}
		w.Write([]byte("unknow"))
	})
	http.ListenAndServe(":8080", nil)
}
