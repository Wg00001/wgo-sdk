package wg

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func ReadYAMLToMap(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRelativePath
// @param: 来自内容根的路径
// @return: 目标文件相对于当前启动目录的位置
func GetRelativePath(pathDir string) string {
	dir, _ := os.Getwd()
	if len(pathDir) == 0 {
		return dir
	}
	root := strings.Split(pathDir, "/")[0]
	if strings.Contains(dir, root) {
		idx := strings.Index(dir, root)
		_ = idx
		s := dir[strings.Index(dir, root):]
		count := strings.Count(s, "/")
		for i := 0; i < count; i++ {
			pathDir = "../" + pathDir
		}
	}
	return pathDir
}
