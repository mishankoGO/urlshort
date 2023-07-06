package urlshort

import (
	"encoding/json"
	"fmt"
	"github.com/mishankoGO/urlshort/repository"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.RequestURI]; ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(ymlPath string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(ymlPath)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(ymlPath string) ([]map[string]string, error) {
	var paths []map[string]string

	ymlFile, err := os.Open(ymlPath)
	if err != nil {
		return nil, err
	}

	yml := yaml.NewDecoder(ymlFile)
	err = yml.Decode(&paths)
	if err != nil {
		return nil, err
	}

	fmt.Println(paths)
	return paths, nil
}

func JSONHandler(jsonPath string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonPath)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonPath string) ([]map[string]string, error) {
	var paths []map[string]string

	ymlFile, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}

	js := json.NewDecoder(ymlFile)
	err = js.Decode(&paths)
	if err != nil {
		return nil, err
	}

	fmt.Println(paths)
	return paths, nil
}

func buildMap(paths []map[string]string) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, path := range paths {
		pathsToUrls[path["path"]] = path["url"]
	}
	fmt.Println(pathsToUrls)
	return pathsToUrls
}

func DBHandler(repo *repository.Repo, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := repo.View(r.RequestURI)
		if url != "" {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
