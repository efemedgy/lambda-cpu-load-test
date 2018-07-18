package main

import (
	"cpuload"
	"time"
	"runtime"

	"github.com/aws/aws-lambda-go/lambda"
)

func GoroutineCPULoadTestLambda() (cpuload.CpuTimesStatComparer, error) {
	before := cpuload.Sample()
	done := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}
	time.Sleep(time.Second * 3)
	close(done)
	after := cpuload.Sample()

	rs := cpuload.CpuTimesStatComparer{
		Before: before,
		After:  after,
	}
	return rs, nil
}

func main() {
	lambda.Start(GoroutineCPULoadTestLambda)
}
