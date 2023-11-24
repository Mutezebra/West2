package model

type Group struct {
	ID        uint   `gorm:"primarykey"`
	GroupName string `gorm:"type:varchar(20);not null,unique"`
}
