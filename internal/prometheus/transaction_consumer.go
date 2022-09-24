package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (p *Prometheus) RegisterTrasactionConsumerMetrics() {
	p.parsedTxs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "parsed_txs"),
			Help: "The amount of parsed txs",
		},
		[]string{labelChain},
	)
}

func (p *Prometheus) SetParsedTxs(chain string, amount int) {
	p.parsedTxs.With(prometheus.Labels{
		labelChain: chain,
	}).Set(float64(amount))
}
