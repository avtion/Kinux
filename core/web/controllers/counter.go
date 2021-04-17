package controllers

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"Kinux/core/web/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CounterRes struct {
	Account    int64 `json:"account"`
	Department int64 `json:"department"`
	Deployment int64 `json:"deployment"`
	Lesson     int64 `json:"lesson"`
	Mission    int64 `json:"mission"`
	Exam       int64 `json:"exam"`
	Checkpoint int64 `json:"checkpoint"`
	Session    int64 `json:"session"`
}

// Counter 数据统计
func Counter(c *gin.Context) {
	db := models.GetGlobalDB().WithContext(c)
	res := new(CounterRes)
	if _err := db.Model(new(models.Account)).Count(&res.Account).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Department)).Count(&res.Department).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Deployment)).Count(&res.Deployment).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Lesson)).Count(&res.Lesson).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Mission)).Count(&res.Mission).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Exam)).Count(&res.Exam).Error; _err != nil {
		logrus.Error(_err)
	}
	if _err := db.Model(new(models.Checkpoint)).Count(&res.Checkpoint).Error; _err != nil {
		logrus.Error(_err)
	}
	res.Session = services.GetScheduleCenterCount()
	c.JSON(http.StatusOK, msg.BuildSuccess(res))
}
