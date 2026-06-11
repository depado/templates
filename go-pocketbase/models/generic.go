package models

import (
	"github.com/pocketbase/pocketbase/core"
)

// FindById finds a single record by its ID.
func FindById(app core.App, collection string, id string) (*core.Record, error) {
	return app.FindRecordById(collection, id)
}

// FindAllByField finds all records matching a field value.
func FindAllByField(app core.App, collection string, field string, value any) ([]*core.Record, error) {
	return app.FindRecordsByFilter(
		collection,
		field+" = {:value}",
		"-created",
		0, 0,
		map[string]any{"value": value},
	)
}

// FindFirst finds the first record matching a filter.
func FindFirst(app core.App, collection string, filter string, params map[string]any) (*core.Record, error) {
	return app.FindFirstRecordByFilter(collection, filter, params)
}
