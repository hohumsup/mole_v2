package converters

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sqlc-dev/pqtype"
)

func NullFloat64ToFloat64(n sql.NullFloat64) *float64 {
	if n.Valid {
		return &n.Float64
	}
	return nil
}

func Float64ToNullFloat64(f *float64) sql.NullFloat64 {
	if f != nil {
		return sql.NullFloat64{Float64: *f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

func Float64ToPGFloat8(value *float64) pgtype.Float8 {
	if value != nil {
		return pgtype.Float8{Float64: *value, Valid: true}
	}
	return pgtype.Float8{Valid: false}
}

func ConvertJSONToPQType(raw json.RawMessage) (pqtype.NullRawMessage, error) {
	if len(raw) == 0 {
		return pqtype.NullRawMessage{Valid: false}, nil
	}

	var js interface{}
	if err := json.Unmarshal(raw, &js); err != nil {
		return pqtype.NullRawMessage{}, fmt.Errorf("invalid JSON: %w", err)
	}

	return pqtype.NullRawMessage{
		RawMessage: raw,
		Valid:      true,
	}, nil
}
