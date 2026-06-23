package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracedDB provides database operations with OpenTelemetry tracing
type TracedDB struct {
	serviceName string
	dbSystem    string
	tracer      trace.Tracer
}

// NewTracedDB creates a new TracedDB instance
func NewTracedDB(serviceName string) *TracedDB {
	return &TracedDB{
		serviceName: serviceName,
		dbSystem:    "in-memory",
		tracer:      otel.Tracer(serviceName + "/db"),
	}
}

// Query traces a database read operation
func (db *TracedDB) Query(ctx context.Context, table string, operation string, filter string) context.Context {
	ctx, span := db.tracer.Start(ctx, fmt.Sprintf("DB %s %s", operation, table),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("db.system", db.dbSystem),
		attribute.String("db.operation", operation),
		attribute.String("db.sql.table", table),
		attribute.String("db.statement", fmt.Sprintf("%s FROM %s WHERE %s", operation, table, filter)),
		attribute.String("db.name", db.serviceName),
	)
	// Simulate DB latency
	time.Sleep(time.Duration(500+time.Now().UnixNano()%1500) * time.Microsecond)
	span.SetStatus(codes.Ok, "")
	span.End()
	return ctx
}

// Insert traces a database write operation
func (db *TracedDB) Insert(ctx context.Context, table string, id string) context.Context {
	ctx, span := db.tracer.Start(ctx, fmt.Sprintf("DB INSERT %s", table),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("db.system", db.dbSystem),
		attribute.String("db.operation", "INSERT"),
		attribute.String("db.sql.table", table),
		attribute.String("db.statement", fmt.Sprintf("INSERT INTO %s (id=%s)", table, id)),
		attribute.String("db.name", db.serviceName),
	)
	time.Sleep(time.Duration(800+time.Now().UnixNano()%2000) * time.Microsecond)
	span.SetStatus(codes.Ok, "")
	span.End()
	return ctx
}

// Update traces a database update operation
func (db *TracedDB) Update(ctx context.Context, table string, id string, fields string) context.Context {
	ctx, span := db.tracer.Start(ctx, fmt.Sprintf("DB UPDATE %s", table),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("db.system", db.dbSystem),
		attribute.String("db.operation", "UPDATE"),
		attribute.String("db.sql.table", table),
		attribute.String("db.statement", fmt.Sprintf("UPDATE %s SET %s WHERE id=%s", table, fields, id)),
		attribute.String("db.name", db.serviceName),
	)
	time.Sleep(time.Duration(600+time.Now().UnixNano()%1500) * time.Microsecond)
	span.SetStatus(codes.Ok, "")
	span.End()
	return ctx
}
