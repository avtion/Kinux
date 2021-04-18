package services

import (
	"Kinux/core/web/models"
	"Kinux/core/web/msg"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sort"
	"time"
)

// NewScoreListener 成绩监听器
func NewScoreListener(account *models.Account, lesson *models.Lesson, exam *models.Exam, mc *models.Mission,
	container string, cps ...*models.Checkpoint) (opt WsPtyWrapperOption) {
	// 校验参数
	if account == nil || account.ID == 0 {
		return nil
	}
	if lesson == nil || lesson.ID == 0 {
		return nil
	}
	if mc == nil || mc.ID == 0 {
		return nil
	}
	if container == "" {
		return nil
	}
	if len(cps) == 0 {
		return nil
	}
	var examID uint
	if exam != nil && exam.ID != 0 {
		examID = exam.ID
	}

	// 定义监听器
	var (
		inOpt, outOpt       WsPtyWrapperOption
		inReader, outReader *TerminalListener
		inMap, outMap       map[string]func(w *WsPtyWrapper) (err error)
	)

	// 回调函数
	newCallbackFn := func(cp *models.Checkpoint) func(w *WsPtyWrapper) (err error) {
		return func(w *WsPtyWrapper) (err error) {
			if err = models.GetGlobalDB().WithContext(w.ChildCtx).Create(&models.Score{
				Model:      gorm.Model{},
				Account:    account.ID,
				Lesson:     lesson.ID,
				Exam:       examID,
				Mission:    mc.ID,
				Checkpoint: cp.ID,
				Container:  container,
			}).Error; err != nil {
				_ = w.ws.SendMsg(msg.BuildFailed(fmt.Sprintf("考点检验失败: %s", err)))
				return
			}
			return w.ws.SendMsg(msg.BuildSuccess(fmt.Sprintf("考点已完成: %s", cp.Name)))
		}
	}

	for _, cp := range cps {
		switch cp.Method {
		case models.MethodExec:
			// 初始化
			if inOpt == nil || inReader == nil {
				inOpt, inReader = NewPtyWrapperListenerOpt(ListenStdin)
			}
			if inMap == nil {
				inMap = make(map[string]func(w *WsPtyWrapper) (err error))
			}
			inMap[cp.In] = newCallbackFn(cp)
		case models.MethodStdout:
			// 初始化
			if outOpt == nil || outReader == nil {
				outOpt, outReader = NewPtyWrapperListenerOpt(ListenStdin)
			}
			if outMap == nil {
				outMap = make(map[string]func(w *WsPtyWrapper) (err error))
			}
			outMap[cp.Out] = newCallbackFn(cp)
		case models.MethodTargetPort:
			// TODO 支持
		}
	}

	// 监听方法
	listenFn := func(listener *TerminalListener, w *WsPtyWrapper,
		callbackMap map[string]func(w *WsPtyWrapper) (err error)) {
		// 比较器
		dmp := diffmatchpatch.New()
		for {
			select {
			case <-listener.Context.Done():
				return
			default:
			}
			line, _err := listener.Reader.ReadLine()
			if _err != nil {
				logrus.WithField("err", _err).Error("监听器发生错误")
				return
			}
			if callback, isExist := callbackMap[line]; isExist {
				if _err := callback(w); _err != nil {
					logrus.WithField("err", _err).Error("监听器执行回调函数失败")
					continue
				}
				delete(callbackMap, line)
			} else {
				for k := range callbackMap {
					diffs := dmp.DiffMain(line, k, false)
					logrus.WithField("diff", dmp.DiffPrettyText(diffs)).Trace("指令差异")
				}
			}
		}
	}

	return CombineWsPtyWrapperOptions(inOpt, outOpt, func(w *WsPtyWrapper) {
		// 输入和输出监听
		for _, tmp := range []struct {
			l *TerminalListener
			m map[string]func(w *WsPtyWrapper) (err error)
		}{{l: inReader, m: inMap}, {l: outReader, m: outMap}} {
			if tmp.l != nil {
				go listenFn(tmp.l, w, tmp.m)
			}
		}
	})
}

// ScoreItem 成绩单项
type ScoreItem struct {
	CheckpointID    uint   `json:"checkpoint_id"`    // 考点ID
	Percent         uint   `json:"percent"`          // 该考点所占成绩比例
	IsFinish        bool   `json:"is_finish"`        // 是否未完成
	FinishTime      int64  `json:"finish_time"`      // 完成时间
	TargetContainer string `json:"target_container"` // 目标容器
	CheckpointName  string `json:"checkpoint_name"`  // 考点名称
	CheckpointDesc  string `json:"checkpoint_desc"`  // 考点描述
}

