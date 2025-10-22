package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// MUDANÇA 1: Adicionado "number" como parâmetro
func GetCoordsFromAddress(street, number, city, uf string) (float64, float64, error) {
	address := fmt.Sprintf("%s %s, %s, %s", street, number, city, uf)
	baseURL := "https://nominatim.openstreetmap.org/search"
	query := url.QueryEscape(address)
	fullURL := fmt.Sprintf("%s?q=%s&format=json&limit=1", baseURL, query)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// MUDANÇA 3: User-Agent formatado corretamente (use seu e-mail)
	// A política do Nominatim exige um User-Agent que identifique sua aplicação.
	req.Header.Set("User-Agent", "MedAssistApp/1.0 (komatsuhenry@gmail.com)")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao fazer requisição ao Nominatim: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Você provavelmente está recebendo um erro 403 (Forbidden) aqui por causa do User-Agent
		return 0, 0, fmt.Errorf("Nominatim retornou status: %s", resp.Status)
	}

	var results []NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, fmt.Errorf("erro ao decodificar resposta do Nominatim: %w", err)
	}

	if len(results) == 0 {
		// Este é o erro que você provavelmente está recebendo por falta do número
		return 0, 0, fmt.Errorf("nenhuma coordenada encontrada para o endereço: %s", address)
	}

	lat, errLat := strconv.ParseFloat(results[0].Lat, 64)
	lon, errLon := strconv.ParseFloat(results[0].Lon, 64)

	if errLat != nil || errLon != nil {
		return 0, 0, fmt.Errorf("erro ao converter coordenadas para float")
	}

	return lat, lon, nil
}
