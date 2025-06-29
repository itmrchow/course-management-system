package course

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/itmrchow/course-management-system/internal/config"
	courseEntity "github.com/itmrchow/course-management-system/internal/domain/course/entity"
)

type CourseRepoTestSuite struct {
	suite.Suite
	courseRepo CourseRepository
	db         *gorm.DB
}

func (s *CourseRepoTestSuite) SetupTest() {
	// logger
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	// db
	db := config.NewMemoryDB(context.Background(), &logger)

	// sqlDB, err := db.DB()
	// s.Require().NoError(err)

	// fixture
	// dir, _ := os.Getwd()
	// fixtures, err := testfixtures.New(
	// 	testfixtures.Database(sqlDB), // 傳入 *sql.DB
	// 	testfixtures.Dialect("postgres"),
	// 	testfixtures.Directory(filepath.Join(dir, "../../repository/teacher/testdata")),
	// 	testfixtures.DangerousSkipTestDatabaseCheck(),
	// )
	// s.Require().NoError(err)
	// s.Require().NoError(fixtures.Load())

	s.db = db
	s.courseRepo = NewCourseRepository(db)
}

// TestCourseRepoSuite 執行測試套件
func TestCourseRepoSuite(t *testing.T) {
	suite.Run(t, new(CourseRepoTestSuite))
}

func (s *CourseRepoTestSuite) TestCreate() {
	// input
	type args struct {
		ctx    context.Context
		course *courseEntity.Course
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, courseID uint, err error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				course: &courseEntity.Course{
					Name:                  "test course name",
					Description:           "test course description",
					Price:                 4000,
					MaxStudents:           10,
					MinStudents:           5,
					RegistrationStartDate: time.Date(2025, 7, 7, 0, 0, 0, 0, time.UTC),     // 2025-07-07
					RegistrationEndDate:   time.Date(2025, 7, 20, 23, 59, 59, 0, time.UTC), // 2025-07-20
					StartDate:             time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC),    // 2025-07-28
					EndDate:               time.Date(2025, 8, 8, 23, 59, 59, 0, time.UTC),  // 2025-08-08
					IsOnline:              false,
					Status:                courseEntity.CourseStatusPending,
					Note:                  "",
				},
			},
			assertFunc: func(t *testing.T, courseID uint, err error) {
				assert.NoError(t, err)
				assert.NotZero(t, courseID)
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			courseID, err := s.courseRepo.Create(test.args.ctx, test.args.course)
			test.assertFunc(s.T(), courseID, err)
		})
	}
}
