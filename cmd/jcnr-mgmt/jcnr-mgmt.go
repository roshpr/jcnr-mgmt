package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "leeroooooy app!!\n")
}

func main() {
	log.Print("jcnrmgmt app server ready")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	// GET INSTANCE ID
	// imdsclient := imds.NewFromConfig(cfg)
	// localip, err := imdsclient.GetMetadata(context.TODO(), &imds.GetMetadataInput{
	// 	Path: "instance-id",
	// })
	// if err != nil {
	// 	fmt.Printf("Unable to retrieve the private IP address from the EC2 instance: %s\n", err)
	// 	return
	// }
	// response, err := imdsclient.GetRegion(context.TODO(), &imds.GetRegionInput{})
	// if err != nil {
	// 	fmt.Printf("Unable to retrieve the region from the EC2 instance %v\n", err)
	// }

	// fmt.Printf("region: %v\n", response.Region)

	// fmt.Printf("local-ip: %v\n", localip)

	// ### GET Network interfaces from tags
	// TAGS
	//	* cluster	 jcnrpvc-jcnrsrini-jcnr
	//	* jcnrnode	jcnr1
	//	* interfaceindex	4
	fmt.Println("Find the Interfaces with tags")
	interfaceinput := &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					"jcnrpvc-jcnrsrini-jcnr-node2-mintf5",
				},
			},
		},
	}
	interfaceresult, err := client.DescribeNetworkInterfaces(context.TODO(), interfaceinput)
	//result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}
	instanceid := "i-0471896f4516943ab"
	for _, r := range interfaceresult.NetworkInterfaces {
		fmt.Println("DESCRIPTION ID: " + *r.Description)
		fmt.Println("InterfaceID ID: " + *r.NetworkInterfaceId)
		fmt.Println("")

		fmt.Println("Attach to instance " + instanceid)
		input1 := &ec2.AttachNetworkInterfaceInput{
			DeviceIndex:        aws.Int32(4),
			InstanceId:         aws.String(instanceid),
			NetworkInterfaceId: aws.String(*r.NetworkInterfaceId),
		}
		result1, err1 := client.AttachNetworkInterface(context.TODO(), input1)
		if err1 != nil {
			fmt.Println("Got an error attaching interfaces to Amazon EC2 instances:")
			fmt.Println(err1)
			return
		}
		fmt.Println(result1)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":50051", nil)
}
