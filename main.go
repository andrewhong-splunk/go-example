package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	
	"context"
	"github.com/signalfx/splunk-otel-go/distro"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

/* Query Params
type QueryParams struct {
	OriginLocationCode string
	DestinationLocationCode string
	DepartureDate string
	Adults string
	Max string
}
*/

// Create request
func createRequest(client *http.Client, apiUrl string, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Context attribute test
	fmt.Println("main 34 req.Context():")
	fmt.Println(req.Context())
	span := trace.SpanFromContext(req.Context())
	if span != nil {
		fmt.Println("main 37 span")
		fmt.Println(span)
		span.SetAttributes(attribute.String("custom_attribute", "test_value"))
	} else {
	fmt.Println("span is nil")
}


// "https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=LAX&destinationLocationCode=TPE&departureDate=2024-10-03&adults=1&max=2"
	q := req.URL.Query()
	q.Add("originLocationCode", "LAX")
	q.Add("destinationLocationCode","TPE")
	q.Add("departureDate","2024-10-03")
	q.Add("adults","1")
	q.Add("max","2")
	req.URL.RawQuery = q.Encode()

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


// If no Span is currently set in ctx an implementation of a Span that performs no operations is returned

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

	// Create tracer and  context
	ctx := context.Background()

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stringAttr", "MyAttribute"))
	fmt.Println("main 99 span:")
	fmt.Println(span)

	// Get token and client to make request. Set url
	token, client := GetAPIkey(ctx)
	client.Transport = otelhttp.NewTransport(http.DefaultTransport)

	url := "https://test.api.amadeus.com/v2/shopping/flight-offers"
	
	// Get response using client, url, and token
	resp, err := createRequest(client, url, token)

	// Process response data
	responseData, err := processResponse(resp) 

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(responseData)

	// Custom Metric Test

	meter := otel.Meter("Example")
	counter, err := meter.Int64Counter(
		"amadeus.price",
		metric.WithUnit("5"),
		metric.WithDescription("Description placeholder"),
	)
	if err != nil {
		log.Fatal(err)
	}

	counter.Add(ctx, 1)
}
