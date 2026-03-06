package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UseCaseID       string              `json:"use_case_id" bson:"use_case_id"`
	FlowID          string              `json:"flow_id" bson:"flow_id"`
	ActionID        string              `json:"action_id" bson:"action_id"`
	JSONPath        string              `json:"json_path" bson:"json_path"`
	Comment         string              `json:"comment" bson:"comment"`
	ParentCommentID *primitive.ObjectID `json:"parent_comment_id,omitempty" bson:"parent_comment_id,omitempty"`
	Resolved        bool                `json:"resolved" bson:"resolved"`
	CreatedBy       string              `json:"created_by" bson:"created_by"`
	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
}
