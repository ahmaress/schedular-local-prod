package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"scheduler-service/internal/app"
)

// Querier abstracts pgx pool/tx for easier testing and transactions.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AvailabilityRepository interface {
	InsertAvailabilityRule(ctx context.Context, q Querier, r *app.AvailabilityRule) error
	ListAvailabilityRules(ctx context.Context, q Querier, userID string) ([]app.AvailabilityRule, error)
	UpdateAvailabilityRule(ctx context.Context, q Querier, userID, ruleID string, r *app.AvailabilityRule) (string, error)
}

type BookingRepository interface {
	ListBookingsInRange(ctx context.Context, q Querier, userID string, from, to AppTime) ([]app.Booking, error)
	ListBookings(ctx context.Context, q Querier, userID string, from, to AppTime, filtered bool) ([]app.Booking, error)
	CheckExistingBookingAtStart(ctx context.Context, q Querier, userID string, start AppTime) (string, error)
	InsertBooking(ctx context.Context, q Querier, b *app.Booking) (string, error)
	GetBookingStatus(ctx context.Context, q Querier, id string) (string, error)
	CancelBooking(ctx context.Context, q Querier, id string) (int64, error)
}

// AppTime is a lightweight alias to avoid importing time here; implemented in impl files.
type AppTime interface{}
