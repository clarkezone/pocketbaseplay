package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"test.com/geo/pkg/pocketbase"
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

type PostGeo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"long"`
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
	router.GET("/", home)
	router.POST("/postgeo", geopost)
	return router
}

func home(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx.Response.BodyWriter(), "Hello!")
	ctx.Logger().Printf("Default")
}

func geopost(ctx *fasthttp.RequestCtx) {
	ctx.Logger().Printf("GeoPost")
	dresp := GeoStruct{}
	err := json.Unmarshal(ctx.PostBody(), &dresp)

	if err != nil {
		ctx.Logger().Printf("unable to unmarshal json %v", err)
	}
	ctx.Logger().Printf("Got a geocoordinate %v", dresp.Locations[0].Geometry.Coordinates[0])

	ctx.SuccessString("", "")

}

func main() {
	user := os.Getenv("POCKET_SHORTEN_USERNAME")
	pass := os.Getenv("POCKET_SHORTEN_PASSWORD")
	url := os.Getenv("POCKET_DB_URL")

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

	gp := PostGeo{}
	gp.Lat = 1.1
	gp.Lon = 1.2
	bytes, err := json.Marshal(gp)
	if err != nil {
		log.Printf("error writing %v", err)
	}

	writeclient := pocketbase.NewClient(url, user, pass)
	err = writeclient.Create("geolog", string(bytes))
	if err != nil {
		log.Printf("error writing %v", err)
	} else {
		log.Printf("Success")
	}

	//log.Println("Listening for geoevents")
	//router := NewGeoHandler()

	//log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	// TODO test locally with CURL
	// TODO dockerize and deploy to cluster
	// TODO add cloudflare ingress
}
