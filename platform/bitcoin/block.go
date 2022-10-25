package bitcoin

func (p *Platform) GetCurrentBlockNumber() (int64, error) {
	return p.client.GetCurrentBlockNumber()
}

func (p *Platform) GetBlockByNumber(num int64) ([]byte, error) {
	return nil, nil
}
