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

type Consumption struct {
	Consumption   float32
	IntervalStart time.Time `mapstructure:"interval_start"`
	IntervalEnd   time.Time `mapstructure:"interval_end"`
}

type ConsumptionRequest struct {
	// TODO: Change to native types (int, time)
	PeriodFrom, PeriodTo, PageSize, OrderBy, GroupBy string
}

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
		return nil, err
	}

	return &gsp, nil
}

func (c *Client) ElectricityConsumption(
	mpan, serial string,
	options ConsumptionRequest,
) ([]Consumption, error) {
	uri := fmt.Sprintf(
		"/electricity-meter-points/%s/meters/%s/consumption/",
		mpan,
		serial,
	)

	params := map[string]string{
		"period_from": options.PeriodFrom,
		"period_to":   options.PeriodTo,
		"page_size":   options.PageSize,
		"order_by":    options.OrderBy,
		"group_by":    options.GroupBy,
	}

	for key, value := range params {
		if value == "" {
			delete(params, key)
		}
	}

	resp, err := c.request("GET", uri, nil, params)
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
		return nil, err
	}

	return consumption, nil
}
