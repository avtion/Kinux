package web

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools/bytesconv"
	"Kinux/tools/cfg"
	"context"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/yanyiwu/gojieba"
	"os"
	"path/filepath"
	"strings"
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

// 导出SQLite测试数据
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

// 添加测试实验
func TestAddMissions(t *testing.T) {
	var basePath = filepath.Join("../../", "example_configs", "command")
	// 查找命令行说明目录
	entries, err := os.ReadDir(basePath)
	if err != nil {
		t.Fatal(err)
	}

	// 查找centos实验配置
	dps, err := models.ListDeployment(context.Background(), "centos", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(dps) == 0 {
		t.Fatal("没有centos实验配置")
	}
	dp := dps[0]

	containers, err := dp.ParseDeploymentContainerNames()
	if err != nil {
		t.Fatal(err)
	}

	x := gojieba.NewJieba()
	defer x.Free()

	for _, f := range entries {
		rawGuide, err := os.ReadFile(filepath.Join(basePath, f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		missionName := strings.Split(f.Name(), ".")[0]
		guide := bytesconv.BytesToString(rawGuide)

		keywordsRaw := x.ExtractWithWeight(guide, 5)
		kvs := make([]string, 0, len(keywordsRaw))
		for _, v := range keywordsRaw {
			kvs = append(kvs, v.Word)
		}

		if err = models.AddMission(context.Background(), missionName, dp.ID,
			models.MissionOptDesc(strings.Join(kvs, " ")),
			models.MissionOptTotal(100),
			models.MissionOptDeployment("", containers[0], containers),
			func(m *models.Mission) (err error) {
				m.Guide = guide
				return nil
			},
		); err != nil {
			t.Fatal(err)
		}
	}
	t.Log("实验导入成功")
}
