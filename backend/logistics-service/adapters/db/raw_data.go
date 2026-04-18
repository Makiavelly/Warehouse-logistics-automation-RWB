package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) InsertRawData(ctx context.Context, points []models.RawDataPoint) (int, error) {
	const op = "db.InsertRawData"

	if len(points) == 0 {
		return 0, nil
	}

	// Build bulk insert
	placeholders := make([]string, 0, len(points))
	args := make([]any, 0, len(points)*11)

	for i, p := range points {
		base := i * 12
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			base+1, base+2, base+3, base+4, base+5, base+6,
			base+7, base+8, base+9, base+10, base+11, base+12,
		))
		args = append(args,
			p.RouteID, p.OfficeFromID, p.Timestamp,
			p.Status1, p.Status2, p.Status3, p.Status4,
			p.Status5, p.Status6, p.Status7, p.Status8,
			p.Target2H,
		)
	}

	query := `INSERT INTO raw_data
		(route_id, office_from_id, timestamp,
		 status_1, status_2, status_3, status_4,
		 status_5, status_6, status_7, status_8,
		 target_2h)
		VALUES ` + strings.Join(placeholders, ",")

	res, err := db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		db.log.Error(op, "error", err)
		return 0, coreErrors.ErrExecQuery
	}

	n, _ := res.RowsAffected()
	return int(n), nil
}

func (db *DB) GetRawData(ctx context.Context, from, to *time.Time) ([]models.RawDataPoint, error) {
	const op = "db.GetRawData"

	query := `SELECT id, route_id, office_from_id, timestamp,
		status_1, status_2, status_3, status_4,
		status_5, status_6, status_7, status_8,
		target_2h, created_at FROM raw_data WHERE 1=1`
	args := []any{}

	if from != nil {
		args = append(args, *from)
		query += fmt.Sprintf(" AND timestamp>=$%d", len(args))
	}
	if to != nil {
		args = append(args, *to)
		query += fmt.Sprintf(" AND timestamp<=$%d", len(args))
	}
	query += " ORDER BY timestamp"

	var list []models.RawDataPoint
	if err := db.conn.SelectContext(ctx, &list, query, args...); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}