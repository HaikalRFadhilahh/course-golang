package models

import "time"

type Task struct {
	Id           *int      `gorm:"type:int; primaryKey; autoIncrement" json:"id,omitempty"`
	UserId       *int      `gorm:"int" json:"userId,omitempty"`
	Title        *string   `gorm:"type:varchar(255)" json:"title,omitempty"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	Status       *string   `gorm:"type:varchar(50)" json:"status,omitempty"`
	Reason       *string   `gorm:"type:text; default:" json:"reason,omitempty"`
	Revision     *int8     `gorm:"type:int; default:0" json:"revision,omitempty"`
	DueDate      *string   `gorm:"type:varchar(50)" json:"dueDate,omitempty"`
	SubmitDate   *string   `gorm:"type:varchar(50)" json:"submitDate,omitempty"`
	RejectedDate *string   `gorm:"type:varchar(50)" json:"rejectedDate,omitempty"`
	ApprovedDate *string   `gorm:"type:varchar(50)" json:"approvedDate,omitempty"`
	Attachment   *string   `gorm:"type:varchar(255)" json:"attachment,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	User         *User     `gorm:"foreignKey:UserId" json:"user,omitempty" binding:"omitempty"` // belongs to
}
