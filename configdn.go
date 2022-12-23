package configdn

type ConfigDN struct {
	settings       settings
	lastUpdateTime int
	fetchedConfig  map[string]any
}

func (c ConfigDN) RefreshConfig() {

}

func (c ConfigDN) Get() any {

}

func (c ConfigDN) GetLocal() any {

}

func (c ConfigDN) ChangeRefreshInterval(refreshInterval int) error {
	err := c.settings.changeRefreshInterval(refreshInterval)
	return err
}
