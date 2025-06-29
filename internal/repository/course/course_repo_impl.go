package course

import (
	"context"

	"gorm.io/gorm"

	courseEntity "github.com/itmrchow/course-management-system/internal/domain/course/entity"
	repo "github.com/itmrchow/course-management-system/internal/repository"
)

type CourseRepositoryImpl struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &CourseRepositoryImpl{db: db}
}

func (r *CourseRepositoryImpl) Create(ctx context.Context, course *courseEntity.Course) (uint, error) {
	result := r.db.WithContext(ctx).Create(course)
	if result.Error != nil {
		return 0, result.Error
	}

	return course.ID, result.Error
}

func (r *CourseRepositoryImpl) GetByID(ctx context.Context, id uint) (*courseEntity.Course, error) {
	return nil, nil
}

func (r *CourseRepositoryImpl) Update(ctx context.Context, course *courseEntity.Course) (int64, error) {
	return 0, nil
}

func (r *CourseRepositoryImpl) Delete(ctx context.Context, id uint) (int64, error) {
	return 0, nil
}

func (r *CourseRepositoryImpl) Find(ctx context.Context, pageInfo *repo.RepoPageInfo, conditions []func(db *gorm.DB) *gorm.DB) ([]*courseEntity.Course, error) {
	return nil, nil
}
