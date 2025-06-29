package teacher

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/itmrchow/course-management-system/internal/config"
	"github.com/itmrchow/course-management-system/internal/domain/teacher/entity"
	"github.com/itmrchow/course-management-system/internal/repository"
)

// TeacherRepoTestSuite 用於 TeacherRepositoryImpl 的測試
// 會在這裡做初始化與測試案例
type TeacherRepoTestSuite struct {
	suite.Suite
	teacherRepo TeacherRepository
	db          *gorm.DB
}

// SetupTest 於每個測試案例前執行，做初始化
func (s *TeacherRepoTestSuite) SetupTest() {
	// logger
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	// db
	db := config.NewMemoryDB(context.Background(), &logger)

	sqlDB, err := db.DB()
	s.Require().NoError(err)

	// fixture
	// 載入 fixture
	dir, _ := os.Getwd()
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB), // 傳入 *sql.DB
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(filepath.Join(dir, "../../repository/teacher/testdata")),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	s.db = db
	s.teacherRepo = NewTeacherRepository(db)

}

// TestTeacherRepoSuite 執行測試套件
func TestTeacherRepoSuite(t *testing.T) {
	suite.Run(t, new(TeacherRepoTestSuite))
}

// 依教師ID查詢教師資料測試
// Test for TeacherRepositoryImpl.GetByID
func (s *TeacherRepoTestSuite) TestGetByID() {
	// input
	teacherID := uint(1)

	// expect
	teacher, err := s.teacherRepo.GetByID(context.Background(), teacherID)
	s.NoError(err, `"teacherRepo.GetByID" should not return error`)
	s.Equal(teacher.Name, "John Doe")
	s.Equal(teacher.Email, "john.doe@example.com")
	s.Equal(teacher.Phone, "1234567890")
	s.Equal(teacher.Bio, "I am a teacher")
	s.EqualValues(teacher.Status, 0)
	s.EqualValues(teacher.ID, teacherID)
	s.EqualValues(teacher.UserID, 1)
}

// 建立教師資料測試
// Test for TeacherRepositoryImpl.Create
func (s *TeacherRepoTestSuite) TestCreate() {
	// input
	type args struct {
		ctx     context.Context
		teacher *entity.Teacher
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, teacherID uint, err error)
	}{
		{
			name: "exists phone",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					UserID: 2,
					Name:   "Test",
					Phone:  "1234567890",
					Email:  "test@example.com",
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, teacherID uint, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.Equal(t, teacherID, uint(0))
			},
		},
		{
			name: "exists email",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					UserID: 3,
					Name:   "Test",
					Phone:  "0987654321",
					Email:  "john.doe@example.com", // 已存在的 email
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, teacherID uint, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.Equal(t, teacherID, uint(0))
			},
		},
		{
			name: "exists user_id",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					UserID: 1, // 已存在的 user_id
					Name:   "Test",
					Phone:  "0987654321",
					Email:  "test2@example.com",
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, teacherID uint, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.Equal(t, teacherID, uint(0))
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					UserID: 100,
					Name:   "Test Success",
					Phone:  "0911222333",
					Email:  "success@example.com",
					Bio:    "I am a teacher",
					Status: 1,
				},
			},
			assertFunc: func(t *testing.T, teacherID uint, err error) {
				assert.NoError(t, err)
				assert.NotEqual(t, teacherID, uint(0))
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			teacherID, err := s.teacherRepo.Create(test.args.ctx, test.args.teacher)
			test.assertFunc(s.T(), teacherID, err)
		})
	}
}

// 更新教師資料測試
// Test for TeacherRepositoryImpl.Update
func (s *TeacherRepoTestSuite) TestUpdate() {
	type args struct {
		ctx     context.Context
		teacher *entity.Teacher
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, rowsAffected int64, err error)
	}{
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					Model: gorm.Model{
						ID: 100,
					},
					UserID: 1,
					Email:  "not_found@example.com",
					Name:   "Not Found",
					Phone:  "1234567890",
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, rowsAffected, 0)
			},
		},
		{
			name: "duplicate phone",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					Model: gorm.Model{
						ID: 1,
					},
					UserID: 2,
					Email:  "duplicate_phone@example.com",
					Name:   "Duplicate Phone",
					Phone:  "1234567890",
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.EqualValues(t, rowsAffected, 0)
			},
		},
		{
			name: "duplicate email",
			args: args{
				ctx: context.Background(),
				teacher: func() *entity.Teacher {
					teacher := &entity.Teacher{}
					teacher.ID = 2
					teacher.UserID = 2
					teacher.Email = "john.doe@example.com"
					teacher.Name = "John Doe 2"
					teacher.Phone = "1234500000"
					teacher.Bio = "I am a teacher 2"
					teacher.Status = 1
					return teacher
				}(),
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.EqualValues(t, rowsAffected, 0)
			},
		},
		{
			name: "duplicate user_id",
			args: args{
				ctx: context.Background(),
				teacher: &entity.Teacher{
					Model: gorm.Model{
						ID: 2,
					},
					UserID: 1,
					Email:  "duplicate_user_id@example.com",
					Name:   "Duplicate User ID",
					Phone:  "1234500000",
					Bio:    "I am a teacher",
					Status: 0,
				},
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.ErrorAs(t, err, &gorm.ErrDuplicatedKey)
				assert.EqualValues(t, rowsAffected, 0)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				teacher: func() *entity.Teacher {
					teacher := &entity.Teacher{}
					teacher.ID = 2
					teacher.UserID = 2
					teacher.Email = "john.doe_update@example.com"
					teacher.Name = "Updated Name"
					teacher.Phone = "1234500000"
					teacher.Bio = "Updated bio"
					teacher.Status = 3
					return teacher
				}(),
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, rowsAffected, 1)
				// 取得更新後的資料做驗證
				updated, getErr := s.teacherRepo.GetByID(context.Background(), 2)
				assert.NoError(t, getErr)
				assert.Equal(t, "Updated Name", updated.Name)
				assert.Equal(t, "john.doe_update@example.com", updated.Email)
				assert.Equal(t, "1234500000", updated.Phone)
				assert.Equal(t, "Updated bio", updated.Bio)
				assert.EqualValues(t, entity.TeacherStatus(3), updated.Status)
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			rowsAffected, err := s.teacherRepo.Update(test.args.ctx, test.args.teacher)
			test.assertFunc(s.T(), rowsAffected, err)
		})
	}
}

