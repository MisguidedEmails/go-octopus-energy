# go-octopus-energy
A very basic and barebones Go client for the Octopus Energy API.

## Usage

You can get your API token, MPAN/MPRN and serial from the [Octopus Energy dev dashboard](https://octopus.energy/dashboard/developer/).
```golang
import (
    "fmt"

    "github.com/misguidedemails/go-octopus-energy"
)

func main() {
    client := octopus.NewClient("API TOKEN HERE")

    consumption, err := client.ElectricityConsumption(
        "MPAN",
        "ELEC SERIAL",
        octopus.ConsumptionRequest{PageSize: 10}
    )
    if err != nil {
        return nil, err
    }

    for _, consumption := range req {
        message := fmt.Sprintf(
            "%f kWh used between '%s' and '%s'",
            consumption.Consumption,
            consumption.IntervalStart,
            consumption.IntervalEnd,
        )

        fmt.Println(message)
    }
}
```
