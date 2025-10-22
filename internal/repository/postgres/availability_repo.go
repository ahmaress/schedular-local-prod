package postgres

import (
	"context"
	"time"

	"scheduler-service/internal/models"
	"scheduler-service/internal/repository"
)

type AvailabilityRepo struct{}

func NewAvailabilityRepo() *AvailabilityRepo { return &AvailabilityRepo{} }

func (r *AvailabilityRepo) InsertAvailabilityRule(ctx context.Context, q repository.Querier, ar *models.AvailabilityRule) error {
	now := time.Now().UTC()
	query := `INSERT INTO availability_rules
		(id, user_id, day_of_week, start_time, end_time, slot_length_minutes, title, available, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return q.QueryRow(ctx, query,
		ar.UserID, ar.DayOfWeek, ar.StartTime, ar.EndTime, ar.SlotLengthMins,
		ar.Title, ar.Available, now, now,
	).Scan(&ar.ID)
}

func (r *AvailabilityRepo) ListAvailabilityRules(ctx context.Context, q repository.Querier, userID string) ([]models.AvailabilityRule, error) {
	query := `SELECT id,user_id,day_of_week,start_time,end_time,slot_length_minutes,title,available,created_at,updated_at
		      FROM availability_rules WHERE user_id=$1 ORDER BY id`
	rows, err := q.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.AvailabilityRule
	for rows.Next() {
		var r models.AvailabilityRule
		var start, end string
		if err := rows.Scan(&r.ID, &r.UserID, &r.DayOfWeek, &start, &end,
			&r.SlotLengthMins, &r.Title, &r.Available, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		r.StartTime = start
		r.EndTime = end
		out = append(out, r)
	}
	return out, nil
}

func (r *AvailabilityRepo) UpdateAvailabilityRule(ctx context.Context, q repository.Querier, userID, ruleID string, ar *models.AvailabilityRule) (string, error) {
	now := time.Now().UTC()
	query := `UPDATE availability_rules
		SET start_time=$1, end_time=$2, slot_length_minutes=$3,
		    title=$4, available=$5, updated_at=$6
		WHERE id=$7 AND user_id=$8
		RETURNING id`
	var updatedID string
	err := q.QueryRow(ctx, query,
		ar.StartTime, ar.EndTime, ar.SlotLengthMins,
		ar.Title, ar.Available, now, ruleID, userID,
	).Scan(&updatedID)
	return updatedID, err
}
