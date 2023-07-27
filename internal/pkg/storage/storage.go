package storage

import (
	"errors"
	genericOptions "financial_statement/internal/pkg/options"
	s "financial_statement/pkg/storage"
	"sync"
)

var (
	FileStorage s.IStorage
	once        sync.Once
)

func NewFileStore(opts genericOptions.IStorageOptions) (fileStorage s.IStorage, err error) {
	if opts == nil && FileStorage == nil {
		return nil, errors.New("failed to get store fatory")
	}
	once.Do(func() {
		switch o := opts.(type) {
		case *genericOptions.NfsStorageOptions:
			o = opts.(*genericOptions.NfsStorageOptions)
			FileStorage = &s.Nfs{
				Bucket:   o.Bucket,
				Upload:   o.Upload,
				Download: o.Download,
			}
		case *genericOptions.LocalStorageOptions:
			o = opts.(*genericOptions.LocalStorageOptions)
			FileStorage = &s.LocalStorage{
				SaveDir: o.SaveDir,
			}
		}
	})
	return FileStorage, nil
}
