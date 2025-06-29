package teacher

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/itmrchow/course-management-system/internal/domain/teacher/entity"
	repo "github.com/itmrchow/course-management-system/internal/repository"
)

var _ TeacherRepository = (*TeacherRepositoryImpl)(nil)

// TeacherRepositoryImpl 實作 TeacherRepository 介面
// 負責教師資料的存取操作
//
// Example:
//
//	repo := teacher.NewTeacherRepository(db)
//	teacher, err := repo.GetByID(ctx, 1)
type TeacherRepositoryImpl struct {
	db *gorm.DB
}

// NewTeacherRepository 建立教師資料庫操作實例
// 參數: db - 資料庫連線
// 回傳: 教師資料庫操作實例
func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &TeacherRepositoryImpl{db: db}
}

func (t *TeacherRepositoryImpl) WithTransaction(tx *gorm.DB) TeacherRepository {
	return &TeacherRepositoryImpl{db: tx}
}

// Create 建立教師資料
// 使用 gorm.Create 建立教師資料
// 如果建立失敗，返回錯誤
// 如果建立成功，返回教師ID
func (t *TeacherRepositoryImpl) Create(ctx context.Context, teacher *entity.Teacher) (uint, error) {
	result := t.db.WithContext(ctx).Create(teacher)
	if result.Error != nil {
		return 0, result.Error
	}

	return teacher.ID, result.Error
}

// GetByID 依教師ID查詢教師資料
// 使用 gorm.First 查詢教師資料
// 如果查詢失敗，返回錯誤
// 如果查詢成功，返回教師資料
func (t *TeacherRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Teacher, error) {
	var teacher entity.Teacher
	if err := t.db.WithContext(ctx).First(&teacher, id).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

// Update 更新教師資料
// 使用 gorm.Model.Updates 更新教師資料
// 如果更新失敗，返回錯誤
// 如果更新成功，返回更新後的教師資料數量
func (t *TeacherRepositoryImpl) Update(ctx context.Context, teacher *entity.Teacher) (int64, error) {
	result := t.db.WithContext(ctx).Model(teacher).Updates(teacher)
	return result.RowsAffected, result.Error
}

// Delete 刪除教師資料
// 使用 gorm.Delete 刪除教師資料
// 如果刪除失敗，返回錯誤
// 如果刪除成功，返回刪除的教師資料數量
func (t *TeacherRepositoryImpl) Delete(ctx context.Context, id uint) (int64, error) {
	result := t.db.WithContext(ctx).Delete(&entity.Teacher{}, id)
	return result.RowsAffected, result.Error
}

// Find 查詢教師資料
// 使用 gorm.Find 查詢教師資料
// 如果查詢失敗，返回錯誤
// 如果查詢成功，返回教師資料
func (t *TeacherRepositoryImpl) Find(ctx context.Context, pageInfo *repo.RepoPageInfo, conditions []func(db *gorm.DB) *gorm.DB) ([]*entity.Teacher, error) {
	var teachers []*entity.Teacher

	dbFuncs := []func(db *gorm.DB) *gorm.DB{repo.Paginate(pageInfo)}
	dbFuncs = append(dbFuncs, conditions...)

	if err := t.db.
		WithContext(ctx).
		Scopes(dbFuncs...).
		Order(fmt.Sprintf("%s %s", pageInfo.Sort, pageInfo.Order)).
		Find(&teachers).
		Error; err != nil {
		return nil, err
	}
	return teachers, nil
}