// MissionScore 实验成绩
type MissionScore struct {
	MissionID          uint         `json:"mission_id"`           // 实验ID
	MissionName        string       `json:"mission_name"`         // 实验名
	MissionDesc        string       `json:"mission_desc"`         // 实验描述
	FinishScoreCounter uint         `json:"finish_score_counter"` // 完成考点数量
	AllScoreCounter    uint         `json:"all_score_counter"`    // 全部考点数量
	Score              float64      `json:"score"`                // 成绩
	ScoreDetails       []*ScoreItem `json:"score_details"`        // 成绩详情
	Total              uint         `json:"total"`                // 总分
}

// ExamScore 考试成绩
type ExamScore struct {
	ExamID        uint    `json:"exam_id"`       // 考试ID
	ExamName      string  `json:"exam_name"`     // 考试名
	ExamDesc      string  `json:"exam_desc"`     // 考试描述
	ExamBeginAt   int64   `json:"exam_begin_at"` // 用户开始考试的时间
	ExamEndAt     int64   `json:"exam_end_at"`   // 用户结束考试的时间
	Score         float64 `json:"score"`         // 用户的成绩
	MissionScores []*struct {
		*MissionScore
		Percent uint `json:"percent"`
	} `json:"mission_scores"` // 实验详情
	Total uint `json:"total"` // 考试总成绩
}

// GetMissionScore 获取实验成绩
func GetMissionScore(c *gin.Context, accountID, lessonID, missionID, examID uint) (res *MissionScore, err error) {
	// 获取用户信息
	ac, err := models.GetAccountByID(c, int(accountID))
	if err != nil {
		return
	}

	// 获取课程信息
	lesson, err := models.GetLesson(c, lessonID)
	if err != nil {
		return
	}

	// 获取实验信息
	mission, err := models.GetMission(c, missionID)
	if err != nil {
		return
	}

	var exam = new(models.Exam)
	if examID != 0 {
		exam, _ = models.GetExam(c, examID)
	}

	// 初始化结果
	res = &MissionScore{
		MissionID:          mission.ID,
		MissionName:        mission.Name,
		MissionDesc:        mission.Desc,
		FinishScoreCounter: 0,
		AllScoreCounter:    0,
		Score:              0,
		ScoreDetails:       make([]*ScoreItem, 0),
	}

	// 找到实验所有的考点
	allMCps, err := models.FindAllMissionCheckpoints(c, mission.ID)
	if err != nil {
		return
	}
	if res.AllScoreCounter = uint(len(allMCps)); res.AllScoreCounter == 0 {
		return
	}

	// 考点信息
	var cpIDs = make([]uint, 0, len(allMCps))
	for _, v := range allMCps {
		cpIDs = append(cpIDs, v.CheckPoint)
	}
	cps, err := models.FindCheckpoints(c, cpIDs...)
	if err != nil {
		return
	}
	var cpsMapping = make(map[uint]*models.Checkpoint, len(cps))
	for k, v := range cps {
		cpsMapping[v.ID] = cps[k]
	}

	// 找到已经完成的成绩
	finishScores, err := models.FindAllAccountFinishScores(c, ac.ID, lesson.ID, exam.ID, mission.ID)
	if err != nil {
		return
	}
	var finishScoresMapping = make(map[string]map[uint]*models.Score) // map[容器名]map[考点ID]成绩
	for k, v := range finishScores {
		if _, isExist := finishScoresMapping[v.Container]; !isExist {
			finishScoresMapping[v.Container] = make(map[uint]*models.Score)
		}
		finishScoresMapping[v.Container][v.Checkpoint] = finishScores[k]
	}
	res.FinishScoreCounter = uint(len(finishScores))

	for _, cp := range allMCps {
		// 找到对应的成绩
		var (
			_score     *models.Score
			isCpFinish bool
		)
		if _, isExist := finishScoresMapping[cp.TargetContainer]; isExist {
			_score, isCpFinish = finishScoresMapping[cp.TargetContainer][cp.CheckPoint]
		}

		_cp, _ := cpsMapping[cp.CheckPoint]

		// 生成详情结果
		var item = &ScoreItem{
			CheckpointID: cp.ID,
			Percent:      cp.Percent,
			IsFinish:     isCpFinish,
			FinishTime: func() int64 {
				if isCpFinish {
					return _score.CreatedAt.Unix()
				}
				return 0
			}(),
			TargetContainer: cp.TargetContainer,
			CheckpointName:  _cp.Name,
			CheckpointDesc:  _cp.Desc,
		}

		// 添加成绩
		if isCpFinish {
			res.Score += float64(mission.Total) * (float64(cp.Percent) / 100)
		}
		res.Total = mission.Total

		// 追加结果
		res.ScoreDetails = append(res.ScoreDetails, item)
	}
	return
}

