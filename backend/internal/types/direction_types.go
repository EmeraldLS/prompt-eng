package types

type Entity struct {
	Locations []string `json:"locations"`
	From      string   `json:"from"`
	To        string   `json:"to"`
}

type Routes struct {
	Routes  []Route  `json:"routes"`
	Notices []Notice `json:"notices,omitempty"`
}

type Route struct {
	ID       string    `json:"id"`
	Sections []Section `json:"sections"`
}

type Section struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Departure Departure `json:"departure"`
	Arrival   Arrival   `json:"arrival"`
	Summary   Summary   `json:"summary"`
	Transport Transport `json:"transport"`
}

type Departure struct {
	Place Place `json:"place"`
}

type Arrival struct {
	Place Place `json:"place"`
}

type Place struct {
	Type             string   `json:"type"`
	Location         Location `json:"location"`
	OriginalLocation Location `json:"originalLocation"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Summary struct {
	Duration     int `json:"duration"`
	Length       int `json:"length"`
	BaseDuration int `json:"baseDuration"`
}

type Transport struct {
	Mode string `json:"mode"`
}

type Notice struct {
	Title string `json:"title"`
	Code  string `json:"code"`
}

type CalculateRouteRequest struct {
}

type GeoEndcodeResponse struct {
	Items []Item `json:"items"`
}

// Item represents each item in the list
type Item struct {
	Title           string        `json:"title"`
	ID              string        `json:"id"`
	ResultType      string        `json:"resultType"`
	HouseNumberType string        `json:"houseNumberType"`
	Address         Address       `json:"address"`
	Position        Coordinates   `json:"position"`
	Access          []Coordinates `json:"access"`
	MapView         MapView       `json:"mapView"`
	Scoring         Scoring       `json:"scoring"`
}

// Address represents the address details
type Address struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	StateCode   string `json:"stateCode"`
	State       string `json:"state"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	HouseNumber string `json:"houseNumber"`
}

// Coordinates represents the geographical coordinates
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// MapView represents the map view details
type MapView struct {
	West  float64 `json:"west"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	North float64 `json:"north"`
}

// Scoring represents the scoring details
type Scoring struct {
	QueryScore float64    `json:"queryScore"`
	FieldScore FieldScore `json:"fieldScore"`
}

// FieldScore represents the field score details
type FieldScore struct {
	City        float64   `json:"city"`
	Streets     []float64 `json:"streets"`
	HouseNumber float64   `json:"houseNumber"`
}
