package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../")

	rootPath := fmt.Sprintf("%s%s", basePath, string(os.PathSeparator))
	return rootPath
}
