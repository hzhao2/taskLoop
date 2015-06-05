package main

import (
	"facenote/signal"
	"facenote/taskLoop"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

func main() {

	// 初始化任务循环
	tl := taskLoop.InitTable()

	// 初始化监听
	ss := signal.InitSignalListen()

	// 注册监听事件, 接收到指定信号, 结束当前任务循环
	handle := func(s os.Signal, arg interface{}) {
		fmt.Printf("信号注册事件执行\n")
		tl.Stop()
	}

	ss.RegisterSignal(syscall.SIGTERM, handle)
	ss.RegisterSignal(syscall.SIGINT, handle)

	// 开始监听
	go ss.StartSignalListen()

	// 注册任务循环
	tl.AddTask(test, "test", 3, taskLoop.TASKROUTINE_SYNC)

	tl.Start()
	fmt.Println("程序结束exit")
}


var testCount int = 0
var mem runtime.MemStats

func test() {

	taskStatus := true
	defer func() {

		runtime.ReadMemStats(&mem)
		testCount++
		fmt.Println(testCount, " test status", taskStatus, " ", mem.Alloc, " ", mem.TotalAlloc, " ", mem.HeapAlloc, " ", mem.HeapSys)

		/*
			fmt.Println(mem.Alloc) // 已经被分配并仍在使用的字节数
			fmt.Println(mem.TotalAlloc)// 从开始运行到现在分配的内存数
			fmt.Println(mem.HeapAlloc)// 堆当前的用量
			fmt.Println(mem.HeapSys)// 堆当前和已经被释放但尚未被释放 但尚未归还操作系统的用量
		*/
	}()

}
