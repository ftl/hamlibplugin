package plugin

import "fmt"

type GlobalSettings struct {
	Radios []RadioSettings `json:"radios"`
}

type RadioSettings struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func parseGlobalSettings(settings map[string]any) (GlobalSettings, error) {
	result := GlobalSettings{}

	radios, ok := settings["radios"].([]any)
	if !ok {
		return GlobalSettings{}, fmt.Errorf("global settings: unexpected radios type %T", settings["radios"])
	}

	for i := range radios {
		radio, ok := radios[i].(map[string]any)
		if !ok {
			return GlobalSettings{}, fmt.Errorf("global settings: unexpected radio type %T", radios[i])
		}
		name, ok := radio["name"].(string)
		if !ok {
			return GlobalSettings{}, fmt.Errorf("global settings: unexpected radio name type %T", radio["name"])
		}
		address, ok := radio["address"].(string)
		if !ok {
			return GlobalSettings{}, fmt.Errorf("global settings: unexpected radio address type %T", radio["address"])
		}
		result.Radios = append(result.Radios, RadioSettings{Name: name, Address: address})
	}

	return result, nil
}
