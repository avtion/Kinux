package services

import (
	"Kinux/core/web/models"
	"context"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestScoreExcelFile(t *testing.T) {
	// 从存档中获取
	res, err := models.ListScoreSave(context.Background(), models.ScoreTypeMission)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Fatal("无存档数据")
	}
	_missionScore := res[0].Data
	var missionScore = make(MissionScoreForAdminSlice, 0)
	if err = jsoniter.Unmarshal(_missionScore, &missionScore); err != nil {
		t.Fatal(err)
	}

	f, err := missionScore.GetExcel()
	if err != nil {
		t.Fatal(err)
	}
	// 保存一下
	if err = f.SaveAs("example_mission.xlsx"); err != nil {
		t.Fatal(err)
	}

	res, err = models.ListScoreSave(context.Background(), models.ScoreTypeExam)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Fatal("无考试存档数据")
	}
	_examScore := res[0].Data
	examScore := make(ExamScoreForAdminSlice, 0)
	if err = jsoniter.Unmarshal(_examScore, &examScore); err != nil {
		t.Fatal(err)
	}

	f, err = examScore.GetExcel()
	if err != nil {
		t.Fatal(err)
	}
	if err = f.SaveAs("example_exam.xlsx"); err != nil {
		t.Fatal(err)
	}
}
