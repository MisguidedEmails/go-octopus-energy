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

func (c *Client) request(
	method,
	path string,
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

	err = decoder.Decode(contents)
	if err != nil {
		return contents, nil
	}

	return results, nil
}

func NewClient(apiKey string) *Client {
	apiClient := &Client{
		apiKey:     apiKey,
		HTTPClient: resty.New(),
	}

	return apiClient
}
