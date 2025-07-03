package handler

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"pitanguinha.com/audio-converter/internal/database"
	"pitanguinha.com/audio-converter/internal/utils"
)

// UpdateDocumentInput holds the input parameters for updating a document in the database.
type UpdateDocumentInput struct {
	ID             string
	CollectionName string
	ContentKey     string
	Duration       float64
	Status         Status
}

// Status represents the status of a document update operation.
type Status uint8

const (
	Success Status = iota
	Failure
)

// SetStatus converts a boolean value to a Status type.
func SetStatus(isSuccess bool) Status {
	if isSuccess {
		return Success
	}
	return Failure
}

// UpdateDocument updates the status of a document in the specified collection.
func (doc *UpdateDocumentInput) UpdateDocument() error {
	db, err := database.GetDatabase()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	collection := db.Collection(doc.CollectionName)

	var setFields any
	switch doc.Status {
	case Success:
		setFields = map[string]any{
			"conversion_status": "SUCCESS",
			"content_key":       doc.ContentKey,
			"duration":          utils.FormatSecondsToTime(doc.Duration),
		}
	case Failure:
		setFields = map[string]any{
			"conversion_status": "ERROR",
		}
	}

	updateBson := bson.M{
		"$set": setFields,
	}

	id, err := bson.ObjectIDFromHex(strings.TrimSpace(doc.ID))
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := collection.UpdateByID(context.TODO(), id, updateBson)
	if err != nil {
		return fmt.Errorf("failed to update document with ID %s: %w", id, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID %s", id)
	}

	return nil
}
