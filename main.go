package main

import (
	"fmt"
	"github.com/mishankoGO/urlshort/conf"
	"github.com/mishankoGO/urlshort/repository"
	"github.com/mishankoGO/urlshort/urlshort"
	"log"
	"net/http"
)

var config conf.ShortenerConfig

func main() {
	err := conf.InitFlags(&config)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	repo, err := repository.NewBoltRepo()
	if err != nil {
		log.Fatal(err)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	for key, val := range pathsToUrls {
		err := repo.Update(key, val)
		if err != nil {
			log.Printf("error adding pair to db: %v", err)
		}
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(config.Path, mapHandler)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JSONHandler(config.Path, yamlHandler)
	if err != nil {
		panic(err)
	}
	dbHandler := urlshort.DBHandler(repo, jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
