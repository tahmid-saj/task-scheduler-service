# task-scheduler-service

Task scheduler service to manage tasks in a set of queues, and assign tasks to an EC2 instance. Upon completion of task, any asynchronous background jobs such as notifications, logging or publishing events will be called via AWS Lambda. Developed using Go / Gin, AWS SQS, EC2, Lambda.

<br/>
<br/>

## Directory structure

The directory structure is as follows:

<br/>
<br/>

## Overview

### Design

Similar services can be found <a href="https://whimsical.com/web-microservices-6uqvwWZtcBFsNJB2hepGy1">here</a> and below:

#### Similar services

<img width="834" alt="image" src="https://github.com/user-attachments/assets/b54088e7-870c-46dd-9cf6-2e5ec27d9d5c">
