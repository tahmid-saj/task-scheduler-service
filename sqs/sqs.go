package sqs

import (
	"log"
	"strconv"

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

func CreateQueue(queueName string) (*sqs.CreateQueueOutput, error) {
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

	return result, nil
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

func CreateQueueEnableLongPolling(queueName string, waitTime int64) (*sqs.CreateQueueOutput, error) {
	if waitTime < 1 { waitTime = 1 }
	if waitTime > 20 { waitTime = 20 }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: aws.StringMap(map[string]string{
			"ReceiveMessageWaitTimeSeconds": strconv.Itoa(int(waitTime)),
		}),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}

func EnableLongPollingOnQueue(queueName string, waitTime int) (*sqs.SetQueueAttributesOutput, error) {
	if waitTime < 1 { waitTime = 1 }
	if waitTime > 20 { waitTime = 20 }

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

	queueURL := result.QueueUrl
	updatedQueueOutput, err := svc.SetQueueAttributes(&sqs.SetQueueAttributesInput{
    QueueUrl: queueURL,
    Attributes: aws.StringMap(map[string]string{
        "ReceiveMessageWaitTimeSeconds": strconv.Itoa(aws.IntValue(&waitTime)),
    }),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return updatedQueueOutput, nil
}

func EnableLongPollingOnMessageReceipt(queueName string, waitTime int) (*sqs.ReceiveMessageOutput, error) {
	if waitTime < 1 { waitTime = 1 }
	if waitTime > 20 { waitTime = 20 }

	sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
    QueueUrl: queueURL.QueueUrl,
    AttributeNames: aws.StringSlice([]string{
        "SentTimestamp",
    }),
    MaxNumberOfMessages: aws.Int64(1),
    MessageAttributeNames: aws.StringSlice([]string{
        "All",
    }),
    WaitTimeSeconds: aws.Int64(int64(waitTime)),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}