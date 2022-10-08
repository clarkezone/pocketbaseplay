package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type authresponse struct {
	Admin string `json:"-"`
	Token string
}

func main() {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email":"", "password":""}`).
		Post("http://clarkezonedevbox3-tr:9099/api/admins/auth-via-email")

	dresp := authresponse{}
	err = json.Unmarshal(resp.Body(), &dresp)
	if err != nil {
		panic(err)
	}
	fmt.Println(dresp.Token)

	resp, err = client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Admin "+dresp.Token).
		Get("http://clarkezonedevbox3-tr:9099/api/collections/urls/records")
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
}
