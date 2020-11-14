package octopus

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mitchellh/mapstructure"
)

type testResponse struct {
	Hello string
}

// Perform a straightforward get request and parsing of the results.
func TestClientBasicRequest(t *testing.T) {
	client := NewClient("HelloApiKey")

	httpmock.ActivateNonDefault(client.HTTPClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.octopus.energy/v1/test-api",
		httpmock.NewJsonResponderOrPanic(200, &testResponse{Hello: "there"}),
	)

	request, err := client.request("GET", "test-api", nil, nil)
	if err != nil {
		t.Error(err)
	}

	var response testResponse
	err = mapstructure.Decode(request, &response)
	if err != nil {
		t.Error(err)
	}

	if response.Hello != "there" {
		t.Errorf("Expected %s, got %s", "there", response.Hello)
	}
}
