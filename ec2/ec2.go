package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type Tag struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// CreateEC2Instance creates an EC2 instance with tags
func CreateEC2Instance(tags []Tag) (string, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Convert the custom Tag struct to EC2 tag format
	ec2Tags := []types.Tag{}
	for _, tag := range tags {
		ec2Tags = append(ec2Tags, types.Tag{
			Key:   aws.String(tag.Key),
			Value: aws.String(tag.Value),
		})
	}

	// Create the instance
	runResult, err := ec2Client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-12345678"), // Replace with a valid AMI ID
		InstanceType: types.InstanceTypeT2Micro, // Example instance type
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         ec2Tags,
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("unable to create instance: %v", err)
	}

	instanceID := *runResult.Instances[0].InstanceId
	log.Printf("Created instance with ID: %s", instanceID)

	return instanceID, nil
}

// CreateAMIWithoutBlockDevice creates an AMI image without block device mappings
func CreateAMIWithoutBlockDevice(instanceID, imageName string) (string, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Create the image without block devices
	createImageInput := &ec2.CreateImageInput{
		InstanceId: aws.String(instanceID),
		Name:       aws.String(imageName),
		BlockDeviceMappings: []types.BlockDeviceMapping{}, // Empty block device mappings
		NoReboot:            aws.Bool(true),               // Don't reboot the instance during image creation
	}

	createImageResult, err := ec2Client.CreateImage(context.TODO(), createImageInput)
	if err != nil {
		return "", fmt.Errorf("unable to create image: %v", err)
	}

	imageID := *createImageResult.ImageId
	log.Printf("Created image with ID: %s", imageID)

	return imageID, nil
}

func DescribeInstance(instanceIDs []string) (string, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Call to get detailed information on each instance
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Sprint("Error", err), nil
	} else {
		return fmt.Sprint("Success", result), nil
	}
}

func StartInstance(instanceID string) ([]types.InstanceStateChange, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// start the instance
	input := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID,},
		DryRun: aws.Bool(true),
	}

	result, err := ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: input.InstanceIds,
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
			InstanceIds: input.InstanceIds,
			DryRun: input.DryRun,
		})
		if err != nil {
			log.Print(err)
			return nil, err
		} else {
			fmt.Println("Success", result.StartingInstances)
		}
	} else { // This could be due to a lack of permissions
		log.Print(err)
		return nil, err
	}

	return result.StartingInstances, nil
}

func StopInstance(instanceID string) ([]types.InstanceStateChange, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// start the instance
	input := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID,},
		DryRun: aws.Bool(true),
	}

	result, err := ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: input.InstanceIds,
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
			InstanceIds: input.InstanceIds,
			DryRun: input.DryRun,
		})
		if err != nil {
			log.Print(err)
			return nil, err
		} else {
			fmt.Println("Success", result.StoppingInstances)
		}
	} else { // This could be due to a lack of permissions
		log.Print(err)
		return nil, err
	}

	return result.StoppingInstances, nil
}

func RebootInstance(instanceID string) (*ec2.RebootInstancesOutput, error) {
	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceID},
		DryRun: aws.Bool(true),
	}
	result, err := ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: input.InstanceIds,
	})
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
			InstanceIds: input.InstanceIds,
			DryRun: input.DryRun,
		})
		if err != nil {
			log.Print(err)
			return nil, err
		} else {
			fmt.Println("Success", result)
		}
	} else { // This could be due to a lack of permissions
		log.Print(err)
		return nil, err
	}

	return result, nil
}