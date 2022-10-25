package binance

func (p *Platform) GetVersion() (string, error) {
	return p.dexClient.GetVersion()
}