// GetExamScore 获取考试成绩
func GetExamScore(c *gin.Context, accountID, lessonID, examID uint) (res *ExamScore, err error) {
	// 获取用户信息
	ac, err := models.GetAccountByID(c, int(accountID))
	if err != nil {
		return
	}

	// 获取课程信息
	lesson, err := models.GetLesson(c, lessonID)
	if err != nil {
		return
	}

	// 获取考试信息
	exam, err := models.GetExam(c, examID)
	if err != nil {
		return
	}

	// 获取用户的考试信息
	var eLog = new(models.ExamLog)
	if err = models.GetGlobalDB().WithContext(c).Where(
		"account = ? AND exam = ?", ac.ID, exam.ID).First(eLog).Error; err != nil {
		return
	}

	// 考试的实验
	examMissions, err := models.ListExamMissions(c, exam.ID, 0, nil)
	if err != nil {
		return
	}

	// 初始化结果
	res = &ExamScore{
		ExamID:      exam.ID,
		ExamName:    exam.Name,
		ExamDesc:    exam.Desc,
		ExamBeginAt: eLog.CreatedAt.Unix(),
		ExamEndAt:   eLog.EndAt.Unix(),
		Score:       0,
		MissionScores: make([]*struct {
			*MissionScore
			Percent uint `json:"percent"`
		}, 0),
		Total: exam.Total,
	}

	for _, v := range examMissions {
		// 查询对应实验的成绩
		var ms *MissionScore
		ms, err = GetMissionScore(c, ac.ID, lesson.ID, v.Mission, exam.ID)
		if err != nil {
			return
		}

		// 按照比例加入考试成绩结果
		res.Score += ms.Score * (float64(v.Percent) / 100)

		// 追加结果
		res.MissionScores = append(res.MissionScores, &struct {
			*MissionScore
			Percent uint `json:"percent"`
		}{
			MissionScore: ms,
			Percent:      v.Percent,
		})
	}

	return
}

// MissionScoreForAdmin 班级级别实验成绩
type MissionScoreForAdmin struct {
	*MissionScore
	Pos uint `json:"pos"` // 排名

	ID           uint   `json:"id"`
	Role         uint   `json:"role"`
	Profile      uint   `json:"profile"`
	Username     string `json:"username"`
	RealName     string `json:"real_name"`
	Department   string `json:"department"`
	DepartmentId uint   `json:"department_id"`
}

// ExamScoreForAdmin 班级级别考试成绩
type ExamScoreForAdmin struct {
	*ExamScore
	Pos uint `json:"pos"` // 排名

	ID           uint   `json:"id"`
	Role         uint   `json:"role"`
	Profile      uint   `json:"profile"`
	Username     string `json:"username"`
	RealName     string `json:"real_name"`
	Department   string `json:"department"`
	DepartmentId uint   `json:"department_id"`
}

type MissionScoreForAdminSlice []*MissionScoreForAdmin
type ExamScoreForAdminSlice []*ExamScoreForAdmin

type ExcelCreator interface {
	GetExcel() (f *excelize.File, err error)
}

var _ ExcelCreator = (MissionScoreForAdminSlice)(nil)
var _ ExcelCreator = (ExamScoreForAdminSlice)(nil)

