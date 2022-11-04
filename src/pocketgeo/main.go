package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/go-resty/resty/v2"
	"github.com/valyala/fasthttp"
)

type authresponse struct {
	Admin string `json:"-"`
	Token string
}

type AutoGenerated struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Items      []struct {
		CollectionID   string `json:"@collectionId"`
		CollectionName string `json:"@collectionName"`
		Created        string `json:"created"`
		ID             string `json:"id"`
		LongURL        string `json:"longurl"`
		ShortURL       string `json:"shorturl"`
		Updated        string `json:"updated"`
	} `json:"items"`
}

type GeoStruct struct {
	Locations []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Timestamp          string   `json:"timestamp"`
			Altitude           int      `json:"altitude"`
			Speed              int      `json:"speed"`
			HorizontalAccuracy int      `json:"horizontal_accuracy"`
			VerticalAccuracy   int      `json:"vertical_accuracy"`
			Motion             []string `json:"motion"`
			Pauses             bool     `json:"pauses"`
			Activity           string   `json:"activity"`
			DesiredAccuracy    int      `json:"desired_accuracy"`
			Deferred           int      `json:"deferred"`
			SignificantChange  string   `json:"significant_change"`
			LocationsInPayload int      `json:"locations_in_payload"`
			DeviceID           string   `json:"device_id"`
			Wifi               string   `json:"wifi"`
			BatteryState       string   `json:"battery_state"`
			BatteryLevel       float64  `json:"battery_level"`
		} `json:"properties"`
	} `json:"locations"`
}

func NewGeoHandler() *router.Router {
	router := router.New()

	router.POST("/postgeo", geopost)
	return router
}

func geopost(ctx *fasthttp.RequestCtx) {
	dresp := GeoStruct{}
	err := json.Unmarshal(ctx.PostBody(), &dresp)
	if err != nil {
		log.Println("unable to unmarshal json %v", err)
	}
	log.Println("Got a geocoordinate %v", dresp.Locations[0])
}

func main() {
	user := os.Getenv("POCKET_SHORTEN_USERNAME")
	pass := os.Getenv("POCKET_SHORTEN_PASSWORD")
	url := os.Getenv("POCKET_DB_URL")

	client := resty.New()

	log.Printf("Geologger.")
	if url == "" {
		log.Printf("%c", errors.New("Environment variable POCKET_DB_URL not set"))
	} else {
		log.Printf("Database url = %v", url)
	}

	if user == "" {
		log.Printf("%v", errors.New("Environment variable POCKET_SHORTEN_USERNAME not set"))
	} else {
		log.Printf("Database un = %v", user)
	}

	if pass == "" {
		log.Printf("%v", errors.New("Environment variable POCKET_SHORTEN_PASSWORD not set"))
	} else {
		log.Printf("Database pw= %v", pass)
	}

	log.Printf("Login URL %v\n", url+"api/admins/auth-via-email")

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email":"` + user + `", "password":"` + pass + `"}`).
		Post(url + "api/admins/auth-via-email")

	if err != nil {
		log.Printf("%v", err)
	}
	if resp.StatusCode() != 200 {
		log.Printf("Authentication failed")
		log.Printf("%v", errors.New(string(resp.Body())))
	}

	dresp := authresponse{}
	err = json.Unmarshal(resp.Body(), &dresp)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Println("Pocketbase Authenticated.  Ready for responses")
	}

	router := NewGeoHandler()

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	// TODO test locally with CURL
	// TODO dockerize and deploy to cluster
	// TODO add cloudflare ingress
}
