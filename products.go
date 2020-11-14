package octopus

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Product struct {
	Code         string
	FullName     string `mapstructure:"full_name"`
	DisplayName  string `mapstructure:"display_name"`
	Description  string
	IsVariable   bool `mapstructure:"is_variable"`
	IsGreen      bool `mapstructure:"is_green"`
	IsTracker    bool `mapstructure:"is_tracker"`
	IsPrepay     bool `mapstructure:"is_prepay"`
	IsBusiness   bool `mapstructure:"is_business"`
	IsRestricted bool `mapstructure:"is_restricted"`
	Term         int
	Brand        string
	// TODO: We should probably translate this to `time.Time`
	AvailableFrom string `mapstructure:"available_from"`
	AvailableTo   string `mapstructure:"available_to"`
	Links         []ProductLink
}

type ProductLink struct {
	Href   string
	Method string
	Rel    string
}

func (c *Client) ProductList(page int) (*[]Product, error) {
	params := map[string]string{
		"page": strconv.Itoa(page),
	}

	resp, err := c.request("GET", "/products", nil, params)
	if err != nil {
		return nil, err
	}

	var products []Product

	err = mapstructure.Decode(resp.(listResponse).Results, &products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}
