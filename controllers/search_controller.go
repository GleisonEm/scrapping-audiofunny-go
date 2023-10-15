// controllers/search_controller.go

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Parâmetro 'name' não especificado na URL", http.StatusBadRequest)
		return
	}

	searchURL := fmt.Sprintf("https://www.myinstants.com/pt/search/?name=%s", name)
	doc, err := htmlquery.LoadURL(searchURL)
	if err != nil {
		http.Error(w, "Erro ao analisar HTML", http.StatusInternalServerError)
		return
	}

	button := htmlquery.FindOne(doc, "//*[@id='instants_container']/div[1]/div[1]/button")
	if button == nil {
		http.Error(w, "Botão não encontrado", http.StatusNotFound)
		return
	}

	onclick := htmlquery.SelectAttr(button, "onclick")
	if onclick == "" {
		http.Error(w, "Atributo 'onclick' não encontrado", http.StatusNotFound)
		return
	}

	soundPath := extractSoundPath(onclick)
	title := htmlquery.SelectAttr(button, "title")
	if title == "" {
		http.Error(w, "Atributo 'title' não encontrado", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"soundPath": soundPath,
		"title":     title,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func extractSoundPath(onclick string) string {
	format := strings.Replace(onclick, "play(", "", -1)
	splitParts := strings.Split(format, ",")

	return strings.Trim(splitParts[0], "'")
}
