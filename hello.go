package main

import (
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
//	"github.com/signalfx/splunk-otel-go/distro"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
//	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/attribute"


	)

// Define the struct for the response body. Used to parse json response to extract access token
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

// Function used to get value from the .env file using key
func goDotEnvVariable(key string) string {
	return os.Getenv(key)
}

// Get clien_id and client_secret from the .env file 
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

// Create request body. Sets client_id, client_secret, grant type and returns body
func createRequestBody(client_id, client_secret string) io.Reader {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)

	encodedData := strings.NewReader(data.Encode())
	return encodedData
}

// Uses client, url, and body to make request. Returns http response and err
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

// Convert http response to json using TokenResponse struct then extract access token value
func processTokenResponse(response io.Reader) string {
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

// Use above functions to get API token. Returns access token as a key and client for further API calls
func GetAPIkey(ctx context.Context) (string, *http.Client) {

/* Sample custom span instrumentation
	// Create a named tracer
	tracer := otel.Tracer("call/GetAPIkey")

	// crate a span with custom attributes
	var span trace.Span
	ctx, span = tracer.Start(ctx, "get API key")
	defer span.End()
*/

// Custom attribute attempt
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("myBool", true))
	defer span.End()

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
	accessToken := processTokenResponse(resp.Body)

	return accessToken, client
}


