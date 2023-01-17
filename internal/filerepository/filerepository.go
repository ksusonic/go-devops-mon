package filerepository

import (
	"bufio"
	"encoding/json"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"os"
)

type FileRepository struct {
	file *os.File
}

func NewFileRepository(filename string, restoreContent bool) (*FileRepository, *[]metrics.Metrics, error) {
	var restoredMetrics *[]metrics.Metrics = nil
	flag := os.O_WRONLY | os.O_APPEND
	if restoreContent {
		restoreMetrics, err := restoreMetrics(filename)
		if err != nil {
			return nil, nil, err
		}
		restoredMetrics = &restoreMetrics
	} else {
		flag |= os.O_CREATE
	}

	file, err := os.OpenFile(filename, flag, 0777)
	if err != nil {
		return nil, restoredMetrics, err
	}

	return &FileRepository{file: file}, restoredMetrics, nil
}

func (p *FileRepository) SaveMetrics(metrics []metrics.Metrics) error {
	log.Printf("Saving %d metrics\n", len(metrics))

	// clear old metrics
	err := p.file.Truncate(0)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		data, err := json.Marshal(&m)
		if err != nil {
			return err
		}
		data = append(data, '\n')
		_, err = p.file.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *FileRepository) Close() error {
	return p.file.Close()
}

func restoreMetrics(filename string) ([]metrics.Metrics, error) {
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
