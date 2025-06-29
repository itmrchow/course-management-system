package course

import "gorm.io/gorm"

type CourseTeacher struct {
	gorm.Model
	CourseID  uint `gorm:"not null;uniqueIndex:idx_course_teacher"` // 課程ID
	TeacherID uint `gorm:"not null;uniqueIndex:idx_course_teacher"` // 教師ID
	IsMain    bool `gorm:"not null;default:false"`                  // 是否為主教師
}
