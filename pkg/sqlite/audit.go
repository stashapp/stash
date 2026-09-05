package sqlite

import (
	"context"
	"time"
)

type AuditEvent struct {
	ID          int64     `db:"id"`
	OccurredAt  time.Time `db:"occurred_at"`
	ActorUserID *int64    `db:"actor_user_id"`
	EventType   string    `db:"event_type"`
	TargetType  *string   `db:"target_type"`
	TargetID    *string   `db:"target_id"`
	Result      string    `db:"result"`
	DetailsJSON *string   `db:"details_json"`
}

type AuditStore struct{}

func (s *AuditStore) Record(ctx context.Context, actorUserID int64, eventType, targetType, targetID, result string) error {
	_, err := dbWrapper.Exec(ctx, `INSERT INTO user_audit_events
		(occurred_at, actor_user_id, event_type, target_type, target_id, result, details_json)
		VALUES (?, ?, ?, ?, ?, ?, NULL)`, time.Now().UTC(), actorUserID, eventType, targetType, targetID, result)
	return err
}

func (s *AuditStore) RecordAuthentication(ctx context.Context, actorUserID *int64, eventType, result string) error {
	_, err := dbWrapper.Exec(ctx, `INSERT INTO user_audit_events
		(occurred_at, actor_user_id, event_type, target_type, target_id, result, details_json)
		VALUES (?, ?, ?, NULL, NULL, ?, NULL)`, time.Now().UTC(), actorUserID, eventType, result)
	return err
}

func (s *AuditStore) List(ctx context.Context) ([]AuditEvent, error) {
	var ret []AuditEvent
	err := dbWrapper.Select(ctx, &ret, `SELECT id, occurred_at, actor_user_id, event_type, target_type, target_id, result, details_json
		FROM user_audit_events ORDER BY id`)
	return ret, err
}

func (s *AuditStore) ListPage(ctx context.Context, limit, offset int) ([]AuditEvent, error) {
	var ret []AuditEvent
	err := dbWrapper.Select(ctx, &ret, `SELECT id, occurred_at, actor_user_id, event_type, target_type, target_id, result, details_json
		FROM user_audit_events ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	return ret, err
}
