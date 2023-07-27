package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(v interface{}) string {
	var data []byte

	switch v.(type) {
	case []byte:
		data = v.([]byte)
	case string:
		data = []byte(v.(string))
	}
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}
