package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/emeraldls/fyp/internal/types"
)

func (rs *RouteService) EncodeText(location string) (*types.GeoEndcodeResponse, error) {
	p := Parser(location)
	fmt.Println("Parser: ", p)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://geocode.search.hereapi.com/v1/geocode?q=%s&apiKey=%s", p, rs.apiKey), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var encodeRes types.GeoEndcodeResponse
	if res.StatusCode != http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil
		}

		fmt.Println("Response = ", string(b))

		return nil, errors.New(string(b))
	}

	err = json.NewDecoder(res.Body).Decode(&encodeRes)
	if err != nil {
		return nil, err
	}

	return &encodeRes, nil
}

func Parser(text string) string {
	text = strings.Replace(text, ",", "%2C", -1)
	text = strings.Replace(text, " ", "+", -1)
	return text
}
