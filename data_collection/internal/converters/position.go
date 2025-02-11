package converters

import (
	model "mole/data_collection/v1/models"
	db "mole/db/sqlc"
)

// Convert a database position to an API position (sql.NullFloat64 to float64)
func DBtoAPI(dbPosition *db.InsertPositionParams) *model.CreatePosition {
	return &model.CreatePosition{
		InstanceID:        dbPosition.InstanceID,
		LatitudeDegrees:   dbPosition.LatitudeDegrees,
		LongitudeDegrees:  dbPosition.LongitudeDegrees,
		HeadingDegrees:    NullFloat64ToFloat64(dbPosition.HeadingDegrees),
		AltitudeHaeMeters: NullFloat64ToFloat64(dbPosition.AltitudeHaeMeters),
		SpeedMps:          NullFloat64ToFloat64(dbPosition.SpeedMps),
	}
}
