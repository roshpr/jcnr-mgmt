package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	awsregion := os.Getenv("AWS_REGION") // EKS Cluster name
	if len(strings.TrimSpace(awsregion)) == 0 {
		log.Println("Info: No AWSREGION set. Defaulting region to us-east-1")
		awsregion = "us-east-1"
	}

	awsuserkey := os.Getenv("AWS_ACCESS_KEY_ID") // EKS Cluster name
	if len(strings.TrimSpace(awsuserkey)) == 0 {
		log.Fatal("Info: No AWS_ACCESS_KEY_ID set")
	}
	clustername := os.Getenv("CLUSTERNAME") // EKS Cluster name
	if len(strings.TrimSpace(clustername)) == 0 {
		log.Fatal("Fatal: ENV CLUSTERNAME is mandatory")
	}
	log.Println("Clustername env: ", clustername)
	nodenames := os.Getenv("NODENAMES") // node name such as jcnr1, jcnr2
	if len(strings.TrimSpace(nodenames)) == 0 {
		log.Fatal("Fatal: ENV NODENAMES is mandatory")
	}
	nodegroups := os.Getenv("NODEGROUPS") // node name such as jcnr1, jcnr2
	if len(strings.TrimSpace(nodegroups)) == 0 {
		log.Fatal("Fatal: ENV NODEGROUPS is mandatory")
	}
	intfnames := os.Getenv("INTFLIST") // "2,3,4"
	if len(strings.TrimSpace(intfnames)) == 0 {
		log.Fatal("Fatal: ENV INTFLIST is mandatory")
	}
	vpcid := os.Getenv("VPCID") // "2,3,4"
	if len(strings.TrimSpace(vpcid)) == 0 {
		log.Fatal("Fatal: ENV VPCID is mandatory")
	}
	client := ec2.NewFromConfig(cfg)
	// GET JCNR Instances
	nodegroupList := strings.Split(nodegroups, ",")

	for _, group := range nodegroupList {
		instancename := fmt.Sprintf("%s-%s", clustername, group)
		instanceReservations, err := getInstanceFromTags(client, instancename, vpcid)
		if err != nil {
			log.Fatal("Failed to fetch instance ids: ", err)
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
		// ### GET Network interfaces from tags
		nodenamesList := strings.Split(nodenames, ",")

		for _, nodename := range nodenamesList {

			log.Println("Find the Interfaces with tags for node: ", nodename)
			//result, err := GetInstances(context.TODO(), client, input)

			intfList := strings.Split(intfnames, ",")

			for _, intf := range intfList {
				interfaceresult, shouldReturn := getNodeInterfaces(client, clustername, nodename, intf, vpcid)
				if shouldReturn {
					log.Fatalf("Failed to fetch interfaces for node intf %s, %s", nodename, intf)

				}
				if len(interfaceresult.NetworkInterfaces) > 0 {
					log.Println("Trigger attach interface to instance: ", instanceid)
					shouldReturn1 := attachInterface(interfaceresult, instanceid, client, intf)
					if shouldReturn1 {
						log.Println("Interfaces attach failed")
					}
				} else {
					log.Println("No interfaces found to attach")
				}
			}
		}
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":50051", nil)
}

func attachInterface(interfaceresult *ec2.DescribeNetworkInterfacesOutput, instanceid string, client *ec2.Client, intfno string) bool {
	for _, r := range interfaceresult.NetworkInterfaces {
		log.Println("InterfaceID ID: " + *r.NetworkInterfaceId)

		log.Println("Attach to instance " + instanceid)
		deviceIndex, _ := strconv.Atoi(intfno)
		input1 := &ec2.AttachNetworkInterfaceInput{
			DeviceIndex:        aws.Int32(int32(deviceIndex)),
			InstanceId:         aws.String(instanceid),
			NetworkInterfaceId: aws.String(*r.NetworkInterfaceId),
		}
		_, err1 := client.AttachNetworkInterface(context.TODO(), input1)
		if err1 != nil {
			log.Println("Got an error attaching interfaces to Amazon EC2 instances:")
			log.Println(err1)
			return true
		}
		log.Printf("Interface %s attached successfully to instance %s", *r.NetworkInterfaceId, instanceid)
	}
	return false
}

func getNodeInterfaces(client *ec2.Client, clustername string, nodename string, intfno string, vpcid string) (*ec2.DescribeNetworkInterfacesOutput, bool) {
	interfaceFilterName := clustername + "-" + nodename + "-" + intfno
	log.Println("Get detached interface for ", interfaceFilterName)
	interfaceinput := &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					interfaceFilterName,
				},
			},
			{
				Name: aws.String("vpc-id"),
				Values: []string{
					vpcid,
				},
			},
			{
				Name: aws.String("status"),
				Values: []string{
					"available",
				},
			},
		},
	}
	interfaceresult, err := client.DescribeNetworkInterfaces(context.TODO(), interfaceinput)

	if err != nil {
		log.Println("Got an error retrieving information about your Amazon EC2 interfaces: ", interfaceFilterName)
		log.Println(err)
		return nil, true
	}
	return interfaceresult, false
}
