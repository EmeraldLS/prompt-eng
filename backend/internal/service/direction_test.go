package service

import (
	"encoding/json"
	"testing"

	"github.com/emeraldls/fyp/internal/types"
)

var rs = RouteService{
	apiKey: "IWkHU9joQd9Fle_Tz5-izygLSrLvycV638LY-28ic_0",
}

func TestRoutingService(t *testing.T) {
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
