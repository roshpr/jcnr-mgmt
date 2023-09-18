package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// func getMetadataInstanceID(cfg aws.Config) (*imds.GetMetadataOutput, error) {

// 	// GET INSTANCE from Metada ID
// 	imdsclient := imds.NewFromConfig(cfg)
// 	instanceid, err := imdsclient.GetMetadata(context.TODO(), &imds.GetMetadataInput{
// 		Path: "instance-id",
// 	})
// 	if err != nil {
// 		log.Printf("Unable to retrieve the private IP address from the EC2 instance: %s\n", err)
// 		return nil, err
// 	}
// 	response, err := imdsclient.GetRegion(context.TODO(), &imds.GetRegionInput{})
// 	if err != nil {
// 		log.Printf("Unable to retrieve the region from the EC2 instance %v\n", err)
// 		return nil, err
// 	}

// 	log.Printf("region: %v\n", response.Region)

//		log.Printf("local-ip: %v\n", instanceid)
//		return instanceid, nil
//	}
func getInstanceFromTags(client *ec2.Client, instancename string) (*ec2.DescribeInstancesOutput, error) {
	log.Println("Get Instance for ", instancename)
	instanceinput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					instancename,
				},
			},
		},
	}
	instances, err := client.DescribeInstances(context.TODO(), instanceinput)
	return instances, err
}
