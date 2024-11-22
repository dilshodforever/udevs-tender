package model

import (
	"time"
)

// User represents the users table.
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"size:255;not null" json:"full_name"`
	Password string `gorm:"size:255;not null" json:"password"`
	Role     string `gorm:"size:50;not null;check:role IN ('client', 'contractor')" json:"role"` // Restrict role to "client" or "contractor"
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
}

// Tender represents the tenders table.
type Tender struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID            int64     `gorm:"not null" json:"client_id"`
	Title               string    `gorm:"size:255;not null" json:"title"`
	Description         string    `gorm:"type:text;not null" json:"description"`
	Deadline            time.Time `gorm:"not null" json:"deadline"`
	Budget              float64   `gorm:"not null" json:"budget"`
	Status              string    `gorm:"size:50;not null;check:status IN ('open', 'closed', 'pending', 'awarded')" json:"status"` // Restrict status to predefined values
	AwardedContractorID int64     `json:"awarded_contractor_id"`
}

// Bid represents the bids table.
type Bid struct {
	ID           int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	TenderID     int64   `gorm:"not null" json:"tender_id"`
	ContractorID int64   `gorm:"not null" json:"contractor_id"`
	Price        float64 `gorm:"not null" json:"price"`
	DeliveryTime int     `gorm:"not null" json:"delivery_time"`
	Comments     string  `gorm:"type:text" json:"comments"`
	Status       string  `gorm:"size:50;not null;check:status IN ('accepted', 'rejected', 'pending')" json:"status"` // Restrict status to predefined values
}

// Notification represents the notifications table.
type Notification struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64      `gorm:"not null" json:"user_id"`
	Message     string     `gorm:"type:text;not null" json:"message"`
	IsDelivered bool       `gorm:"not null" json:"is_delivered"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
}
