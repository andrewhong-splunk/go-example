package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"encoding/json"
	"io"

	"context"
	"github.com/signalfx/splunk-otel-go/distro"
//	"github.com/signalfx/splunk-otel-go/instrumentation/net/http/splunkhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"


)


type TokenResponse struct {
	Type string `json:"type"`
	Username string `json:"username"`
	ApplicationName string `json:"application_name"`
	ClientID string `json:"client_id"`
	TokenType string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	State string `json:"state"`
	Scope string `json:"scope"`
}

func goDotEnvVariable(key string) string {
	return os.Getenv(key)
}

func clientVariables() (string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	client_id := goDotEnvVariable("CLIENT_ID")
	client_secret := goDotEnvVariable("CLIENT_SECRET")
	
	if client_id == "" || client_secret == "" {
		log.Fatal("client_id or client_secret not found in environment")
	}

	return client_id, client_secret	
}

func createRequestBody(client_id, client_secret string) io.Reader {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)

	encodedData := strings.NewReader(data.Encode())
	return encodedData
}

func sendPostRequest(client *http.Client, apiUrl string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return nil, err	
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func processResponse(response io.Reader) string {
	responseData, err := ioutil.ReadAll(response)
        if err != nil {
                log.Fatal(err)
	}


        var tokenResponse TokenResponse
        err = json.Unmarshal(responseData, &tokenResponse)
        if err != nil {
                log.Fatalf("Error unmarshalling JSON: %v", err)
	}

        accessToken := tokenResponse.AccessToken
        return accessToken
}

func sendGetRequest(client *http.Client, apiUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiUrl)
	if err != nil {
		return nil, err
	}

	

func main() {
// otel instrumentation
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

// Load env variables	
	client_id, client_secret := clientVariables()

// Set url and token/secret	
	body := createRequestBody(client_id, client_secret)

// Create request and add header
	apiUrl := "https://test.api.amadeus.com/v1/security/oauth2/token"


// Create client and send request
client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := sendPostRequest(client, apiUrl, body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

// Read response body, get access token
	accessToken := processResponse(resp.Body)

// Make API call
	
}


