package handlers

import (
	"bytes"
	"encoding/json"
	"time"

	"automation-developer-guide/src/database"
	"automation-developer-guide/src/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validateNote = validator.New()

// HandleCreateNote adds a new note with strict validation
func HandleCreateNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	type CreatePayload struct {
		UseCaseID string `json:"use_case_id" validate:"required"`
		FlowID    string `json:"flow_id" validate:"required"`
		ActionID  string `json:"action_id" validate:"required"`
		JSONPath  string `json:"json_path" validate:"required"`
		Note      string `json:"note" validate:"required"`
	}

	payload := new(CreatePayload)
	// Use json.NewDecoder to access DisallowUnknownFields
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields() // This throws an error if extra fields are present

	if err := decoder.Decode(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json or unknown fields"})
	}

	// Validate required fields
	if err := validateNote.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	note := models.Note{
		UseCaseID: payload.UseCaseID,
		FlowID:    payload.FlowID,
		ActionID:  payload.ActionID,
		JSONPath:  payload.JSONPath,
		Note:      payload.Note,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := database.CreateOne("notes", note)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save note"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "note created", "id": result.InsertedID})
}

// HandleUpdateNote updates a note with strict validation
func HandleUpdateNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	noteID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid note id"})
	}

	type UpdatePayload struct {
		Note string `json:"note" validate:"required"`
	}

	payload := new(UpdatePayload)
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json or unknown fields"})
	}

	if err := validateNote.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.M{"_id": objID, "created_by": userID}
	update := bson.M{
		"$set": bson.M{
			"note":       payload.Note,
			"updated_at": time.Now(),
		},
	}

	result, err := database.UpdateOne("notes", filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update note"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "note not found or unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "note updated"})
}

// HandleGetNotes retrieves notes (no body validation needed)
func HandleGetNotes(c *fiber.Ctx) error {
	filter := bson.M{}

	if val := c.Query("use_case_id"); val != "" {
		filter["use_case_id"] = val
	}
	if val := c.Query("flow_id"); val != "" {
		filter["flow_id"] = val
	}
	if val := c.Query("action_id"); val != "" {
		filter["action_id"] = val
	}
	if val := c.Query("json_path"); val != "" {
		filter["json_path"] = val
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.M{"created_at": -1}}},
		{{Key: "$addFields", Value: bson.M{"created_by_oid": bson.M{"$toObjectId": "$created_by"}}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "created_by_oid",
			"foreignField": "_id",
			"as":           "user_details",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$user_details", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$addFields", Value: bson.M{
			"user": bson.M{"email": "$user_details.email", "username": "$user_details.login"},
		}}},
		{{Key: "$project", Value: bson.M{"user_details": 0, "created_by_oid": 0}}},
	}

	var notes []bson.M = []bson.M{}
	if err := database.Aggregate("notes", pipeline, &notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch notes"})
	}

	return c.JSON(notes)
}

// HandleDeleteNote deletes a note (no body validation needed)
func HandleDeleteNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	noteID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid note id"})
	}

	filter := bson.M{"_id": objID, "created_by": userID}

	result, err := database.DeleteOne("notes", filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete note"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "note not found or unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "note deleted"})
}
