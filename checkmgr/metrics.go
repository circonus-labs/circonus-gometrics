package checkmgr

import (
	"github.com/circonus-labs/circonus-gometrics/api"
)

// is this metric currently active
func (cm *CheckManager) IsMetricActive(name string) bool {
	_, ok := cm.activeMetrics[name]
	return ok
}

// Add new metrics to an existing check
func (cm *CheckManager) AddNewMetrics(newMetrics map[string]*api.CheckBundleMetric) {
	// only if check manager is enabled
	if !cm.enabled {
		return
	}

	// only if checkBundle has been populated
	if cm.checkBundle == nil {
		return
	}

	newCheckBundle := cm.checkBundle
	numCurrMetrics := len(newCheckBundle.Metrics)
	numNewMetrics := len(newMetrics)

	if numCurrMetrics+numNewMetrics >= cap(newCheckBundle.Metrics) {
		nm := make([]api.CheckBundleMetric, numCurrMetrics+numNewMetrics)
		copy(nm, newCheckBundle.Metrics)
		newCheckBundle.Metrics = nm
	}

	newCheckBundle.Metrics = newCheckBundle.Metrics[0 : numCurrMetrics+numNewMetrics]

	i := 0
	for _, metric := range newMetrics {
		newCheckBundle.Metrics[numCurrMetrics+i] = *metric
		i++
	}

	checkBundle, err := cm.apih.UpdateCheckBundle(newCheckBundle)
	if err != nil {
		cm.Log.Printf("[ERROR] updating check bundle with new metrics %v", err)
		return
	}

	for _, metric := range newMetrics {
		cm.activeMetrics[metric.Name] = true
	}

	cm.checkBundle = checkBundle
}
