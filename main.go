package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	
	"context"
	"github.com/signalfx/splunk-otel-go/distro"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)


func createRequest(client *http.Client, apiUrl string, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)	
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Process response data
func processResponse(resp *http.Response) (string, error) {
	responseData, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(responseData), err

}

func main() {
	// Otel instrumentation
	sdk, err := distro.Run()
	if err != nil {
		panic(err)
	}
	// Flush all spans before the application exits
	defer func() {
		if err := sdk.Shutdown(context.Background()); err != nil {
			 panic(err)
		}
	}()

	// Get token and client to make request. Set url
	token, client := GetAPIkey()
	client.Transport = otelhttp.NewTransport(http.DefaultTransport)

	url := "https://test.api.amadeus.com/v1/shopping/flight-destinations?origin=PAR&maxPrice=200"
	
	// Get response using client, url, and token
	resp, err := createRequest(client, url, token)

	// Process response data
	responseData, err := processResponse(resp) 

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(responseData)

}
