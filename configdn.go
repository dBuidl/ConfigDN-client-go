package configdn

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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

var ErrRefreshConfig = errors.New("error refreshing config")

const version = "0.0.2"

func (c ConfigDN) RefreshConfig(errorOnFail bool) error {
	client := &http.Client{}
	urlPath, err := url.JoinPath(c.Settings.Endpoint, "api/custom/v1/get_config/")

	if err != nil {
		if errorOnFail {
			return ErrRefreshConfig
		}
		return nil
	}

	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Add("ConfigDN-Client-Version", "ConfigDN-Go/"+version)
	req.Header.Add("Authorization", c.Settings.AuthKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		if errorOnFail {
			return ErrRefreshConfig
		}
		return nil
	}
	var decodedResponse = InitialResponse{}
	err = json.NewDecoder(resp.Body).Decode(&decodedResponse)
	if !decodedResponse.S || err != nil {
		if errorOnFail {
			return ErrRefreshConfig
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
	if c.LastUpdateTime+c.Settings.RefreshInterval > int(time.Now().Unix()) {
		defer c.RefreshConfig(false)
	}

	val, ok := c.FetchedConfig[key]
	if ok {
		return val.V
	}
	return nil
}

func (c ConfigDN) Get(key string) any {
	if c.LastUpdateTime+c.Settings.RefreshInterval > int(time.Now().Unix()) {
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
