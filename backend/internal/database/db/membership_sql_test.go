package db

import (
    "context"
    "testing"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgtype"
    "github.com/jackc/pgx/v5/pgconn"
)

// fakeDBTX captures the last Exec query and arguments.
type fakeDBTX struct {
    lastQuery string
}

func (f *fakeDBTX) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
    f.lastQuery = sql
    var zero pgconn.CommandTag
    return zero, nil
}

func (f *fakeDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
    return nil, nil
}

func (f *fakeDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
    return nil
}

func TestCancelActiveMembershipsUsesActiveTimeWindow(t *testing.T) {
    f := &fakeDBTX{}
    q := New(f)

    // Build params similar to generated CancelActiveMembershipsByUserIdParams
    var userID pgtype.UUID
    _ = userID.Scan("00000000-0000-0000-0000-000000000000")
    q.CancelActiveMembershipsByUserId(context.Background(), CancelActiveMembershipsByUserIdParams{
        UserID: userID,
        CancelledAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
    })

    sql := f.lastQuery
    if sql == "" {
        t.Fatalf("expected Exec to be called, but no query captured")
    }

    // Ensure the WHERE clause includes the active time window predicate.
    if !contains(sql, "started_at <= NOW()") {
        t.Fatalf("expected query to contain started_at <= NOW(), got: %s", sql)
    }
    if !contains(sql, "expires_at > NOW()") {
        t.Fatalf("expected query to contain expires_at > NOW(), got: %s", sql)
    }
    if !contains(sql, "cancelled_at IS NULL") {
        t.Fatalf("expected query to contain cancelled_at IS NULL, got: %s", sql)
    }
}

func contains(s, sub string) bool {
    return len(s) >= len(sub) && (indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
    for i := 0; i+len(sub) <= len(s); i++ {
        if s[i:i+len(sub)] == sub {
            return i
        }
    }
    return -1
}
