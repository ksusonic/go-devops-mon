package filemanage

import (
	"encoding/json"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"log"
	"os"
	"time"
)

type FileProducer struct {
	file     *os.File
	DropChan <-chan time.Time
}

func NewFileProducer(filename string, restoreContent bool) *FileProducer {
	flag := os.O_WRONLY | os.O_APPEND
	if !restoreContent {
		flag |= os.O_CREATE
	}
	file, err := os.OpenFile(filename, flag, 0777)
	if err != nil {
		log.Fatal(err)
	}

	return &FileProducer{file: file}
}

func (p *FileProducer) WriteMetrics(metrics []metrics.Metrics) error {
	log.Printf("Saving %d metrics\n", len(metrics))

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

func (p *FileProducer) Close() error {
	return p.file.Close()
}

func (p *FileProducer) DropRoutine(storage metrics.ServerMetricStorage, interval string) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(duration)
	for {
		_ = <-ticker.C
		err := p.WriteMetrics(storage.GetAllMetrics())
		if err != nil {
			log.Println("Error saving metrics: ", err)
		}
	}
}
