package course

import (
	"time"

	"gorm.io/gorm"
)

// Course 課程
type Course struct {
	gorm.Model
	Name                  string       `gorm:"type:varchar(100);not null"` // 課程名稱
	Description           string       `gorm:"type:text;not null"`         // 課程描述
	Price                 uint         `gorm:"not null;default:0"`         // 課程價格
	MaxStudents           uint         `gorm:"not null;default:0"`         // 最大學生人數
	MinStudents           uint         `gorm:"not null;default:0"`         // 最小學生人數
	RegistrationStartDate time.Time    `gorm:"not null"`                   // 報名開始時間
	RegistrationEndDate   time.Time    `gorm:"not null"`                   // 報名結束時間
	StartDate             time.Time    `gorm:"not null"`                   // 上課開始時間
	EndDate               time.Time    `gorm:"not null"`                   // 上課結束時間
	IsOnline              bool         `gorm:"not null;default:false"`     // 是否是線上課程
	Status                CourseStatus `gorm:"not null;default:0"`         // 0: 草稿, 1: 審核中, 2: 開放報名, 3: 已結束 , 4: 暫停報名
	Note                  string       `gorm:"type:text"`                  // 課程備註
}

type CourseStatus uint

const (
	CourseStatusDraft   CourseStatus = iota // 草稿
	CourseStatusPending                     // 審核中
	CourseStatusOnline                      // 開放報名
	CourseStatusEnd                         // 已結束
	CourseStatusPause                       // 暫停報名
)
