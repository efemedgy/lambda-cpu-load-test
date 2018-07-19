#### CPU LOAD TEST FOR AWS LAMBDA

Test to compare CPU load for three different lambda functions under sleep, sqrt and goroutine packages to see the
effect of memory size increase on AWS Lambda.

You can visit [this blog post](https://medium.com/@dnamal/why-cpu-load-increases-with-memory-size-on-aws-lambda-c38235c36938) to learn more about it.

I've done these tests to calibrate Thundra's agents cpu load calculation. If you want to observe the system level metrics
 of your own lambda functions you can use [Thundra](https://www.thundra.io/) to do that for you.

If you want to run your own tests just run the script by passing a reasonable memory size as a parameter.

```shell
sh build-deploy-run.sh 512
```

This script will build all lambda functions and deploy them with `512 MB` memory size to `us-west-2` region.
You need serverless framework to use it.

You will see the new results under results folder probably overriding mine. You can avoid it by changing filenames in
the script. First parameter is the filename.
```shell
./invoker-test sqrt $mSize
```

Give it a try and let me know if you have any suggestions for improvement.