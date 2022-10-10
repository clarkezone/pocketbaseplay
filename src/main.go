package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
		LongURL        string `json:"long_url"`
		ShortURL       string `json:"short_url"`
		Updated        string `json:"updated"`
	} `json:"items"`
}

type Service interface {
	Store(string, string) error
	Lookup(string) (string, error)
}

//TODO implement service using dictionary

type DictStore struct {
	m map[string]string
}

func (store *DictStore) Store(short string, long string) error {
	store.m[short] = long
	return nil
}

func (store *DictStore) Lookup(short string) (string, error) {
	val, pres := store.m[short]
	if pres {
		return val, nil
	} else {
		return "", errors.New("Key not found")
	}
}

func NewDictstore() Service {
	ds := &DictStore{}
	ds.m = make(map[string]string)
	return ds
}

type handler struct {
	schema  string
	host    string
	storage Service
}

func NewHandler(schema string, host string, storage Service) *router.Router {
	router := router.New()

	h := handler{schema, host, storage}
	router.GET("/{shortLink}", h.redirect)
	return router
}

func (h handler) redirect(ctx *fasthttp.RequestCtx) {
	code := ctx.UserValue("shortLink").(string)

	uri, err := h.storage.Lookup(code)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(http.StatusNotFound)
		return
	}

	ctx.Redirect(uri, http.StatusMovedPermanently)
}

func main() {
	user := os.Getenv("POCKET_SHORTEN_USERNAME")
	pass := os.Getenv("POCKET_SHORTEN_PASSWORD")

	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email":"` + user + `", "password":"` + pass + `"}`).
		Post("http://clarkezonedevbox3-tr:9099/api/admins/auth-via-email")

	if resp.StatusCode() != 200 {
		panic(errors.New(string(resp.Body())))
	}
	if err != nil {
		panic(err)
	}

	dresp := authresponse{}
	err = json.Unmarshal(resp.Body(), &dresp)
	if err != nil {
		panic(err)
	}

	resp, err = client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Admin "+dresp.Token).
		Get("http://clarkezonedevbox3-tr:9099/api/collections/urls/records")
	if resp.StatusCode() != 200 {
		panic(errors.New(string(resp.Body())))
	}
	if err != nil {
		panic(err)
	}

	dresp2 := AutoGenerated{}
	err = json.Unmarshal(resp.Body(), &dresp2)
	if err != nil {
		panic(err)
	}
	ds := NewDictstore()
	for _, element := range dresp2.Items {
		ds.Store(element.ShortURL, element.LongURL)
	}

	router := NewHandler("", "", ds)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	// TODO test locally with CURL
	// TODO dockerize and deploy to cluster
	// TODO add cloudflare ingress
}
