package models

import (
	"context"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestExamLogsCreate(t *testing.T) {
	eLog := &ExamLog{
		Model:   gorm.Model{},
		Account: 3,
		Exam:    4,
		TickAt:  time.Time{},
		EndAt:   time.Time{},
	}
	err := GetGlobalDB().WithContext(context.Background()).Create(eLog).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Log(eLog)
	t.Log(eLog.CreatedAt)
	t.Log(eLog.TickAt.IsZero())
	t.Log(eLog.EndAt.IsZero())
}
