package lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func ListFunctions(region string) ([]*lambda.FunctionConfiguration, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess, &aws.Config{Region: aws.String(region)})

	result, err := svc.ListFunctions(nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result.Functions, nil
}

func CreateFunction(region, zipFile, bucket, functionName, handler, roleARN, runtime string) (*lambda.FunctionConfiguration, error) {
	if zipFile == "" || bucket == "" || functionName == "" || handler == "" || roleARN == "" || runtime == "" {
    return nil, fmt.Errorf("You must supply a zip file name, bucket name, function name, handler (package) name, role ARN, and runtime value.")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess, &aws.Config{Region: aws.String(region)})

	contents, err := os.Open(zipFile)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	zipFileContents, err := json.Marshal(contents)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	createCode := &lambda.FunctionCode{
		S3Bucket:        aws.String(bucket),
		S3Key:           aws.String(zipFile),
		S3ObjectVersion: aws.String("1"),
		ZipFile: zipFileContents,
	}

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String(functionName),
		Handler:      aws.String(handler),
		Role:         aws.String(roleARN),
		Runtime:      aws.String(runtime),
	}

	result, err := svc.CreateFunction(createArgs)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}

func RunFunction(region string) (*lambda.InvokeOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String(region)})

	request := struct{
		time string `json:"time"`
		sorting string `json:"sorting"`
		timing int `json"timing"`
	}{
		"time",
		"descending",
		10,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("MyGetItemsFunction"), Payload: payload})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}