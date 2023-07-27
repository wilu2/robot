// 文件存储接口
package storage

type IStorage interface {
	Save(file []byte, path string) error
	Get(path string) ([]byte, error)
	Delete(path string) error
	Copy(path string, targetPath string) error
}
