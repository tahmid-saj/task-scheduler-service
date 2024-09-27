package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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