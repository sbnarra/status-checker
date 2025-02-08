package prometheus

import (
	"fmt"
	"status-checker/internal/checker"
	"status-checker/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type checkMetrics struct {
	total     prometheus.Counter
	success   prometheus.Counter
	recovered prometheus.Counter
	failed    prometheus.Counter
}

var metrics = make(map[string]checkMetrics)

func Publish(name string, result checker.Result) {
	if config.PrometheusEnabled {
		return
	}

	checkMetrics := getCheckMetrics(name)
	checkMetrics.total.Inc()

	if result.CheckError == nil {
		checkMetrics.success.Inc()
	} else if result.RecoverError != nil || result.RecheckError != nil {
		checkMetrics.failed.Inc()
	} else {
		checkMetrics.recovered.Inc()
	}
}

func getCheckMetrics(name string) checkMetrics {
	metric, ok := metrics[name]
	if !ok {
		metric = checkMetrics{
			total:     createMetric(name, "total"),
			success:   createMetric(name, "success"),
			recovered: createMetric(name, "recovered"),
			failed:    createMetric(name, "failed"),
		}
		metrics[name] = metric
	}
	return metric
}

func createMetric(name string, ctype string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("check_%s_%s", name, ctype),
		Help: fmt.Sprintf("The total number of %s checks", ctype),
	})
}
