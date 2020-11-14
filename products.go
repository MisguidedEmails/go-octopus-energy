package octopus

import (
	"fmt"
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

type ProductRequest struct {
	// Specify what products to receive from the API. Note: This is OR not AND, so doing
	// both business and green will show all green and all business products.
	variable, green, tracker, prepay, business bool
}

func (c *Client) ProductList(
	page int,
	productOptions *ProductRequest,
) (*[]Product, error) {
	// TODO: Add support for `available_at` as time
	// TODO: Possibly implement a generator/iterator?
	params := map[string]string{
		"page":        strconv.Itoa(page),
		"is_variable": strconv.FormatBool(productOptions.variable),
		"is_green":    strconv.FormatBool(productOptions.green),
		"is_tracker":  strconv.FormatBool(productOptions.tracker),
		"is_prepay":   strconv.FormatBool(productOptions.prepay),
		"is_business": strconv.FormatBool(productOptions.business),
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

// Get a specific product by it's productCode.
func (c *Client) Product(productCode string) (*Product, error) {
	// TODO: Add support for `tarrifs_active_at`

	uri := fmt.Sprintf("/products/%s", productCode)

	resp, err := c.request("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	var product Product

	err = mapstructure.Decode(resp, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
