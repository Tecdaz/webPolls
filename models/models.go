package models

type User struct {
	ID       uint   `gorm:"primaryKey;column:id"`
	Username string `gorm:"unique;not null;column:username"`
	Password string `gorm:"not null;column:password"`
	Email    string `gorm:"unique;not null;column:email"`
	Polls    []Poll `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Poll struct {
	ID      uint     `gorm:"primaryKey;column:id"`
	Title   string   `gorm:"not null;column:title"`
	UserID  uint     `gorm:"not null;column:user_id"` // Clave foránea
	Options []Option `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE"`
}

type Option struct {
	ID      uint   `gorm:"primaryKey;column:id"`
	Content string `gorm:"not null;column:content"`
	PollID  uint   `gorm:"not null;column:poll_id"` // Clave foránea
	Correct bool   `gorm:"default:false;column:correct"`
}
