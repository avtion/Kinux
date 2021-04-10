package services

import (
	"Kinux/core/k8s"
	"Kinux/core/web/models"
	"Kinux/tools"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

type mcOpt func(mc *MissionController) error

// MissionController 任务控制器
type MissionController struct {
	ctx     context.Context
	Ac      *models.Account
	Mission *models.Mission
	Exam    *models.Exam
	Lesson  *models.Lesson
	opts    []mcOpt

	// 容器部署配置
	dp         *models.Deployment
	dpCfg      *appV1.Deployment
	dpSelector labels.Set
}

// NewMissionController 创建任务控制器
func NewMissionController(ctx context.Context) (mc *MissionController) {
	mc = &MissionController{ctx: ctx}
	return
}

// SetAc 设置任务执行用户
func (mc *MissionController) SetAc(ac *models.Account) *MissionController {
	if ac == nil {
		return mc
	}
	mc.Ac = ac
	return mc
}

// SetExam 设置考试
func (mc *MissionController) SetExam(e *models.Exam) *MissionController {
	if e == nil || e.ID == 0 {
		return mc
	}
	mc.Exam = e
	return mc
}

// SetLesson  设置课程
func (mc *MissionController) SetLesson(l *models.Lesson) *MissionController {
	if l == nil || l.ID == 0 {
		return mc
	}
	mc.Lesson = l
	return mc
}

// SetMission 设置任务
func (mc *MissionController) SetMission(m *models.Mission) *MissionController {
	if m == nil {
		return mc
	}
	mc.Mission = m
	return mc.getDpCfg().parseDpCfg().fixNamespace().generateSelector(nil).applySelector()
}

// NewDeployment 创建新的Deployment
func (mc *MissionController) NewDeployment() (err error) {
	defer func() {
		if err != nil {
			logrus.WithField("err", err).Error("创建新的Deployment失败")
			return
		}
		logrus.WithField("deployment名称", mc.dpCfg.GetName()).Debug("创建新的Deployment成功")
	}()

	// 初始化控制器配置
	mc.generateAndApplyDpName()

	// 创建新的Deployment之前先清除正在进行的任务
	if err = mc.ClearAllMission(); err != nil {
		return
	}

	mc.appendOpt(
		func(mc *MissionController) error {
			if mc.dpCfg == nil || mc.dpSelector == nil {
				return errors.New("创建新的Deployment失败: 配置未初始化")
			}
			// 调用K8S模块方法
			_err := k8s.NewDeployment(mc.ctx, mc.dpCfg, mc.dpSelector)
			return _err
		})
	return mc.exec()
}

// ClearAllMission 清除用户正在进行的Deployment
func (mc *MissionController) ClearAllMission() (err error) {
	return mc.destroyDeployment(missionLabel, deploymentLabel, examLabel, lessonLabel)
}

// DestroyDeployment 删除用户指定的dp
func (mc *MissionController) DestroyDeployment() (err error) {
	if mc.Mission == nil || mc.Mission.ID == 0 {
		return errors.New("无法删除用户指定的任务: 任务信息未初始化")
	}
	return mc.destroyDeployment()
}

// 删除操作的内部实现
func (mc *MissionController) destroyDeployment(withoutLabels ...string) (err error) {
	mc.appendOpt(func(mc *MissionController) error {
		// FIX 排除deployment的label
		var m = mc.dpSelector
		if len(withoutLabels) > 0 {
			filter := make(map[string]struct{}, len(withoutLabels))
			for _, v := range withoutLabels {
				filter[v] = struct{}{}
			}
			var tmp = make(labels.Set, len(m))
			for k, v := range m {
				if _, isExist := filter[k]; isExist {
					continue
				}
				tmp[k] = v
			}
			m = tmp
		}

		if mc.Ac == nil {
			return errors.New("无法删除用户指定的任务: 用户信息未初始化")
		}

		dps, _err := k8s.ListDeployments(mc.ctx, mc.dpCfg.GetNamespace(), m)
		if _err != nil {
			return _err
		}
		for _, dp := range dps.Items {
			if _err = k8s.DeleteDeployment(mc.ctx, mc.dpCfg.GetNamespace(), dp.Name); _err != nil {
				return _err
			}
		}
		return nil
	})
	return mc.exec()
}

// 追加任务操作
func (mc *MissionController) appendOpt(opt mcOpt) {
	mc.opts = append(mc.opts, opt)
	return
}

// 执行操作
func (mc *MissionController) exec() (err error) {
	defer func() {
		// 当目前执行队列所有的任务执行完毕的时候应当清空队列
		// 以释放所有的函数空间，便于GC
		mc.opts = make([]mcOpt, 0, 0)
	}()
	for _, opt := range mc.opts {
		// 上下文判断结束
		select {
		case <-mc.ctx.Done():
			return mc.ctx.Err()
		default:
			if err = opt(mc); err != nil {
				return
			}
		}
	}
	return
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
		dpCfg, err := k8s.ParseDeploymentConfig(mc.dp.Raw, true)
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
		nameBuilder.WriteString(strings.ToLower(tools.GetRandomString(6)))
		mc.dpCfg.GetObjectMeta().SetName(nameBuilder.String())
		return nil
	})
	return mc
}

