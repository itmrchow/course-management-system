package course

import (
	"time"

	"gorm.io/gorm"
)

// CoursePattern 課程模式表
type CoursePattern struct {
	gorm.Model
	CourseID  uint      `gorm:"not null"` // 課程ID
	DayOfWeek uint      `gorm:"not null"` // 星期幾
	StartTime time.Time `gorm:"not null"` // 上課開始時間
	EndTime   time.Time `gorm:"not null"` // 上課結束時間
}
