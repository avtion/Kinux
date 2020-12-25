package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools"
	"context"
	"errors"
	"github.com/spf13/cast"
	appV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

type mcOpt func(mc *MissionController) error

// 任务控制器
type MissionController struct {
	ctx     context.Context
	Ac      *models.Account
	Mission *models.Mission
	opts    []mcOpt

	// 容器部署配置
	dp         *models.Deployment
	dpCfg      *appV1.Deployment
	dpSelector labels.Set
}

// 创建任务控制器
func NewMissionController(ctx context.Context) (mc *MissionController) {
	mc = &MissionController{ctx: ctx}
	return
}

// 追加任务操作
func (mc *MissionController) appendOpt(opt mcOpt) {
	mc.opts = append(mc.opts, opt)
	return
}

// 执行操作
func (mc *MissionController) exec() (err error) {
	for _, opt := range mc.opts {
		// 上下文判断结束
		select {
		case <-mc.ctx.Done():
			return mc.ctx.Err()
		default:
			if err = opt(mc); err != nil {

			}
			return
		}
	}
	return
}

// 设置任务执行用户
func (mc *MissionController) SetAc(ac *models.Account) *MissionController {
	if ac == nil {
		return mc
	}
	mc.appendOpt(func(mc *MissionController) error {
		mc.Ac = ac
		return nil
	})
	return mc
}

// 设置任务
func (mc *MissionController) SetMission(m *models.Mission) *MissionController {
	if m == nil {
		return mc
	}
	mc.appendOpt(func(mc *MissionController) error {
		mc.Mission = m
		return nil
	})
	return mc
}

// 获取K8S Deployment部署配置
func (mc *MissionController) getDpCfg() *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.Mission == nil {
			return errors.New("mission is null")
		}
		if mc.Mission.Deployment == 0 {
			return errors.New("mission没有对应的deployment")
		}

		// 从数据库中获取Deployment
		dp, err := models.GetDeployment(mc.ctx, mc.Mission.Deployment)
		if err != nil {
			return err
		}
		mc.dp = dp
		return nil
	})
	return mc
}

// 解析K8S Deployment
func (mc *MissionController) parseDpCfg() *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dp == nil {
			return errors.New("mission的deployment为空")
		}
		dpCfg, err := k8s.ParseDeploymentConfig(mc.ctx, mc.dp.Raw, true)
		if err != nil {
			return err
		}
		mc.dpCfg = dpCfg
		return nil
	})
	return mc
}

// 生成并应用K8S Deployment的资源名
// 格式: {DeploymentName}-{AccountName}-{随机6位字符}
func (mc *MissionController) generateAndApplyDpName() *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dpCfg == nil {
			return errors.New("生成部署配置失败: mission没有对应的Deployment配置信息")
		}
		if mc.Ac == nil {
			return errors.New("生成部署配置失败: mission没有对应的用户信息")
		}
		nameBuilder := strings.Builder{}
		nameBuilder.WriteString(mc.dpCfg.GetObjectMeta().GetName() + "-")
		nameBuilder.WriteString(mc.Ac.Username + "-")
		nameBuilder.WriteString(tools.GetRandomString(6))
		mc.dpCfg.GetObjectMeta().SetName(nameBuilder.String())
		return nil
	})
	return mc
}

// 生成K8S Deployment的Selector用于检索对应的资源
// 格式：map[string]string{"accountID", "departmentID", "missionID"，other...}
func (mc *MissionController) generateSelector(other map[string]string) *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		mc.dpSelector = make(labels.Set, 3)
		if mc.Ac != nil && mc.Ac.ID != 0 {
			mc.dpSelector["accountID"] = cast.ToString(mc.Ac.ID)
		}
		if mc.Mission != nil {
			mc.dpSelector["missionID"] = cast.ToString(mc.Mission.ID)
		}
		if mc.Mission != nil && mc.Mission.Deployment != 0 {
			mc.dpSelector["deploymentID"] = cast.ToString(mc.Mission.Deployment)
		}
		for k, v := range other {
			mc.dpSelector[k] = v
		}
		return nil
	})
	return mc
}

// 应用K8S Deployment选择器，在此之前应该生成对应的选择器
func (mc *MissionController) applySelector() *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dpCfg == nil {
			return errors.New("应用选择器失败: mission未载入部署配置信息")
		}
		if mc.dpSelector == nil {
			return errors.New("应用选择器失败: mission未生成选择器")
		}

		// 应用
		mc.dpCfg.SetLabels(mc.dpSelector)
		return nil
	})
	return mc
}

// 生成并应用K8S Deployment节点选择器用于测试环境
func (mc *MissionController) generateAndApplyNodeSelector(cpu string) *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dpCfg == nil {
			return errors.New("生成节点选择器失败: mission的部署配置信息未初始化")
		}
		mc.dpCfg.Spec.Template.Spec.NodeSelector["cpu"] = cpu
		return nil
	})
	return mc
}