// 生成K8S Deployment的Selector用于检索对应的资源
// 格式：map[string]string{"accountID", "departmentID", "missionID"，other...}
func (mc *MissionController) generateSelector(other map[string]string) *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		// 创建标签生成器
		l := NewLabelMarker(3)
		if mc.Ac != nil && mc.Ac.ID != 0 {
			l.WithAccount(mc.Ac.ID)
		}
		if mc.Mission != nil {
			l.WithMission(mc.Mission.ID)
		}
		if mc.Exam != nil {
			l.WithExam(mc.Exam.ID)
		} else {
			l.WithExam(0)
		}
		if mc.Lesson != nil {
			l.WithLesson(mc.Lesson.ID)
		}
		if mc.Mission != nil && mc.Mission.Deployment != 0 {
			l.WithDeployment(mc.Mission.Deployment)
		}
		for k, v := range other {
			l.WithString(k, v)
		}
		mc.dpSelector = l.Do()

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

// GenerateAndApplyNodeSelector 特殊对外的方法, 用于在测试环境生成并应用K8S Deployment节点选择器
func (mc *MissionController) GenerateAndApplyNodeSelector(cpu string) *MissionController {
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dpCfg == nil {
			return errors.New("生成节点选择器失败: mission的部署配置信息未初始化")
		}
		mc.dpCfg.Spec.Template.Spec.NodeSelector["cpu"] = cpu
		return nil
	})
	return mc
}

// 修改Deployment配置文件的命名空间，防止错误执行
func (mc *MissionController) fixNamespace(pass ...bool) *MissionController {
	if len(pass) > 0 && pass[0] {
		return mc
	}
	mc.appendOpt(func(mc *MissionController) error {
		if mc.dpCfg != nil {
			mc.dpCfg.SetNamespace(k8s.GetDefaultNS())
		}
		return nil
	})
	return mc
}

// GetPods 获取POD列表
func (mc *MissionController) GetPods() (p *coreV1.PodList, err error) {
	// 初始化控制器配置
	mc.getDpCfg().parseDpCfg().generateAndApplyDpName().generateSelector(nil).applySelector().fixNamespace()
	if err = mc.exec(); err != nil {
		return
	}
	return k8s.GetPods(mc.ctx, mc.dpCfg.GetNamespace(), mc.dpSelector)
}

// ResetMission 重置任务的容器（即删除POD）
func (mc *MissionController) ResetMission(ns string) (err error) {
	pods, err := mc.GetPods()
	if err != nil {
		return err
	}
	if len(pods.Items) == 0 {
		return
	}
	if ns == "" {
		ns = k8s.GetDefaultNS()
	}
	for _, pod := range pods.Items {
		if err = k8s.DeletePod(mc.ctx, ns, pod.Name); err != nil {
			return
		}
	}
	return
}

// WatchDeploymentToReady 监视Deployment
func (mc *MissionController) WatchDeploymentToReady(ns string) (resCh <-chan error) {
	if ns == "" {
		ns = k8s.GetDefaultNS()
	}
	errCh := make(chan error)
	k8s.WatchDeploymentsToReady(mc.ctx, ns, mc.dpSelector, errCh)
	return errCh
}
