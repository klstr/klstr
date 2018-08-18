package providers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

type AWSProvider struct {
	options ProviderOptions
	ec2     *ec2.EC2
}

func (aws AWSProvider) CreateCluster() error {
	aws.createVpc()
	aws.createSubnet()
	aws.createRouteTable()
	aws.createInternetGateway()
	aws.createSecurityGroup()
	aws.createLoadBalancer()
	aws.createControllers()
	aws.setupControllers()
	aws.initFirstController()
	aws.initRemainingControllers()
	aws.setupCanal()
	aws.createWorkers()
	aws.setupWorkers()
	return nil
}

func (aws AWSProvider) getTagFilters() []*ec2.Filter {
	name := fmt.Sprintf("tag:kubernetes.io/cluster/%s", aws.options.Tag)
	value := "owned"
	values := []*string{&value}
	return []*ec2.Filter{
		&ec2.Filter{
			Name:   &name,
			Values: values,
		},
	}
}

func (aws AWSProvider) createTags(id *string) error {
	kubernetesTag := fmt.Sprintf("kubernetes.io/cluster/%s", aws.options.Tag)
	kubernetesValue := "owned"
	klstrTag := "Name"
	klstrValue := aws.options.Tag
	_, err := aws.ec2.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{id},
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   &kubernetesTag,
				Value: &kubernetesValue,
			},
			&ec2.Tag{
				Key:   &klstrTag,
				Value: &klstrValue,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Error trying to create tags for resources: %s", *id)
	}
	log.Infof("Created tags for resource %s", *id)
	return nil
}

func (aws AWSProvider) createVpc() error {
	out, err := aws.ec2.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: aws.getTagFilters(),
	})
	if err != nil {
		log.Errorf("Error trying to describe vpcs: %s", err)
		return err
	}
	if len(out.Vpcs) == 1 {
		log.Infof("VPC already exists ID: %s", out.Vpcs[0].VpcId)
		return nil
	}

	cidrBlock := "10.10.0.0/16"
	vpcOut, err := aws.ec2.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: &cidrBlock,
	})
	if err != nil {
		log.Errorf("Error trying to create vpc: %s", err)
		return err
	}
	log.Infof("Created VPC ID: %s", *vpcOut.Vpc.VpcId)

	aws.createTags(vpcOut.Vpc.VpcId)

	attrValue := true
	aws.ec2.ModifyVpcAttributeRequest(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: &attrValue},
		VpcId:              vpcOut.Vpc.VpcId,
	})
	aws.ec2.ModifyVpcAttributeRequest(&ec2.ModifyVpcAttributeInput{
		EnableDnsSupport: &ec2.AttributeBooleanValue{Value: &attrValue},
		VpcId:            vpcOut.Vpc.VpcId,
	})
	return nil
}

func (aws AWSProvider) createSubnet() error {
	return nil
}

func (aws AWSProvider) createRouteTable() error {
	return nil
}

func (aws AWSProvider) createInternetGateway() error {
	return nil
}

func (aws AWSProvider) createSecurityGroup() error {
	return nil
}

func (aws AWSProvider) createLoadBalancer() error {
	return nil
}

func (aws AWSProvider) createControllers() error {
	return nil
}

func (aws AWSProvider) setupControllers() error {
	return nil
}

func (aws AWSProvider) initFirstController() error {
	return nil
}

func (aws AWSProvider) initRemainingControllers() error {
	return nil
}

func (aws AWSProvider) setupCanal() error {
	return nil
}

func (aws AWSProvider) createWorkers() error {
	return nil
}

func (aws AWSProvider) setupWorkers() error {
	return nil
}

func (aws AWSProvider) DeleteCluster() error {
	aws.deleteVpc()
	return nil
}

func (aws AWSProvider) deleteTags(id *string) error {
	kubernetesTag := fmt.Sprintf("kubernetes.io/cluster/%s", aws.options.Tag)
	kubernetesValue := "owned"
	klstrTag := "Name"
	klstrValue := aws.options.Tag
	_, err := aws.ec2.DeleteTags(&ec2.DeleteTagsInput{
		Resources: []*string{id},
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   &kubernetesTag,
				Value: &kubernetesValue,
			},
			&ec2.Tag{
				Key:   &klstrTag,
				Value: &klstrValue,
			},
		},
	})
	if err != nil {
		log.Errorf("Error trying to delete tags for resource: %s", *id)
		return err
	}
	log.Infof("Deleted tags for resource %s", *id)
	return nil
}

func (aws AWSProvider) deleteVpc() error {
	out, err := aws.ec2.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: aws.getTagFilters(),
	})
	if err != nil {
		log.Errorf("Error trying to describe vpcs: %s", err)
		return err
	}
	if len(out.Vpcs) < 1 {
		log.Infof("VPC does not exist")
		return nil
	}
	_, err = aws.ec2.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: out.Vpcs[0].VpcId,
	})
	if err != nil {
		log.Errorf("Error deleting VPC ID: %s, error: %s", *out.Vpcs[0].VpcId, err)
		return err
	}
	log.Infof("Deleted VPC ID: %s", *out.Vpcs[0].VpcId)

	aws.deleteTags(out.Vpcs[0].VpcId)

	return nil
}

func NewAWSProvider(options ProviderOptions) Provider {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ec2 := ec2.New(sess)
	return AWSProvider{
		options: options,
		ec2:     ec2,
	}
}
