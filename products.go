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

type ProductDetailed struct {
	Product `mapstructure:",squash"`
	// TODO: time.time
	TarrifsActiveAt string `mapstructure:"tariffs_active_at"`

	SingleRegisterElecTariffs map[string]struct {
		DirectDebityMonthly singleTarrif `mapstructure:"direct_debit_monthly"`
	} `mapstructure:"single_register_electricity_tariffs"`

	DualRegisterElecTariffs map[string]struct {
		DirectDebityMonthly dualElecTarrif `mapstructure:"direct_debit_monthly"`
	} `mapstructure:"dual_register_electricity_tariffs"`

	SingleRegisterGasTariffs map[string]struct {
		DirectDebityMonthly singleTarrif `mapstructure:"direct_debit_monthly"`
	} `mapstructure:"single_register_gas_tariffs"`

	SampleQuotes map[string]struct {
		DirectDebityMonthly sampleQuotes `mapstructure:"direct_debit_monthly"`
	} `mapstructure:"sample_quotes"`

	SampleConsumption sampleConsumption `mapstructure:"sample_consumption"`
}

type tariffBase struct {
	Code string

	StandingChargeExcVat float32 `mapstructure:"standing_charge_exc_vat"`
	StandingChargeIncVat float32 `mapstructure:"standing_charge_inc_vat"`

	OnlineDiscountExcVat float32 `mapstructure:"online_discount_exc_vat"`
	OnlineDiscountIncVat float32 `mapstructure:"online_discount_inc_vat"`

	DualFuelDiscountExcVat float32 `mapstructure:"dual_fuel_discount_exc_vat"`
	DualFuelDiscountIncVat float32 `mapstructure:"dual_fuel_discount_inc_vat"`

	ExitFeesExcVat float32 `mapstructure:"exit_fees_exc_vat"`
	ExitFeesIncVat float32 `mapstructure:"exit_fees_inc_vat"`

	Links []ProductLink
}

type singleTarrif struct {
	tariffBase `mapstructure:",squash"`

	StandardUnitRateExcVat float32 `mapstructure:"standard_unit_rate_exc_vat"`
	StandardUnitRateIncVat float32 `mapstructure:"standard_unit_rate_inc_vat"`
}

type dualElecTarrif struct {
	tariffBase `mapstructure:",squash"`

	DayUnitRateExcVat float32 `mapstructure:"day_unit_rate_exc_vat"`
	DayUnitRateIncVat float32 `mapstructure:"day_unit_rate_inc_vat"`

	NightUnitRateExcVat float32 `mapstructure:"night_unit_rate_exc_vat"`
	NightUnitRateIncVat float32 `mapstructure:"night_unit_rate_inc_vat"`
}

type sampleQuotes struct {
	ElectricitySingleRate sampleQuote `mapstructure:"electricity_single_rate"`
	ElectricityDualRate   sampleQuote `mapstructure:"electricity_dual_rate"`
	DualFuelSingleRate    sampleQuote `mapstructure:"dual_fuel_single_rate"`
	DualFuelDualRate      sampleQuote `mapstructure:"dual_fuel_dual_rate"`
}

type sampleQuote struct {
	AnnualCostIncVat float32 `mapstructure:"annual_cost_inc_vat"`
	AnnualCostExcVat float32 `mapstructure:"annual_cost_exc_vat"`
}

type sampleConsumption struct {
	ElectricitySingleRate struct {
		Standard int `mapstructure:"electricity_standard"`
	} `mapstructure:"electricity_single_rate"`

	ElectricityDualRate struct {
		Day   int `mapstructure:"electricity_day"`
		Night int `mapstructure:"electricity_night"`
	} `mapstructure:"electricity_dual_rate"`

	DualFuelSingleRate struct {
		ElectricityStandard int `mapstructure:"electricity_standard"`
		Gas                 int `mapstructure:"gas_standard"`
	} `mapstructure:"dual_fuel_single_rate"`

	DualFuelDuelRate struct {
		ElectricityDay   int `mapstructure:"electricity_day"`
		ElectricityNight int `mapstructure:"electricity_night"`
		Gas              int `mapstructure:"gas_standard"`
	} `mapstructure:"dual_fuel_dual_rate"`
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
func (c *Client) Product(productCode string) (*ProductDetailed, error) {
	// TODO: Add support for `tarrifs_active_at`

	uri := fmt.Sprintf("/products/%s", productCode)

	resp, err := c.request("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	var product ProductDetailed

	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &product,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(resp)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
