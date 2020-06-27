package model

import "strings"

// MoodPayload for creating a new Mood
type MoodPayload struct {
	Name  string `json:"name,omitempty" validate:"required,lte=64"`
	Color string `json:"color,omitempty" validate:"required,hexcolor"`
}

// Validate struct fields
func (p *MoodPayload) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	return validate.Struct(p)
}

// MoodChanges for updating a Mood
type MoodChanges struct {
	Name  string `json:"name" validate:"lte=64"`
	Color string `json:"color" validate:"hexcolor"`
}

// Validate struct fields
func (p *MoodChanges) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	return validate.Struct(p)
}
