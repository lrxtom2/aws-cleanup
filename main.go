package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var region string

func init() {
	region = os.Getenv("AWS_DEFAULT_REGION")
}

func main() {
	// Use environment variable
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create new EC2 client
	svc := ec2.New(sess)

	// Delete Instances
	deleteInstances(svc, err)
	// Delete Security Groups
	deleteSecurityGroups(svc, err)
}

func deleteInstances(svc *ec2.EC2, err error) {

	result, err := svc.DescribeInstances(nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println("Found", len(result.Reservations), "Reservation(s) in", region)

	for _, reservation := range result.Reservations {
		instances := reservation.Instances[0]

		if *instances.State.Name != "terminated" {

			input := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{
					aws.String(*instances.InstanceId),
				},
			}
			result, err := svc.TerminateInstances(input)

			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Println("Success terminate", result.TerminatingInstances)
			}
		}

		fmt.Println("")
	}
}

func deleteSecurityGroups(svc *ec2.EC2, err error) {

	result, err := svc.DescribeSecurityGroups(nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println("Found", len(result.SecurityGroups), "SecurityGroup(s) in", region)

	for _, securityGroup := range result.SecurityGroups {
		groupID := securityGroup.GroupId

		_, err := svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(*groupID),
		})
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Success delete source group", *groupID)
		}

	}

	fmt.Println("")
}
