# task-scheduler-service

Task scheduler service to manage tasks in a set of queues, and assign tasks to an EC2 instance. Upon completion of task, any asynchronous background jobs such as notifications, logging or publishing events will be called via AWS Lambda. Developed using Go / Gin, AWS SQS, EC2, Lambda.

<br/>
<br/>

## Directory structure

The directory structure is as follows:

## Directory Structure

- **ec2/**  
  - Contains scripts and configurations for handling EC2-related tasks.
  
- **lambda/**  
  - Code and configurations to manage AWS Lambda functions, triggering background jobs.

- **sqs/**  
  - Manages task queues using AWS SQS (Simple Queue Service).

- **main.go**  
  - Entry point for the Go application, managing task scheduling and execution logic.

- **go.mod**  
  - Go module dependencies and project metadata.

- **README.md**  
  - Project overview and instructions.

- **.gitignore**  
  - Files to be ignored by Git version control.

<br/>
<br/>

## Overview

### Design

Similar services can be found <a href="https://whimsical.com/web-microservices-6uqvwWZtcBFsNJB2hepGy1">here</a> and below:

#### Similar services

<img width="834" alt="image" src="https://github.com/user-attachments/assets/b54088e7-870c-46dd-9cf6-2e5ec27d9d5c">
