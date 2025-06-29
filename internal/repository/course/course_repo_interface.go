package course

import (
	"context"

	"gorm.io/gorm"

	courseEntity "github.com/itmrchow/course-management-system/internal/domain/course/entity"
	repo "github.com/itmrchow/course-management-system/internal/repository"
)

// CourseRepository 定義課程資料存取的介面
// 負責課程的新增、查詢、更新、刪除等操作
//
// Example:
//
//	var repo CourseRepository
//	course, err := repo.GetByID(ctx, 1)
type CourseRepository interface {
	// Create 新增課程資料
	// 參數: ctx - context, course - 課程實體
	// 回傳: 新增後的課程ID, 錯誤訊息
	Create(ctx context.Context, course *courseEntity.Course) (uint, error)

	// GetByID 依課程ID查詢課程資料
	// 參數: ctx - context, id - 課程ID
	// 回傳: 課程實體, 錯誤訊息
	GetByID(ctx context.Context, id uint) (*courseEntity.Course, error)

	// Update 更新課程資料
	// 參數: ctx - context, course - 課程實體
	// 回傳: 影響數量, 錯誤訊息
	Update(ctx context.Context, course *courseEntity.Course) (int64, error)

	// Delete 刪除課程資料
	// 參數: ctx - context, id - 課程ID
	// 回傳: 影響數量, 錯誤訊息
	Delete(ctx context.Context, id uint) (int64, error)

	// Find 取得課程清單（可加分頁、條件查詢）
	// 參數: ctx - context
	// 回傳: 課程實體切片, 錯誤訊息
	Find(ctx context.Context, pageInfo *repo.RepoPageInfo, conditions []func(db *gorm.DB) *gorm.DB) ([]*courseEntity.Course, error)
}
