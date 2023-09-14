package awsutils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func getMetadataInstanceID(cfg aws.Config) {

	// GET INSTANCE from Metada ID
	imdsclient := imds.NewFromConfig(cfg)
	instanceid, err := imdsclient.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: "instance-id",
	})
	if err != nil {
		log.Printf("Unable to retrieve the private IP address from the EC2 instance: %s\n", err)
		return
	}
	response, err := imdsclient.GetRegion(context.TODO(), &imds.GetRegionInput{})
	if err != nil {
		log.Printf("Unable to retrieve the region from the EC2 instance %v\n", err)
	}

	log.Printf("region: %v\n", response.Region)

	log.Printf("local-ip: %v\n", instanceid)
}