// 刪除教師資料測試
// Test for TeacherRepositoryImpl.Delete
func (s *TeacherRepoTestSuite) TestDelete() {
	type args struct {
		ctx context.Context
		id  uint
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, rowsAffected int64, err error)
	}{
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				id:  100,
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, rowsAffected, 0)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			assertFunc: func(t *testing.T, rowsAffected int64, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, rowsAffected, 1)
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			rowsAffected, err := s.teacherRepo.Delete(test.args.ctx, test.args.id)
			test.assertFunc(s.T(), rowsAffected, err)
		})
	}
}

// 取得教師清單測試
// Test for TeacherRepositoryImpl.List
func (s *TeacherRepoTestSuite) TestList() {
	type args struct {
		ctx        context.Context
		pageInfo   *repository.RepoPageInfo
		conditions []func(db *gorm.DB) *gorm.DB
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, teachers []*entity.Teacher, err error)
	}{
		{
			name: "success_page_size",
			args: args{
				ctx: context.Background(),
				pageInfo: &repository.RepoPageInfo{
					Page:     1,
					PageSize: 1,
					Sort:     "id",
					Order:    "desc",
				},
				conditions: []func(db *gorm.DB) *gorm.DB{},
			},
			assertFunc: func(t *testing.T, teachers []*entity.Teacher, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(teachers), 1)
				assert.Equal(t, teachers[0].ID, uint(10))
				assert.Equal(t, teachers[0].Name, "Thomas Anderson")
				assert.Equal(t, teachers[0].Email, "thomas.anderson@example.com")
				assert.Equal(t, teachers[0].Phone, "9012345678")
				assert.Equal(t, teachers[0].Bio, "Music instructor teaching classical piano and composition")
				assert.EqualValues(t, entity.TeacherStatus(2), teachers[0].Status)
				assert.EqualValues(t, teachers[0].UserID, 10)
			},
		},
		{
			name: "success_like_name",
			args: args{
				ctx: context.Background(),
				pageInfo: &repository.RepoPageInfo{
					Page:     1,
					PageSize: 10,
					Sort:     "id",
					Order:    "desc",
				},
				conditions: []func(db *gorm.DB) *gorm.DB{
					entity.LikeTeacherName("John"),
				},
			},
			assertFunc: func(t *testing.T, teachers []*entity.Teacher, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(teachers), 3)

				for _, teacher := range teachers {
					assert.Contains(t, teacher.Name, "John")
				}
			},
		},
		{
			name: "success_in_status",
			args: args{
				ctx: context.Background(),
				pageInfo: &repository.RepoPageInfo{
					Page:     1,
					PageSize: 10,
					Sort:     "id",
					Order:    "desc",
				},
				conditions: []func(db *gorm.DB) *gorm.DB{
					entity.InTeacherStatus([]entity.TeacherStatus{entity.TeacherStatusPending, entity.TeacherStatusRejected}),
				},
			},
			assertFunc: func(t *testing.T, teachers []*entity.Teacher, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(teachers), 4)

				for _, teacher := range teachers {
					assert.Contains(t, []entity.TeacherStatus{entity.TeacherStatusPending, entity.TeacherStatusRejected}, teacher.Status)
				}
			},
		},
		{
			name: "success_by_email",
			args: args{
				ctx: context.Background(),
				pageInfo: &repository.RepoPageInfo{
					Page:     1,
					PageSize: 10,
					Sort:     "id",
					Order:    "desc",
				},
				conditions: []func(db *gorm.DB) *gorm.DB{
					func(db *gorm.DB) *gorm.DB {
						return db.Where("email = ?", "jennifer.lee@example.com")
					},
				},
			},
			assertFunc: func(t *testing.T, teachers []*entity.Teacher, err error) {
				assert.NoError(t, err)
				assert.Equal(t, len(teachers), 1)
				assert.Equal(t, teachers[0].Email, "jennifer.lee@example.com")
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			teachers, err := s.teacherRepo.Find(test.args.ctx, test.args.pageInfo, test.args.conditions)
			test.assertFunc(s.T(), teachers, err)
		})
	}

}
