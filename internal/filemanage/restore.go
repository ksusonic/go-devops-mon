package filemanage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

func RestoreMetrics(filename string) ([]metrics.Metrics, error) {
	var result []metrics.Metrics
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var metric metrics.Metrics
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return nil, err
		}
		result = append(result, metric)
	}
	return result, nil
}
