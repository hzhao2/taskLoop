package signal

import (
	"fmt"
	"os"
	"os/signal"
)

type signalHandler func(s os.Signal, arg interface{})

type signalSet struct {
	m map[os.Signal]signalHandler
}

func (set *signalSet) handle(sig os.Signal, arg interface{}) (err error) {
	if _, found := set.m[sig]; found {
		set.m[sig](sig, arg)
		return nil
	} else {
		return fmt.Errorf("未注册信号 %v", sig)
	}

	panic("won't reach here")
}

func InitSignalListen() *signalSet {
	set := new(signalSet)
	set.m = make(map[os.Signal]signalHandler)
	return set
}

func (set *signalSet) RegisterSignal(s os.Signal, handler signalHandler) {
	if _, found := set.m[s]; !found {
		set.m[s] = handler
	}
}

func (set *signalSet) StartSignalListen() {

	for {
		c := make(chan os.Signal)

		signal.Notify(c)
		sig := <-c

		err := set.handle(sig, nil)
		if err != nil {
			fmt.Printf("接收信号 %v,并无对应事件 \n", sig)
			//os.Exit(1)
		}
	}
}
