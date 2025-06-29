package entity

import (
	"gorm.io/gorm"
)

// Teacher 代表教師實體
// 包含教師的基本資料與關聯帳號資訊
type Teacher struct {
	gorm.Model
	UserID uint          `gorm:"not null;uniqueIndex"`                   // 關聯 user（帳號）ID
	Name   string        `gorm:"type:varchar(100);not null"`             // 教師姓名 , 原則上跟User的Name一樣
	Phone  string        `gorm:"type:varchar(20);not null;uniqueIndex"`  // 聯絡電話 , 原則上跟User的Phone一樣
	Email  string        `gorm:"type:varchar(100);not null;uniqueIndex"` // 教師信箱 , 原則上跟User的Email一樣
	Bio    string        `gorm:"type:text"`                              // 教師簡介/自我介紹
	Status TeacherStatus `gorm:"not null"`                               // 狀態 , 0: 審核中 , 1: 審核通過 , 2: 審核失敗 , 3: 已停用
}

type TeacherStatus uint

const (
	TeacherStatusPending  TeacherStatus = iota // 審核中
	TeacherStatusApproved                      // 審核通過
	TeacherStatusRejected                      // 審核失敗
	TeacherStatusDisabled                      // 已停用
)

// 查詢LikeName 查詢LikeName
func LikeTeacherName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name LIKE ?", "%"+name+"%")
	}
}

// 查詢狀態 查詢狀態
func InTeacherStatus(statuses []TeacherStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", statuses)
	}
}
