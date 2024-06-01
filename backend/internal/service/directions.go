package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/emeraldls/fyp/internal/types"
)

type RouteService struct {
	apiKey string
}

func meshOrigin(origin []float64) string {
	url := "&origin="
	for i := 0; i < len(origin); i++ {
		url += fmt.Sprintf("%f", origin[i])
		if i != len(origin)-1 {
			url += ","
		}
	}

	return url
}

func meshDestination(destination []float64) string {
	url := "&destination="
	for i := 0; i < len(destination); i++ {
		url += fmt.Sprintf("%f", destination[i])
		if i != len(destination)-1 {
			url += ","
		}
	}

	return url
}

// A route calculation consists of a single GET request. The only required parameters are an origin and a destination, given by two pairs of WGS84 coordinates in the form <latitude>,<longitude>; and a transportation mode, which can be bicycle, bus, car,pedestrian,scooter, taxi, or truck.
func (rs RouteService) GetRoute(transportationMode types.TransportMode, origin, destination []float64) (*types.Routes, error) {
	url := fmt.Sprintf("https://router.hereapi.com/v8/routes?transportMode=%s%s%s", transportationMode, meshOrigin(origin), meshDestination(destination))

	url += "&return=summary&apikey=" + rs.apiKey

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil
		}

		return nil, errors.New(string(b))
	}

	var routes types.Routes
	err = json.NewDecoder(res.Body).Decode(&routes)
	if err != nil {
		return nil, err
	}

	return &routes, nil
}

func (rs RouteService) CalculateRoute(transportationMode types.TransportMode, origin, destination []float64) (*types.Routes, error) {
	return nil, nil
}
