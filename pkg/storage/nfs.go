// 网络文件存储
package storage

type Nfs struct {
	Bucket   string
	Upload   string
	Download string
}

func (n *Nfs) Save(file []byte, path string) error {
	return nil
}
func (l *Nfs) Get(path string) ([]byte, error) {
	return nil, nil
}
func (l *Nfs) Delete(path string) error {
	return nil
}
func (l *Nfs) Copy(path string, targetPath string) error {
	return nil
}
