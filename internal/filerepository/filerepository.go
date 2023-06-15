package filerepository

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"go.uber.org/zap"
)

type FileRepository struct {
	file   *os.File
	logger *zap.Logger
}

func NewFileRepository(filename string, logger *zap.Logger) (*FileRepository, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &FileRepository{file: file, logger: logger}, nil
}

func (p *FileRepository) DropRoutine(ctx context.Context, getMetricsFunc func(context.Context) ([]metrics.Metric, error), duration time.Duration) {
	p.logger.Info("Started repository drop routine ", zap.String("destination", p.DebugInfo()), zap.Duration("drop-duration", duration))

	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			allMetrics, err := getMetricsFunc(ctx)
			if err != nil {
				p.logger.Error("Error while getting all metrics: ", zap.Error(err))
				continue
			}
			if err = p.SaveMetrics(allMetrics); err != nil {
				p.logger.Error("Error while saving metrics to repository: ", zap.Error(err))
			} else {
				p.logger.Debug("Saved metrics to repository")
			}
		case <-ctx.Done():
			p.logger.Info("Finished repository routine")
			return
		}
	}
}

func (p *FileRepository) ReadCurrentState() []*metrics.Metric {
	var result []*metrics.Metric

	scanner := bufio.NewScanner(p.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var metric metrics.Metric
		err := json.Unmarshal(data, &metric)
		if err == nil {
			result = append(result, &metric)
		}
		// ignore incorrect metric records
	}
	return result
}

func (p *FileRepository) SaveMetrics(metrics []metrics.Metric) error {
	// clear old metrics
	err := p.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.file.Seek(0, 0)
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

func (p *FileRepository) DebugInfo() string {
	return "file: " + p.file.Name()
}

func (p *FileRepository) Close() error {
	return p.file.Close()
}