// GetMissionScoreForAdmin 管理员获取实验成绩
func GetMissionScoreForAdmin(c *gin.Context, department, lessonID, missionID uint) (
	res MissionScoreForAdminSlice, err error) {
	// 查找班级所有成员
	accounts, err := models.ListAccountsWithProfiles(c, nil, func(db *gorm.DB) *gorm.DB {
		return db.Where("department_id = ?", department)
	})
	if err != nil {
		return
	}

	// 初始化结果
	res = make([]*MissionScoreForAdmin, 0)

	// 逐个获取对应的课程成绩
	for _, ac := range accounts {
		// 查询成绩
		ms, _err := GetMissionScore(c, ac.ID, lessonID, missionID, 0)
		if _err != nil {
			ms = new(MissionScore)
		}
		res = append(res, &MissionScoreForAdmin{
			MissionScore: ms,
			Pos:          0,

			ID:           ac.ID,
			Role:         ac.Role,
			Profile:      ac.Profile,
			Username:     ac.Username,
			RealName:     ac.RealName,
			Department:   ac.Department,
			DepartmentId: ac.DepartmentId,
		})
	}

	// 技术排名
	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	for k := range res {
		res[k].Pos = uint(k + 1)
	}
	return
}

// GetExamScoreForAdmin 管理员获取考试成绩
func GetExamScoreForAdmin(c *gin.Context, department, lessonID, examID uint) (
	res ExamScoreForAdminSlice, err error) {
	// 查找班级所有成员
	accounts, err := models.ListAccountsWithProfiles(c, nil, func(db *gorm.DB) *gorm.DB {
		return db.Where("department_id = ?", department)
	})
	if err != nil {
		return
	}

	// 初始化结果
	res = make([]*ExamScoreForAdmin, 0)

	// 逐个获取对应的课程成绩
	for _, ac := range accounts {
		es, _err := GetExamScore(c, ac.ID, lessonID, examID)
		if _err != nil {
			es = new(ExamScore)
		}
		res = append(res, &ExamScoreForAdmin{
			ExamScore: es,
			Pos:       0,

			ID:           ac.ID,
			Role:         ac.Role,
			Profile:      ac.Profile,
			Username:     ac.Username,
			RealName:     ac.RealName,
			Department:   ac.Department,
			DepartmentId: ac.DepartmentId,
		})
	}

	// 技术排名
	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	for k := range res {
		res[k].Pos = uint(k + 1)
	}
	return
}

const __defaultSheetName = "Sheet1"

// GetExcel 生成Excel文件
func (ms MissionScoreForAdminSlice) GetExcel() (f *excelize.File, err error) {
	// 创建文件
	f = excelize.NewFile()

	// 顶部标题
	for k, v := range [...]string{"排名", "班级", "用户名", "姓名", "实验", "成绩", "完成考点数"} {
		if err = f.SetCellStr(__defaultSheetName,
			fmt.Sprintf("%s%s", string(rune('A'+k)), "1"),
			v,
		); err != nil {
			return
		}
	}

	// 写入每行数据
	for i := 0; i < len(ms); i++ {
		var yIndex = 2 + i
		for k, v := range [...]interface{}{
			ms[i].Pos,
			ms[i].Department,
			ms[i].Username,
			ms[i].RealName,
			ms[i].MissionName,
			ms[i].Score,
			ms[i].FinishScoreCounter,
		} {
			if err = f.SetCellValue(
				__defaultSheetName,
				fmt.Sprintf("%s%d", string(rune('A'+k)), yIndex),
				v,
			); err != nil {
				return
			}
		}
	}
	return
}

// GetExcel 生成Excel文件
func (ex ExamScoreForAdminSlice) GetExcel() (f *excelize.File, err error) {
	// 创建文件
	f = excelize.NewFile()

	// 顶部标题
	for k, v := range [...]string{"排名", "班级", "用户名", "姓名", "考试", "成绩", "考试开始时间", "考试结束时间"} {
		if err = f.SetCellStr(__defaultSheetName,
			fmt.Sprintf("%s%s", string(rune('A'+k)), "1"),
			v,
		); err != nil {
			return
		}
	}

	// 写入每行数据
	for i := 0; i < len(ex); i++ {
		var yIndex = 2 + i
		for k, v := range [...]interface{}{
			ex[i].Pos,
			ex[i].Department,
			ex[i].Username,
			ex[i].RealName,
			ex[i].ExamName,
			ex[i].Score,
			time.Unix(ex[i].ExamBeginAt, 0).Format("2006-01-02 15:04:05"),
			time.Unix(ex[i].ExamEndAt, 0).Format("2006-01-02 15:04:05"),
		} {
			if err = f.SetCellValue(
				__defaultSheetName,
				fmt.Sprintf("%s%d", string(rune('A'+k)), yIndex),
				v,
			); err != nil {
				return
			}
		}
	}
	return
}
