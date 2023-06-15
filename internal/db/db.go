package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ksusonic/go-devops-mon/internal/metrics"
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

func (d DB) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d DB) Close() error {
	return d.db.Close()
}

func (d DB) SetMetric(ctx context.Context, m *metrics.Metric) (_ *metrics.Metric, err error) {
	switch m.Type {
	case metrics.CounterType:
		current, err := d.GetMetric(ctx, m.Type, m.ID)
		if err == nil {
			m.Delta = current.Delta
		}
		_, err = d.db.ExecContext(ctx,
			`INSERT INTO metrics
			(id, type, delta)
			VALUES ($1, $2, $3)
			ON CONFLICT(id, type) DO UPDATE SET delta=$3`,
			m.ID,
			m.Type,
			m.GetDelta(),
		)
		if err != nil {
			return nil, fmt.Errorf("error in SetMetric: %v", err)
		}
	case metrics.GaugeType:
		_, err = d.db.ExecContext(
			ctx,
			`INSERT INTO metrics
			(id, type, value)
			VALUES ($1, $2, $3)
			ON CONFLICT(id, type) DO UPDATE SET value=$3`,
			m.ID,
			m.Type,
			m.GetValue(),
		)
	default:
		return nil, fmt.Errorf("unknown metric type for db-insertion: %s", m.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("error in SetMetric: %v", err)
	}

	// get metric to check real value
	metric, err := d.GetMetric(ctx, m.Type, m.ID)
	if err != nil {
		return nil, fmt.Errorf("error in SetMetric: %v", err)
	}
	return metric, nil
}

func (d DB) SetMetrics(ctx context.Context, m []*metrics.Metric) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	gaugeStmt, err := tx.Prepare(`INSERT INTO metrics (id, type, value) VALUES($1, 'gauge', $2) ON CONFLICT(id, type) DO UPDATE SET value=$2;`)
	if err != nil {
		return err
	}
	counterStmt, err := tx.Prepare(`INSERT INTO metrics (id, type, delta) VALUES($1, 'counter', $2) ON CONFLICT(id, type) DO UPDATE SET delta=$2;`)
	if err != nil {
		return err
	}
	defer gaugeStmt.Close()
	defer counterStmt.Close()

	for _, metric := range m {
		switch metric.Type {
		case metrics.GaugeType:
			_, err = gaugeStmt.ExecContext(ctx, metric.ID, metric.GetValue())
		case metrics.CounterType:
			metricRow := tx.QueryRowContext(ctx, "SELECT id, type, value, delta FROM metrics WHERE type = $1 AND id = $2 LIMIT 1;", metric.Type, metric.ID)
			currentMetric, err2 := rowToMetric(metricRow)
			if err2 == nil {
				// plus current delta
				metric.Delta = currentMetric.Delta
			}
			_, err = counterStmt.ExecContext(ctx, metric.ID, metric.GetDelta())
		default:
			err = fmt.Errorf("unknown metric type: %s", metric.Type)
		}
		if err != nil {
			if err = tx.Rollback(); err != nil {
				return fmt.Errorf("SetMetrics: unable to rollback: %v", err)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("SetMetrics: unable to commit: %v", err)
	}

	return nil
}

func (d DB) GetMetric(ctx context.Context, type_ string, name string) (*metrics.Metric, error) {
	row := d.db.QueryRowContext(
		ctx,
		"SELECT id, type, value, delta FROM metrics WHERE type = $1 AND id = $2 LIMIT 1;",
		type_,
		name,
	)
	return rowToMetric(row)
}

func rowToMetric(row *sql.Row) (*metrics.Metric, error) {
	res := metrics.Metric{}
	var gaugeValue sql.NullFloat64
	var counterValue sql.NullInt64
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&res.ID, &res.Type, &gaugeValue, &counterValue)
	if err != nil {
		return nil, fmt.Errorf("error in rowToMetric: %v", err)
	}
	if gaugeValue.Valid {
		res.Value = &gaugeValue.Float64
	}
	if counterValue.Valid {
		res.Delta = &counterValue.Int64
	}

	return &res, nil
}

func (d DB) GetAllMetrics(ctx context.Context) (res []metrics.Metric, err error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, type, value, delta FROM metrics;")
	if err != nil {
		return nil, fmt.Errorf("error in GetAllMetrics: %v", err)
	}

	for rows.Next() {
		var m metrics.Metric
		var gaugeValue sql.NullFloat64
		var counterValue sql.NullInt64
		err = rows.Scan(&m.ID, &m.Type, &gaugeValue, &counterValue)
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
		_, ok := res[m.Type]
		if !ok {
			res[m.Type] = make(map[string]interface{})
		}
		if m.Type == metrics.GaugeType {
			res[m.Type][m.ID] = m.GetValue()
		} else if m.Type == metrics.CounterType {
			res[m.Type][m.ID] = m.GetDelta()
		}
	}
	return res, nil
}
