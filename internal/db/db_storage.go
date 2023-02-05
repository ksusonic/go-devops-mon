package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func NewDB(connectString string) (*DB, error) {
	ctx := context.Background()
	db, err := sql.Open("pgx", connectString)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS metrics
		(
			id    VARCHAR NOT NULL,
			type  VARCHAR NOT NULL,
			value DOUBLE PRECISION,
			delta BIGINT,
			hash  varchar,
			UNIQUE (id, type)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (d DB) Close() error {
	return d.db.Close()
}

func (d DB) SetMetric(ctx context.Context, m metrics.Metrics, h metrics.HashService) (_ metrics.Metrics, err error) {
	switch m.MType {
	case metrics.CounterMType:
		current, err := d.GetMetric(ctx, m.MType, m.ID)
		if err != nil {
			_, err = d.db.ExecContext(ctx,
				`INSERT INTO metrics
			(id, type, delta, hash)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT(id, type) DO UPDATE SET delta=$3, hash=$4`,
				m.ID,
				m.MType,
				*m.Delta,
				m.Hash,
			)
			if err != nil {
				return metrics.Metrics{}, fmt.Errorf("error in SetMetric: %v", err)
			}
		} else {
			*m.Delta += *current.Delta
			if err = h.SetHash(&m); err != nil {
				return metrics.Metrics{}, err
			}
			_, err = d.db.ExecContext(ctx,
				`INSERT INTO metrics
			(id, type, delta, hash)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT(id, type) DO UPDATE SET delta=$3, hash=$4`,
				m.ID,
				m.MType,
				*m.Delta,
				m.Hash,
			)
			if err != nil {
				return metrics.Metrics{}, fmt.Errorf("error in SetMetric: %v", err)
			}
		}
	case metrics.GaugeMType:
		_, err = d.db.ExecContext(
			ctx,
			`INSERT INTO metrics
			(id, type, value, hash)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT(id, type) DO UPDATE SET value=$3, hash=$4`,
			m.ID,
			m.MType,
			*m.Value,
			m.Hash,
		)
	default:
		return metrics.Metrics{}, fmt.Errorf("unknown metric type for db-insertion: %s", m.MType)
	}
	if err != nil {
		return metrics.Metrics{}, fmt.Errorf("error in SetMetric: %v", err)
	}

	// get metric to check real value
	metric, err := d.GetMetric(ctx, m.MType, m.ID)
	if err != nil {
		return metrics.Metrics{}, fmt.Errorf("error in SetMetric: %v", err)
	}
	return metric, nil
}

func (d DB) GetMetric(ctx context.Context, type_, name string) (res metrics.Metrics, err error) {
	row := d.db.QueryRowContext(
		ctx,
		"SELECT id, type, value, delta, hash FROM metrics WHERE type = $1 AND id = $2;",
		type_,
		name,
	)
	var gaugeValue sql.NullFloat64
	var counterValue sql.NullInt64
	err = row.Scan(&res.ID, &res.MType, &gaugeValue, &counterValue, &res.Hash)
	if err != nil {
		return metrics.Metrics{}, fmt.Errorf("error in GetMetric: %v", err)
	}

	if gaugeValue.Valid {
		res.Value = &gaugeValue.Float64
	}
	if counterValue.Valid {
		res.Delta = &counterValue.Int64
	}

	return res, nil
}

func (d DB) GetAllMetrics(ctx context.Context) (res []metrics.Metrics, err error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, type, value, delta, hash FROM metrics;")
	if err != nil {
		return nil, fmt.Errorf("error in GetAllMetrics: %v", err)
	}

	for rows.Next() {
		var m metrics.Metrics
		var gaugeValue sql.NullFloat64
		var counterValue sql.NullInt64
		err = rows.Scan(&m.ID, &m.MType, &gaugeValue, &counterValue, &m.Hash)
		if err != nil {
			return nil, fmt.Errorf("error in GetAllMetrics: %v", err)
		}
		if gaugeValue.Valid {
			m.Value = &gaugeValue.Float64
		}
		if counterValue.Valid {
			m.Delta = &counterValue.Int64
		}

		res = append(res, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error in GetAllMetrics: %v", err)
	}
	return res, nil
}

func (d DB) GetMappedByTypeAndNameMetrics(ctx context.Context) (map[string]map[string]interface{}, error) {
	res := make(map[string]map[string]interface{})
	allMetrics, err := d.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in GetMappedByTypeAndNameMetrics: %v", err)
	}

	for _, m := range allMetrics {
		_, ok := res[m.MType]
		if !ok {
			res[m.MType] = make(map[string]interface{})
		}
		if m.MType == metrics.GaugeMType {
			res[m.MType][m.ID] = *m.Value
		} else if m.MType == metrics.CounterMType {
			res[m.MType][m.ID] = *m.Delta
		}
	}
	return res, nil
}
