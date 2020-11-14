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
	requestURL := fmt.Sprintf("%s%s", "https://api.octopus.energy/v1/", path)

	var err interface{}

	var request *resty.Response

	var contents interface{}
	baseRequest := c.
		HTTPClient.
		R().
		ForceContentType("application/json").
		SetResult(&contents).
		SetBody(body).
		SetQueryParams(parameters)

	switch method {
	case "GET":
		request, err = baseRequest.Get(requestURL)
	case "POST":
		request, err = baseRequest.Post(requestURL)
	case "DELETE":
		request, err = baseRequest.Delete(requestURL)
	case "PATCH":
		request, err = baseRequest.Patch(requestURL)
	case "PUT":
		request, err = baseRequest.Put(requestURL)
	default:
		return nil, &HTTPError{
			msg: fmt.Sprintf("given HTTP method '%s' not implemented", method),
		}
	}

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
	client := &Client{
		apiKey:     apiKey,
		HTTPClient: resty.New(),
	}

	client.HTTPClient.SetHeader("Accept", "application/json")
	client.HTTPClient.SetAuthToken(client.apiKey)

	return client
}
