package configdn

import "errors"

type settings struct {
	endpoint        string
	refreshInterval int
	authKey         string
}

func NewSettings(authKey string, endpoint string, refreshInterval int) (*settings, error) {
	if refreshInterval <= 0 {
		return nil, errors.New("refresh interval should be a posotive integer")
	}
	newSettings := settings{
		authKey:         authKey,
		endpoint:        endpoint,
		refreshInterval: refreshInterval,
	}
	return &newSettings, nil
}

func (s settings) changeRefreshInterval(refreshInterval int) error {
	if refreshInterval <= 0 {
		return errors.New("refresh interval should be a posotive integer")
	}
	s.refreshInterval = refreshInterval
	return nil
}
