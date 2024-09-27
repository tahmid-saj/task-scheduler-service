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

func SendMessage(message, messageTitle, messageAuthor, queueURL string) (*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	_, err := svc.SendMessage(&sqs.SendMessageInput{
    DelaySeconds: aws.Int64(10),
    MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(messageTitle),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(messageAuthor),
			},
    },
    MessageBody: aws.String(message),
    QueueUrl:    aws.String(queueURL),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &queueURL, nil
}

func ReceiveMessage(queueName string, timeout int64) (*sqs.ReceiveMessageOutput, error) {
	if timeout < 0 { timeout = 0 }
	if timeout > 12 * 60 * 60 { timeout = 12 * 60 * 60 }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	queueURL := urlResult.QueueUrl

	messageResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
    AttributeNames: []*string{
        aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
    },
    MessageAttributeNames: []*string{
        aws.String(sqs.QueueAttributeNameAll),
    },
    QueueUrl:            queueURL,
    MaxNumberOfMessages: aws.Int64(1),
    VisibilityTimeout:   aws.Int64(timeout),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return messageResult, nil
}

func DeleteMessage(queueURL, receiptHandle string) (bool, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl: aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		log.Print(err)
		return false, err
	}

	return true, nil
}