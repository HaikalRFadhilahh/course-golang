package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	Id        int       `gorm:"type:int;primaryKey;autoIncrement" json:"id,omitempty"`
	Role      string    `gorm:"type:varchar(10)" json:"role"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Email     string    `gorm:"type:varchar(50)" json:"email" binding:"required,email"`
	Password  string    `gorm:"type:varchar(255)" json:"password" binding:"required,min=5"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Tasks     []Task    `gorm:"constraint:OnDelete:CASCADE" json:"tasks,omitempty"` // Has Many
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.Id).Delete(&Task{})
	return
}
