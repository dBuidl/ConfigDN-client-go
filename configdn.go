package configdn

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type InitialResponse struct {
	S bool                   `json:"s"`
	D map[string]ConfigValue `json:"d"`
}

type ConfigValue struct {
	V any `json:"v"`
}

type ConfigDN struct {
	Settings       Settings
	LastUpdateTime int
	FetchedConfig  map[string]ConfigValue
}



func (c ConfigDN) RefreshConfig(errorOnFail bool) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", c.Settings.Endpoint + "public_api/v1/get_config/", nil)
	req.Header.Add("User-Agent", "ConfigDN-JS/0.0.1")
	req.Header.Add("Authorization", c.Settings.AuthKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		if errorOnFail {
			return errors.New("error refreshing config")
		}
		return nil
	}
	var decodedResponse = InitialResponse{}
	json.NewDecoder(resp.Body).Decode(&decodedResponse)
	if !decodedResponse.S {
		if errorOnFail {
			return errors.New("error refreshing config")
		}
		return nil
	}
	var finalMap = map[string]ConfigValue{}
	for key, value := range decodedResponse.D {
		finalMap[key] = value
	}
	c.FetchedConfig = finalMap
	c.LastUpdateTime = int(time.Now().Unix())
	return nil
}

func (c ConfigDN) GetLocal(key string) any {
	if c.LastUpdateTime + c.Settings.RefreshInterval > int(time.Now().Unix()) {
		defer c.RefreshConfig(false)
	}

	val, ok := c.FetchedConfig[key]
	if ok {
		return val.V
	}
	return nil
}

func (c ConfigDN) Get(key string) any {
	if c.LastUpdateTime + c.Settings.RefreshInterval > int(time.Now().Unix()) {
		c.RefreshConfig(false)
	}
	return c.GetLocal(key)
}

func (c ConfigDN) ChangeRefreshInterval(refreshInterval int) error {
	err := c.Settings.ChangeRefreshInterval(refreshInterval)
	return err
}

func NewConfigDN(authKey string) ConfigDN {
	return NewCustomConfigDN(authKey, "https://cdn.configdn.com/", 60)
}

func NewCustomConfigDN(authKey string, apiEndpoint string, refreshInterval int) ConfigDN {
	s, _ := NewSettings(authKey, apiEndpoint, refreshInterval)
	return ConfigDN{
		*s,
		int(time.Now().Unix()),
		map[string]ConfigValue{},
	}
}
