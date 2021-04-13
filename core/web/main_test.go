package web

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools/cfg"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../")
	})
	cfg.DefaultConfig.Kubernetes.KubeConfigPath = "../../kubeConfig"
	k8s.InitKubernetes()
}

func TestInitWebService(t *testing.T) {
	tests := []struct {
		name     string
		noFinish bool
	}{
		{
			name:     "test",
			noFinish: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitWebService()
			if tt.noFinish {
				select {}
			} else {
				<-time.After(10 * time.Second)
			}
		})
	}
}

type colsInfo struct {
	Cid     int
	Name    string
	Type    string
	Notnull int
	Pk      int
}

type colsInfoTrans struct {
	Cid     int    `json:"序号"`
	Name    string `json:"字段名"`
	Type    string `json:"类型"`
	Notnull string `json:"是否为空"`
	Pk      string `json:"是否为主键"`
	Desc    string `json:"说明"`
}

func (c *colsInfo) Translate() *colsInfoTrans {
	return &colsInfoTrans{
		Cid:  c.Cid + 1,
		Name: c.Name,
		Type: c.Type,
		Notnull: func() string {
			switch c.Notnull {
			case 1:
				return "是"
			default:
				return "否"
			}
		}(),
		Pk: func() string {
			switch c.Pk {
			case 1:
				return "是"
			default:
				return "否"
			}
		}(),
	}
}

func TestGetSQLiteDBInfo(t *testing.T) {
	// 先获取表名
	rows, err := models.GetGlobalDB().Raw(
		"select name from sqlite_master where type='table'").Rows()
	if err != nil {
		t.Fatal(err)
	}
	var names = make([]string, 0)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	t.Log(names)

	var tbs = make([][]*colsInfoTrans, 0, len(names))
	for _, v := range names {
		var colsInfos = make([]*colsInfo, 0)
		if err = models.GetGlobalDB().Raw(
			`PRAGMA table_info ('` + v + `')`).Scan(&colsInfos).Error; err != nil {
			t.Fatal(err)
		}
		var colsInfosTrans = make([]*colsInfoTrans, 0, len(colsInfos))
		for _, v := range colsInfos {
			colsInfosTrans = append(colsInfosTrans, v.Translate())
		}
		tbs = append(tbs, colsInfosTrans)
	}
	//_ = os.Mkdir("./dbInfos", 0777)
	//for k, v := range tbs {
	//	data, err := jsoniter.Marshal(v)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if err = os.WriteFile(filepath.Join("dbInfos", names[k]+`.json`), data, 0777); err != nil {
	//		return
	//	}
	//}

	// 可以通过 https://tableconvert.com/ 转换成markdown表格
	const sheetName = "Sheet1"
	var colNames = [...]string{"序号", "字段名", "类型", "是否为空", "是否为主键", "说明"}
	var colsIndexMapping = map[int]string{
		0: "A",
		1: "B",
		2: "C",
		3: "D",
		4: "E",
		5: "F",
	}
	f := excelize.NewFile()
	f.NewSheet(sheetName)
	for k := range colNames {
		if err = f.SetCellStr(sheetName, colsIndexMapping[k]+"1", colNames[k]); err != nil {
			t.Fatal(err)
		}
	}
	var rowIndex = 2
	for k, colsInfosTrans := range tbs {
		if err = f.MergeCell(sheetName, "A"+cast.ToString(rowIndex), "E"+cast.ToString(rowIndex)); err != nil {
			return
		}
		if err = f.SetCellStr(sheetName, "A"+cast.ToString(rowIndex), names[k]); err != nil {
			t.Fatal(err)
		}
		rowIndex++
		for _, v := range colsInfosTrans {
			_ = f.SetCellInt(sheetName, "A"+cast.ToString(rowIndex), v.Cid)
			_ = f.SetCellStr(sheetName, "B"+cast.ToString(rowIndex), v.Name)
			_ = f.SetCellStr(sheetName, "C"+cast.ToString(rowIndex), v.Type)
			_ = f.SetCellStr(sheetName, "D"+cast.ToString(rowIndex), v.Notnull)
			_ = f.SetCellStr(sheetName, "E"+cast.ToString(rowIndex), v.Pk)
			rowIndex++
		}
	}
	if err = f.SaveAs("dbInfo.xlsx"); err != nil {
		t.Fatal(err)
	}
}
