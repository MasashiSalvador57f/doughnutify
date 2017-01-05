package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	keyCurrentValue = "current_value"
	keyTargetValue  = "target_value"
	keyTitle        = "title"
	keyCurrentColor = "current_color"
	keyRestColor    = "rest_color"
)

const templatePath = "index.html"
const defaultPort = 9999

// Index returns rendered html to client :)
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	currentValue := query.Get(keyCurrentValue)
	targetValue := query.Get(keyTargetValue)
	currentColor := query.Get(keyCurrentColor)
	restColor := query.Get(keyRestColor)
	percent := getPercent(currentValue, targetValue)
	title := query.Get(keyTitle)

	if currentColor == "" {
		currentColor = "2ecc71"
	}
	if restColor == "" {
		restColor = "ffffff"
	}

	variables := map[string]interface{}{
		"Percent":      percent,
		"Title":        title,
		"CurrentColor": currentColor,
		"RestColor":    restColor,
	}
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// NOTE: Don't panic .... :)
		panic(err)
	}
	err = tmpl.Execute(w, variables)
	if err != nil {
		// NOTE: Don't panic :)
		panic(err)
	}
}

func getPercent(currentStr, targetStr string) string {
	currentValue, _ := strconv.Atoi(currentStr)
	targetValue, _ := strconv.Atoi(targetStr)

	if targetValue <= 0 {
		return "0.0"
	}
	p := 100 * float32(currentValue) / float32(targetValue)
	return fmt.Sprintf("%.2f", p)
}

func main() {
	port := os.Getenv("DOUGHNUTIFY_PORT")
	if port == "" {
		port = ":9999"
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	log.Fatal(http.ListenAndServe(port, router))
}
