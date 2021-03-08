package services

import (
	"github.com/spf13/cast"
	"k8s.io/apimachinery/pkg/labels"
)

// label标签定义
const (
	accountLabel    = "account-id"
	missionLabel    = "mission-id"
	deploymentLabel = "deployment-id"
	examLabel       = "exam-id"
)

// labels标签生成器
type labelMaker struct {
	raw labels.Set
}

// 创建标签生成器
func NewLabelMarker(size ...int) *labelMaker {
	if len(size) > 0 && size[0] > 0 {
		return &labelMaker{
			raw: make(labels.Set, size[0]),
		}
	}
	return &labelMaker{
		raw: labels.Set{},
	}
}

// 执行生成方法
func (l *labelMaker) Do() labels.Set {
	if l.raw == nil {
		return labels.Set{}
	}
	return l.raw
}

/*
	方法
*/
func (l *labelMaker) With(k, v interface{}) *labelMaker {
	l.raw[cast.ToString(k)] = l.raw[cast.ToString(v)]
	return l
}

func (l *labelMaker) WithString(k, v string) *labelMaker {
	l.raw[k] = v
	return l
}

func (l *labelMaker) WithAccount(id interface{}) *labelMaker {
	return l.WithString(accountLabel, cast.ToString(id))
}

func (l *labelMaker) WithMission(id interface{}) *labelMaker {
	return l.WithString(missionLabel, cast.ToString(id))
}

func (l *labelMaker) WithDeployment(id interface{}) *labelMaker {
	return l.WithString(deploymentLabel, cast.ToString(id))
}

func (l *labelMaker) WithExam(id interface{}) *labelMaker {
	return l.WithString(examLabel, cast.ToString(id))
}
