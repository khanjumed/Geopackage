package services

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
)

type GeocodeResult struct {
    Lat float64
    Lng float64
}

func GeocodeAddress(address string) (GeocodeResult, error) {
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
    reqURL := fmt.Sprintf("%s?address=%s&key=%s", baseURL, url.QueryEscape(address), apiKey)

    resp, err := http.Get(reqURL)
    if err != nil {
        return GeocodeResult{}, err
    }
    defer resp.Body.Close()

    var result struct {
        Results []struct {
            Geometry struct {
                Location GeocodeResult `json:"location"`
            } `json:"geometry"`
        } `json:"results"`
        Status string `json:"status"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return GeocodeResult{}, err
    }

    if result.Status != "OK" || len(result.Results) == 0 {
        return GeocodeResult{}, fmt.Errorf("Google Maps API returned: %s", result.Status)
    }

    return result.Results[0].Geometry.Location, nil
}