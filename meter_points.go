package octopus

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type GSP struct {
	Gsp, Mpan    string
	ProfileClass int `mapstructure:"profile_class"`
}

// Consumption represents a half hour consumption of gas or electricity.
type Consumption struct {
	// Consumption in kWh
	Consumption   float32
	IntervalStart time.Time `mapstructure:"interval_start"`
	IntervalEnd   time.Time `mapstructure:"interval_end"`
}

type ConsumptionRequest struct {
	PeriodFrom time.Time `url:"period_from,omitempty"`
	PeriodTo   time.Time `url:"period_to,omitempty"`
	PageSize   int       `url:"page_size,omitempty"`
	OrderBy    string    `url:"order_by,omitempty"`
	GroupBy    string    `url:"group_by,omitempty"`
}

// ElectricityMeterPoint returns the GSP for an electricity meter.
func (c *Client) ElectricityMeterPoint(mpan string) (*GSP, error) {
	uri := fmt.Sprintf("/electricity-meter-points/%s", mpan)

	resp, err := c.request("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	var gsp GSP

	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &gsp,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(resp)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal response into GSP struct: %w",
			err,
		)
	}

	return &gsp, nil
}

// consumption collates the common logic between the gas and electricity consumption
// functions.
// meterPointNumber is the MPRN or MPAN of the meter.
// electricity determines if we query the electricity meter endpoint or the gas one.
func (c *Client) consumption(
	meterPointNumber, serial string,
	options ConsumptionRequest,
	electricity bool,
) ([]Consumption, error) {
	var baseURL string
	if electricity {
		baseURL += "/electricity-meter-points"
	} else {
		baseURL += "/gas-meter-points"
	}

	uri := fmt.Sprintf(
		"%s/%s/meters/%s/consumption/",
		baseURL,
		meterPointNumber,
		serial,
	)

	resp, err := c.request("GET", uri, nil, options)
	if err != nil {
		return nil, err
	}

	var consumption []Consumption

	decoderConfig := mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &consumption,
		DecodeHook: mapstructure.StringToTimeHookFunc(
			"2006-01-02T15:04:05Z07:00",
		),
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(resp.(listResponse).Results)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal response into slice of Consumption struct: %w",
			err,
		)
	}

	return consumption, nil
}

// ElectricityConsumption returns the consumption of an electricity meter.
// mpan and serial are the MPAN and serial of the electricity meter.
// options are the parameters to use when querying for consumption.
func (c *Client) ElectricityConsumption(
	mpan, serial string,
	options ConsumptionRequest,
) ([]Consumption, error) {
	request, err := c.consumption(mpan, serial, options, true)

	return request, err
}

// GasConsumption returns the consumption of a gas meter.
// mprn and serial are the MPRN and serial of the gas meter.
// options are the parameters to use when querying for consumption.
func (c *Client) GasConsumption(
	mprn, serial string,
	options ConsumptionRequest,
) ([]Consumption, error) {
	request, err := c.consumption(mprn, serial, options, false)

	return request, err
}
