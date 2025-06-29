package teacher

import (
	"os"
	"testing"

	"github.com/itmrchow/course-management-system/internal/repository"
)

// TestMain 初始化 repository 測試環境
// repo 測試呼叫RepoTestInit 初始化測試環境
func TestMain(m *testing.M) {
	code := repository.RepoTestInit(m)

	os.Exit(code)
}
