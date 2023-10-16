// controllers/search_controller.go

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"encoding/base64"
	"io/ioutil"
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

func SearchAllHandler(w http.ResponseWriter, r *http.Request) {

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

	buttons := htmlquery.Find(doc, "//*[@class='instant-link link-secondary']")
	if buttons == nil {
		http.Error(w, "Botões não encontrados", http.StatusNotFound)
		return
	}

	// println(htmlquery.OutputHTML(divContent, true))
	println(len(buttons))

	// Itere sobre os botões encontrados
	// for _, button := range buttons {
	// 	println(htmlquery.OutputHTML(button, true))
	// 	fmt.Println(htmlquery.InnerText(button))
	// }
	// println(len(buttons))


	var titles []string

	// Itere sobre os botões encontrados
	for _, button := range buttons {
		title := htmlquery.InnerText(button)
		titles = append(titles, title)
	}

	// response := strings.Join(titles, ", ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(titles)
}

func SearchReturnBase64(w http.ResponseWriter, r *http.Request) {
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
	url := fmt.Sprintf("https://www.myinstants.com%s", soundPath)
	println(url)
	base64String, err := getBase64Media(url)
	if err != nil {
		http.Error(w, "Erro ao obter o vídeo em base64", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"soundBase64": base64String,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
func getBase64Media(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Converte para base64
	base64String := base64.StdEncoding.EncodeToString(body)

	return base64String, nil
}
// func getBase64Media(url string) (string, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer response.Body.Close()

// 	var base64String string

// 	encoder := base64.NewEncoder(base64.StdEncoding, &base64String)
// 	defer encoder.Close()

// 	_, err = io.Copy(encoder, response.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	return base64String, nil
// }

func extractSoundPath(onclick string) string {
	format := strings.Replace(onclick, "play(", "", -1)
	splitParts := strings.Split(format, ",")

	return strings.Trim(splitParts[0], "'")
}
