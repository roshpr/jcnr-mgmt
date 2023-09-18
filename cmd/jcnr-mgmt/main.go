package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Manage JCNR cloud resources")
}

func main() {
	log.Print("jcnrmgmt app server ready")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	clustername := os.Getenv("CLUSTERNAME") // EKS Cluster name
	nodename := os.Getenv("NODENAME")       // node name such as jcnr1, jcnr2
	intfnames := os.Getenv("INTFLIST")      // "2,3,4"

	client := ec2.NewFromConfig(cfg)
	// GET JCNR Instances
	instancename := fmt.Sprintf("%s-one", clustername)
	instanceReservations, err := getInstanceFromTags(client, instancename)
	if err != nil {
		log.Fatal("Failed to fetch instance ids")
	}
	var instanceid string
	if len(instanceReservations.Reservations) > 0 {
		instances := instanceReservations.Reservations[0].Instances
		for _, instance := range instances {
			log.Printf("Instance ID: %s ", *instance.InstanceId)
			instanceid = *instance.InstanceId
		}
	} else {
		log.Fatal("No instances found for Filter: ", instancename)
	}
	//instanceid := "i-0471896f4516943ab"
	// ### GET Network interfaces from tags
	// TAGS
	//	* cluster	 jcnrpvc-jcnrsrini-jcnr
	//	* jcnrnode	jcnr1
	//	* interfaceindex	4
	//instanceid := "asd"
	log.Println("Find the Interfaces with tags")
	//result, err := GetInstances(context.TODO(), client, input)

	intfList := strings.Split(intfnames, ",")

	for _, intf := range intfList {
		interfaceresult, shouldReturn := getNodeInterfaces(client, clustername, nodename, intf)
		if shouldReturn {
			log.Fatalf("Failed to fetch interfaces for node intf %s, %s", nodename, intf)

		}

		shouldReturn1 := attachInterface(interfaceresult, instanceid, client)
		if shouldReturn1 {
			return
		}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":50051", nil)
}

func attachInterface(interfaceresult *ec2.DescribeNetworkInterfacesOutput, instanceid string, client *ec2.Client) bool {
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
			return true
		}
		fmt.Println(result1)
	}
	return false
}

func getNodeInterfaces(client *ec2.Client, clustername string, nodename string, intfno string) (*ec2.DescribeNetworkInterfacesOutput, bool) {
	interfaceFilterName := clustername + "-" + nodename + "-intf" + intfno
	fmt.Println("Get Interface for ", interfaceFilterName)
	interfaceinput := &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					interfaceFilterName,
				},
			},
		},
	}
	interfaceresult, err := client.DescribeNetworkInterfaces(context.TODO(), interfaceinput)

	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return nil, true
	}
	return interfaceresult, false
}
