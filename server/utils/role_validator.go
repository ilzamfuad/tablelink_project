package utils

import (
	"errors"
	"net/http"
	"tablelink_project/server/model"

	"gorm.io/gorm"
)

// RoleRightsValidator is a struct to validate role rights
type RoleRightsValidator struct {
	db *gorm.DB
}

// NewRoleRightsValidator creates a new instance of RoleRightsValidator
func NewRoleRightsValidator(db *gorm.DB) *RoleRightsValidator {
	return &RoleRightsValidator{db: db}
}

// ValidateRoleRights validates the role rights based on section, route, and method
func (v *RoleRightsValidator) ValidateRoleRights(section, route, method string) error {
	var roleRight model.RoleRight

	// Check if the section and route exist in the role_rights table
	err := v.db.Where("section = ? AND route = ?", section, route).First(&roleRight).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("access denied: invalid section or route")
		}
		return err
	}

	// Validate the HTTP method based on the role_rights table
	switch method {
	case http.MethodPost:
		if roleRight.RCreate != 1 {
			return errors.New("access denied: POST method not allowed")
		}
	case http.MethodGet:
		if roleRight.RRead != 1 {
			return errors.New("access denied: GET method not allowed")
		}
	case http.MethodPut:
		if roleRight.RUpdate != 1 {
			return errors.New("access denied: PUT method not allowed")
		}
	case http.MethodDelete:
		if roleRight.RDelete != 1 {
			return errors.New("access denied: DELETE method not allowed")
		}
	default:
		return errors.New("access denied: unsupported HTTP method")
	}

	return nil
}
