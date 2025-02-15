package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const govAPI = "https://api.portaldatransparencia.gov.br/api-de-dados/safra-codigo-por-cpf-ou-nis"
const apiKey = "ba38a7dcb1e17354204bdceaea765c5f" // Substitua pela sua chave de API

func fetchGarantiaSafra(codigo string, pagina string) ([]byte, error) {
	url := fmt.Sprintf("%s?codigo=%s&pagina=%s", govAPI, codigo, pagina)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("chave-api-dados", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	codigo := query.Get("codigo")
	pagina := query.Get("pagina")
	if codigo == "" || pagina == "" {
		http.Error(w, "Parâmetros 'codigo' e 'pagina' são obrigatórios", http.StatusBadRequest)
		return
	}

	data, err := fetchGarantiaSafra(codigo, pagina)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar dados: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/garantia-safra", handler)
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
