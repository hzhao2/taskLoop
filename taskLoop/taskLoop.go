package taskLoop

import (
	"fmt"
	"sync"
	"time"
)

const (
	TASKSTOP = false
	TASKRUN  = true
)

const (
	TASKROUTINE_SYNC  = true
	TASKROUTINE_ASYNC = false
)

func init() {

}

type task struct {
	taskName  string
	taskFunc  func()
	sleepTime time.Duration
	isSync    bool
}

type TaskTable struct {
	taskList   map[string]task
	taskGroup  sync.WaitGroup
	taskStatus bool
}

func InitTable() *TaskTable {
	tt := new(TaskTable)
	tt.taskList = make(map[string]task)
	return tt
}

// 新增任务事件
func (tt *TaskTable) AddTask(taskFunc func(), taskName string, sleepTime time.Duration, isSync bool) {

	var t task
	t.taskName = taskName
	t.taskFunc = taskFunc
	t.sleepTime = sleepTime
	t.isSync = isSync

	tt.add(t)
}

// 新增任务
func (t *task) run(taskGroup sync.WaitGroup, taskStatus bool) {

	defer taskGroup.Done()

	for {

		// 判断状态, 是否结束循环
		if taskStatus {
			break
		}

		//		fmt.Println("准备执行函数")

		// 不同异步
		if t.isSync {
			t.taskFunc()
		} else {
			go t.taskFunc()
		}

		//		fmt.Println("进入休眠")

		// 休眠
		time.Sleep(time.Second * t.sleepTime)
	}
}

// 新增任务
func (tt *TaskTable) add(t task) {

	// 判断是否已经有同名任务
	if _, ok := tt.taskList[t.taskName]; !ok {
		tt.taskList[t.taskName] = t
	}
}

// 开启任务循环, 进入堵塞模式
func (tt *TaskTable) Start() {

	fmt.Println("开启任务循环, 进入堵塞模式")
	tt.taskStatus = TASKRUN
	for _, v := range tt.taskList {
		tt.taskGroup.Add(1)
		// 传入 taskGroup 和 status
		go func(t task) {
			defer tt.taskGroup.Done()

			for {

				if !tt.taskStatus {
					break
				}

				fmt.Println("准备执行函数")
				// 判断状态, 是否结束循环
				t.taskFunc()

				// 不同异步

				//				fmt.Println("进入休眠")

				// 休眠
				time.Sleep(time.Second * t.sleepTime)
			}
		}(v)
	}
	tt.taskGroup.Wait()

	fmt.Println("任务循环结束, 优雅退出程序")
}

// 停止任务循环
func (tt *TaskTable) Stop() {
	fmt.Println("开始结束任务循环")
	tt.taskStatus = TASKSTOP
}
