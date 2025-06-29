package teacher

import (
	"context"

	"github.com/itmrchow/course-management-system/internal/domain/teacher/entity"
	repo "github.com/itmrchow/course-management-system/internal/repository"
	"gorm.io/gorm"
)

// TeacherRepository 定義教師資料存取的介面
// 負責教師的新增、查詢、更新、刪除等操作
//
// Example:
//
//	var repo TeacherRepository
//	teacher, err := repo.GetByID(ctx, 1)
type TeacherRepository interface {
	// Create 新增教師資料
	// 參數: ctx - context, teacher - 教師實體
	// 回傳: 新增後的教師ID, 錯誤訊息
	Create(ctx context.Context, teacher *entity.Teacher) (uint, error)

	// GetByID 依教師ID查詢教師資料
	// 參數: ctx - context, id - 教師ID
	// 回傳: 教師實體, 錯誤訊息
	GetByID(ctx context.Context, id uint) (*entity.Teacher, error)

	// Update 更新教師資料
	// 參數: ctx - context, teacher - 教師實體
	// 回傳: 影響數量, 錯誤訊息
	Update(ctx context.Context, teacher *entity.Teacher) (int64, error)

	// Delete 刪除教師資料
	// 參數: ctx - context, id - 教師ID
	// 回傳: 影響數量, 錯誤訊息
	Delete(ctx context.Context, id uint) (int64, error)

	// Find 取得教師清單（可加分頁、條件查詢）
	// 參數: ctx - context
	// 回傳: 教師實體切片, 錯誤訊息
	Find(ctx context.Context, pageInfo *repo.RepoPageInfo, conditions []func(db *gorm.DB) *gorm.DB) ([]*entity.Teacher, error)
}
