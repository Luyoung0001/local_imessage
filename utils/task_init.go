package utils

import "time"

type TimerFunc func(interface{}) bool

//delay: 表示初始延迟时间，即第一次执行函数前的等待时间。
//tick: 表示定时执行函数的时间间隔。
//fun: 是一个自定义的函数类型TimerFunc，它接受一个参数（interface{}类型），并返回一个bool值。
//param: 是传递给定时执行函数的参数，类型为interface{}。

func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		t := time.NewTimer(delay)
		for {
			select {
			// 监听 t
			case <-t.C:
				if fun(param) == false {
					return
				}
				// 重置 t
				t.Reset(tick)
			}
		}
	}()
}
