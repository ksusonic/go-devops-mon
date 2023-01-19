package filerepository

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
)

type FileRepository struct {
	file *os.File
}

func NewFileRepository(filename string) (*FileRepository, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &FileRepository{file: file}, nil
}

func (p *FileRepository) ReadCurrentState() []metrics.Metrics {
	var result []metrics.Metrics

	scanner := bufio.NewScanner(p.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var metric metrics.Metrics
		err := json.Unmarshal(data, &metric)
		if err == nil {
			result = append(result, metric)
		}
		// ignore incorrect metric records
	}
	return result
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
