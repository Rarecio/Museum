package models

import (
	"time"
)

type UserRole string

const (
	RoleStaff UserRole = "staff"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ExhibitStatus string

const (
	PendingApprovalStatus ExhibitStatus = "pending_approval"
	ApproveStatus         ExhibitStatus = "approve"
)

// экспонат
type Exhibit struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
	//Material        *string       `json:"material"`
	SaveBy          *int64        `json:"save_by"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	StorageLocation string        `json:"storage_location"`
	Status          ExhibitStatus `json:"status"`
	Length          int32         `json:"length"`
	Width           int32         `json:"width"`
	Height          int32         `json:"height"`
	Weight          int32         `json:"weight"`
	CreationDate    time.Time     `json:"creation_date"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type Category struct {
	Name string `json:"name"`
}

type Exhibition struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
