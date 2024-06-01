package service

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/emeraldls/fyp/internal/types"
	"github.com/joho/godotenv"
)

func TestRoutingService(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	var apiKey = os.Getenv("HERE_API_KEY")

	var rs = RouteService{
		apiKey: apiKey,
	}

	t.Run("get_route", func(t *testing.T) {
		res, err := rs.GetRoute(types.BICYCLE, []float64{53.32556, 14.65314}, []float64{53.65422, 14.66636})
		if err != nil {
			t.Error(err)
			return
		}
		jsonVal, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			t.Errorf("unable to marshal into json: %v\n", err)
			return
		}

		t.Logf(string(jsonVal))
	})

	t.Run("geo_encode", func(t *testing.T) {
		res, err := rs.EncodeText("Lagos, Nigeria")
		if err != nil {
			t.Error(err)
			return
		}

		jsonVal, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			t.Errorf("unable to marshal into json: %v\n", err)
			return
		}

		t.Logf(string(jsonVal))
	})
}
