package main

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"fmt"
	"os"
	"log"
	"cpuload"
	"encoding/json"
)

func main() {
	fName := os.Args[1] + "-" + os.Args[2]

	// Create a csv file and add headers
	file, err := os.Create("./results/" + fName + ".csv")
	checkError("Cannot create file", err)
	defer file.Close()
	file.WriteString("dUtime,dStime,dUser,dSystem,dIdle,dNice,dIowait,dIrq,dSoftirq,dSteal,dGuest,dGuest_nice,sys_usage,proc_usage\n")

	for i := 0; i < 100; i++ {
		result := invokeLambda(fName)
		parseResultsAndWriteToCsv(file, result)
	}
}

func parseResultsAndWriteToCsv(file *os.File, result *lambda.InvokeOutput) {

	stats := cpuload.CpuTimesStatComparer{}
	json.Unmarshal(result.Payload, &stats)

	// Calculate delta time for all
	dUtime := stats.After.Utime - stats.Before.Utime
	dStime := stats.After.Stime - stats.Before.Stime
	dUser := stats.After.User - stats.Before.User
	dSystem := stats.After.System - stats.Before.System
	dIdle := stats.After.Idle - stats.Before.Idle
	dNice := stats.After.Nice - stats.Before.Nice
	dIowait := stats.After.Iowait - stats.Before.Iowait
	dIrq := stats.After.Irq - stats.Before.Irq
	dSoftirq := stats.After.SoftIrq - stats.Before.SoftIrq
	dSteal := stats.After.Steal - stats.Before.Steal
	dGuest := stats.After.Guest - stats.Before.Guest
	dGuest_nice := stats.After.Guest_Nice - stats.Before.Guest_Nice

	sys_usage := cpuload.CalculateSystemUsagePercent(&stats)
	proc_usage := cpuload.CalculateProcessUsagePercent(&stats)

	// Convert to csv format
	row := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		dUtime, dStime, dUser, dSystem, dIdle, dNice, dIowait, dIrq, dSoftirq, dSteal, dGuest, dGuest_nice, sys_usage, proc_usage)
	file.WriteString(row)
}

func invokeLambda(lambdaFunctionName string) *lambda.InvokeOutput {
	cfg, err := external.LoadDefaultAWSConfig()
	cfg.Region = "us-west-2"
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	svc := lambda.New(cfg)
	input := &lambda.InvokeInput{
		FunctionName: aws.String("CPULoadTest-dev-" + lambdaFunctionName),
	}
	req := svc.InvokeRequest(input)
	result, err := req.Send()
	if err != nil {
		handleError(err)
	}
	return result
}

func handleError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case lambda.ErrCodeServiceException:
			fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
		case lambda.ErrCodeResourceNotFoundException:
			fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRequestContentException:
			fmt.Println(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
		case lambda.ErrCodeRequestTooLargeException:
			fmt.Println(lambda.ErrCodeRequestTooLargeException, aerr.Error())
		case lambda.ErrCodeUnsupportedMediaTypeException:
			fmt.Println(lambda.ErrCodeUnsupportedMediaTypeException, aerr.Error())
		case lambda.ErrCodeTooManyRequestsException:
			fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
		case lambda.ErrCodeInvalidParameterValueException:
			fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
		case lambda.ErrCodeEC2UnexpectedException:
			fmt.Println(lambda.ErrCodeEC2UnexpectedException, aerr.Error())
		case lambda.ErrCodeSubnetIPAddressLimitReachedException:
			fmt.Println(lambda.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
		case lambda.ErrCodeENILimitReachedException:
			fmt.Println(lambda.ErrCodeENILimitReachedException, aerr.Error())
		case lambda.ErrCodeEC2ThrottledException:
			fmt.Println(lambda.ErrCodeEC2ThrottledException, aerr.Error())
		case lambda.ErrCodeEC2AccessDeniedException:
			fmt.Println(lambda.ErrCodeEC2AccessDeniedException, aerr.Error())
		case lambda.ErrCodeInvalidSubnetIDException:
			fmt.Println(lambda.ErrCodeInvalidSubnetIDException, aerr.Error())
		case lambda.ErrCodeInvalidSecurityGroupIDException:
			fmt.Println(lambda.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
		case lambda.ErrCodeInvalidZipFileException:
			fmt.Println(lambda.ErrCodeInvalidZipFileException, aerr.Error())
		case lambda.ErrCodeKMSDisabledException:
			fmt.Println(lambda.ErrCodeKMSDisabledException, aerr.Error())
		case lambda.ErrCodeKMSInvalidStateException:
			fmt.Println(lambda.ErrCodeKMSInvalidStateException, aerr.Error())
		case lambda.ErrCodeKMSAccessDeniedException:
			fmt.Println(lambda.ErrCodeKMSAccessDeniedException, aerr.Error())
		case lambda.ErrCodeKMSNotFoundException:
			fmt.Println(lambda.ErrCodeKMSNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRuntimeException:
			fmt.Println(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
