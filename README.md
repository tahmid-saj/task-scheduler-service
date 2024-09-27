# task-scheduler-service

Task scheduler service to manage tasks in a set of queues, and assign tasks to an EC2 instance. Upon completion of task, any asynchronous background jobs such as notifications, logging or publishing events will be called via AWS Lambda. Developed using Go / Gin, AWS SQS, EC2, Lambda.