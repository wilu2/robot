// 本地文件存储
package storage

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type LocalStorage struct {
	SaveDir string
}

func (l *LocalStorage) Save(file []byte, path string) error {
	filePath := filepath.Join(l.SaveDir, path)
	fileDir := filepath.Dir(filePath)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filePath, file, 0644)
}

func (l *LocalStorage) Get(path string) ([]byte, error) {
	filePath := filepath.Join(l.SaveDir, path)
	return ioutil.ReadFile(filePath)
}

func (l *LocalStorage) Delete(path string) error {
	filePath := filepath.Join(l.SaveDir, path)
	return os.Remove(filePath)
}

func (l *LocalStorage) Copy(path string, targetPath string) error {
	f := filepath.Join(l.SaveDir, path)
	t := filepath.Join(targetPath, path)
	cmd := exec.Command("cp", "-rf", f, t)
	return cmd.Run()
}
