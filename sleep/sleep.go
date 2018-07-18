package main

import (
	"time"
	"cpuload"

	"github.com/aws/aws-lambda-go/lambda"
)

func SleepCPULoadTestLambda() (cpuload.CpuTimesStatComparer, error) {
	before := cpuload.Sample()
	time.Sleep(time.Second * 3)
	after := cpuload.Sample()

	rs := cpuload.CpuTimesStatComparer{
		Before: before,
		After:  after,
	}
	return rs, nil
}

func main() {
	lambda.Start(SleepCPULoadTestLambda)
}
