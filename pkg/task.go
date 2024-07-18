package task

import (
	"encoding/json"
	"fmt"
	"github.com/lhdhtrc/func-go/file"
	"path/filepath"
	"reflect"
)

func (core *CoreEntity) InitCert(dir string, config interface{}) {
	dirPath := filepath.Join("dep", "cert", dir)

	// 遍历config的字段
	valueOfConfig := reflect.ValueOf(config).Elem()
	typeOfConfig := valueOfConfig.Type()
	for i := 0; i < valueOfConfig.NumField(); i++ {
		fieldValue := valueOfConfig.Field(i)
		fieldType := typeOfConfig.Field(i)
		if fieldValue.IsValid() && !fieldValue.IsZero() && fieldType.Type.Kind() == reflect.String {
			remote := fieldValue.String()

			// 分割路径，得到文件名部分
			f := filepath.Base(remote)
			local := filepath.Join(dirPath, f)
			fieldValue.SetString(local)

			core.Add(&RawEntity{
				Id: fmt.Sprintf("INIT_CERT_%d", i),
				Handle: func() {
					read, err := file.ReadRemote(remote)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					_ = file.WriteLocal(local, read)
				},
			})
		}
	}
}

func (core *CoreEntity) InitConfig(source []string, config []interface{}) {
	for i, it := range source {
		core.Add(&RawEntity{
			Id: fmt.Sprintf("INIT_CONFIG_%d", i),
			Handle: func() {
				bytes, err := file.ReadRemote(it)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				_ = json.Unmarshal(bytes, config[i])
				return
			},
		})
	}
}
