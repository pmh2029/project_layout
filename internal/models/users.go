package models

// UsersTableName TableName
var UsersTableName = "users"

type User struct {
	ID       uint   `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Username string `gorm:"column:username;type:varchar(50);not null;unique" mapstructure:"username"`
	Email    string `gorm:"column:email;type:varchar(50);not null;unique" mapstructure:"email"`
	Password string `gorm:"column:password;not null" mapstructure:"password"`
	BaseModel
}

// TableName func
func (i *User) TableName() string {
	return UsersTableName
}
