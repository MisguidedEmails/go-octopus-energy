// Package octopus implements a client for the Octopus energy REST API.
package octopus

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

type Client struct {
	HTTPClient *resty.Client

	apiKey string
}

type HTTPError struct {
	msg string
}

func (e *HTTPError) Error() string {
	return e.msg
}

type listResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []interface{}
}

// request sends a request to the Octopus API.
// method is the HTTP method to use, e.g "GET", "POST", etc
// path is the path after "/v1" of the URI to request.
// body to send to the API, and will auto marshal to JSON where possible.
// parameters are the query parameters to send.
// If the response can be parsed into a listResponse, then that will be returned,
// otherwise just the raw response interface{} will be returned.
func (c *Client) request(
	method, path string,
	body interface{},
	parameters map[string]string,
) (interface{}, error) {
	requestURL := fmt.Sprintf("%s%s", "https://api.octopus.energy/v1", path)

	var err interface{}

	var request *resty.Response

	var contents interface{}
	baseRequest := c.
		HTTPClient.
		R().
		ForceContentType("application/json").
		SetResult(&contents).
		SetBody(body).
		SetQueryParams(parameters).
		SetHeader("Accept", "application/json").
		SetBasicAuth(c.apiKey, "")

	request, err = baseRequest.Execute(method, requestURL)

	if !request.IsSuccess() {
		return request, &HTTPError{
			msg: fmt.Sprintf(
				"Error %d making %s request to '%s': Body: %+v. Resp: %s",
				request.StatusCode(),
				method,
				requestURL,
				body,
				request.String(),
			),
		}
	} else if err != nil {
		// If the request was "successful" but an error occurred
		return nil, fmt.Errorf("Error making request to '%s': %w", requestURL, err)
	}

	var results listResponse

	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &results,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		panic(err)
	}

	// Try to decode to ListResponse and if it fails, just bail and return the interface
	// instead.
	err = decoder.Decode(contents)
	if err != nil {
		return contents, nil
	}

	return results, nil
}

// NewClient creates a new Client with the given apiKey.
func NewClient(apiKey string) *Client {
	apiClient := &Client{
		apiKey:     apiKey,
		HTTPClient: resty.New(),
	}

	return apiClient
}
