package near

func (p *Platform) GetVersion() (string, error) {
	return p.client.GetVersion()
}
