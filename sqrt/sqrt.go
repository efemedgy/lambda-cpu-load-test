package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"math"
	"cpuload"
)

func SqrtCPULoadTestLambda() (cpuload.CpuTimesStatComparer, error) {
	before := cpuload.Sample()

	for i := 0.0; i < 1000000000; i++ {
		math.Sqrt(i)
	}

	after := cpuload.Sample()

	rs := cpuload.CpuTimesStatComparer{
		Before: before,
		After:  after,
	}

	return rs, nil
}

func main() {
	lambda.Start(SqrtCPULoadTestLambda)
}
