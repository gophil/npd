package npd

type TaskType int

//任务类型
const (
	TASK_NORMAL TaskType = iota //普通任务
	TASK_REPORT
)

//任务结构
type Task struct {
	TaskId     string
	Type       TaskType
	TargetObj  *interface{}
	TargetFunc string
}

//创建任务
func MakeTask(taskType TaskType, targetObj interface{}, targetFunc string) Task {
	uuid, _ := GenerateUUID()
	return Task{
		TaskId:     uuid,
		Type:       taskType,
		TargetObj:  &targetObj,
		TargetFunc: targetFunc,
	}
}
