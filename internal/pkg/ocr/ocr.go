package ocr

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/spf13/viper"
)

var (
	once              sync.Once
	recognizeTableApi string
)

// 表格识别
func RecognizeTable(file []byte) ([]byte, error) {
	once.Do(func() {
		recognizeTableApi = viper.GetString("ocr.recognize_table_api")
	})
	req, err := http.NewRequest(http.MethodPost, recognizeTableApi, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK || err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}
