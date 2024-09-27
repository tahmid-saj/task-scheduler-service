package sqs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func ListQueues() ([]*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.ListQueues(nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var resQueueUrls []*string
	resQueueUrls = append(resQueueUrls, result.QueueUrls...)

	return resQueueUrls, nil
}

func CreateQueue(queueName string) (*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"DelaySeconds": aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result.QueueUrl, nil
}

func GetURLOfQueue(queueName string) (*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result.QueueUrl, nil
}

func DeleteQueue(queueName string) (*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result.QueueUrl, nil
}