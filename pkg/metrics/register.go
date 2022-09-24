package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

//nolint:errorlint
func Register(labels prometheus.Labels, reg prometheus.Registerer, collectors ...prometheus.Collector) {
	for _, c := range collectors {
		err := prometheus.WrapRegistererWith(labels, reg).Register(c)
		if err != nil {
			if _, ok := err.(*prometheus.AlreadyRegisteredError); !ok {
				log.WithError(err).
					Error("failed to register job duration metrics with prometheus")
			}
		}
	}
}
