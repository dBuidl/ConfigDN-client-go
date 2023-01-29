package configdn

import "errors"

type Settings struct {
	Endpoint        string
	RefreshInterval int
	AuthKey         string
}

func NewSettings(authKey string, endpoint string, refreshInterval int) (*Settings, error) {
	if refreshInterval <= 0 {
		return nil, errors.New("refresh interval should be a positive integer")
	}
	newSettings := Settings{
		AuthKey:         authKey,
		Endpoint:        endpoint,
		RefreshInterval: refreshInterval,
	}
	return &newSettings, nil
}

func (s Settings) ChangeRefreshInterval(refreshInterval int) error {
	if refreshInterval <= 0 {
		return errors.New("refresh interval should be a positive integer")
	}
	s.RefreshInterval = refreshInterval
	return nil
}
