package models

type PackageSchema struct {
	ID             string `bson:"_id,omitempty"`
	Name           string `bson:"name,omitempty"`
	Description    string `bson:"description,omitempty"`
}
