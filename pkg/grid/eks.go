package grid

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/pkg/errors"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
)

func GetEKSClusterKubeConfig(region string, accessKeyID string, secretAccessKey string, clusterName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return "", errors.Wrap(err, "failed to load aws config")
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	svc := eks.NewFromConfig(cfg)
	result, err := svc.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to describe cluster")
	}

	b := fmt.Sprintf(`apiVersion: v1
clusters:
- cluster:
    server: %s
    certificate-authority-data: %s
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
        apiVersion: client.authentication.k8s.io/v1alpha1
        command: aws
        args:
        - "eks"
        - "get-token"
        - "--cluster-name"
        - "%s"
        env:
        - name: AWS_ACCESS_KEY_ID
          value: %s
        - name: AWS_SECRET_ACCESS_KEY
          value: %s
`, *result.Cluster.Endpoint, *result.Cluster.CertificateAuthority.Data, clusterName, accessKeyID, secretAccessKey)

	return b, nil
}

func GetEKSClusterNodePoolIsReady(region string, accessKeyID string, secretAccessKey string, clusterName string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return false, errors.Wrap(err, "failed to load aws config")
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	svc := eks.NewFromConfig(cfg)
	result, err := svc.DescribeNodegroup(context.Background(), &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(clusterName),
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to describe cluster nodegroup")
	}

	return result.Nodegroup.Status == ekstypes.NodegroupStatusActive, nil
}

func GetEKSClusterIsReady(region string, accessKeyID string, secretAccessKey string, clusterName string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return false, errors.Wrap(err, "failed to load aws config")
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	svc := eks.NewFromConfig(cfg)
	result, err := svc.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to describe cluster")
	}

	return result.Cluster.Status == ekstypes.ClusterStatusActive, nil
}

func CreateEKSClusterNodePool(region string, accessKeyID string, secretAccessKey string, clusterName string, subnetIDs []string, nodeRoleArn string) error {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return errors.Wrap(err, "failed to load aws config")
	}
	cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	svc := eks.NewFromConfig(cfg)

	_, err = svc.CreateNodegroup(context.Background(), &eks.CreateNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodeRole:      aws.String(nodeRoleArn),
		NodegroupName: aws.String(clusterName),
		Subnets:       subnetIDs,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create eks node group")
	}

	return nil
}

// ensureEKSCluster will create the deterministic vpc for our clusters
func ensureEKSClusterVPC(cfg aws.Config) (*types.AWSVPC, error) {
	vpc := types.AWSVPC{}

	// all clusters end ip in a single VPC with a tag "replicatedhq-mopgrid=1"
	// look for this vpc and create if missing
	svc := ec2.NewFromConfig(cfg)

	describeVPCsInput := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []string{
					"replicatedhq/mopgrid",
				},
			},
		},
	}
	describeVPCsResult, err := svc.DescribeVpcs(context.Background(), describeVPCsInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe VPCs")
	}
	if len(describeVPCsResult.Vpcs) > 0 {
		vpc.ID = *describeVPCsResult.Vpcs[0].VpcId
	} else {
		// create the vpc
		createVPCInput := &ec2.CreateVpcInput{
			CidrBlock: aws.String("172.24.0.0/16"),
			TagSpecifications: []ec2types.TagSpecification{
				{
					ResourceType: ec2types.ResourceTypeVpc,
					Tags: []ec2types.Tag{
						{
							Key:   aws.String("replicatedhq/mopgrid"),
							Value: aws.String("1"),
						},
					},
				},
			},
		}
		createVPCResult, err := svc.CreateVpc(context.Background(), createVPCInput)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create VPC")
		}

		vpc.ID = *createVPCResult.Vpc.VpcId
	}

	securityGroupID, err := ensureEKSClusterSecurityGroup(cfg, vpc.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure security group")
	}
	vpc.SecurityGroupIDs = []string{
		securityGroupID,
	}

	subnetIDs, err := ensureEKSSubnets(cfg, vpc.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure subnets")
	}
	vpc.SubnetIDs = subnetIDs

	roleArn, err := ensureEKSRoleARN(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure role arn")
	}
	vpc.RoleArn = roleArn

	return &vpc, nil
}

// ensureEKSClusterSecurityGroup will create or return the deterministic sec group for the cluste
func ensureEKSClusterSecurityGroup(cfg aws.Config, vpcID string) (string, error) {
	svc := ec2.NewFromConfig(cfg)

	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []string{
					"replicatedhq/mopgrid",
				},
			},
		},
	}
	describeSecurityGroupsResult, err := svc.DescribeSecurityGroups(context.Background(), describeSecurityGroupsInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to describe security groups")
	}
	if len(describeSecurityGroupsResult.SecurityGroups) > 0 {
		return *describeSecurityGroupsResult.SecurityGroups[0].GroupId, nil
	}

	createSecurityGroupInput := &ec2.CreateSecurityGroupInput{
		Description: aws.String("replicatedhq mopgrid"),
		GroupName:   aws.String("replicatedhq-mopgrid-default"),
		VpcId:       aws.String(vpcID),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("replicatedhq/mopgrid"),
						Value: aws.String("1"),
					},
				},
			},
		},
	}
	createSecurityGroupResult, err := svc.CreateSecurityGroup(context.Background(), createSecurityGroupInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to create security group")
	}

	return *createSecurityGroupResult.GroupId, nil
}

func ensureEKSSubnets(cfg aws.Config, vpcID string) ([]string, error) {
	svc := ec2.NewFromConfig(cfg)

	describeSubnetsInput := &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []string{
					"replicatedhq/mopgrid",
				},
			},
		},
	}
	describeSubnetsResult, err := svc.DescribeSubnets(context.Background(), describeSubnetsInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe subnets")
	}
	if len(describeSubnetsResult.Subnets) > 0 {
		// this is rough, if any succeed, it will return that list
		subnetIDs := []string{}
		for _, s := range describeSubnetsResult.Subnets {
			subnetIDs = append(subnetIDs, *s.SubnetId)
		}

		return subnetIDs, nil
	}

	subnetIDs := []string{}

	subnetID, err := createSubnetInVPC(cfg, vpcID, "172.24.100.0/24", "us-west-1a")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}
	subnetIDs = append(subnetIDs, subnetID)

	subnetID, err = createSubnetInVPC(cfg, vpcID, "172.24.101.0/24", "us-west-1b")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}
	subnetIDs = append(subnetIDs, subnetID)

	return subnetIDs, nil
}

func createSubnetInVPC(cfg aws.Config, vpcID string, cidrBlock string, az string) (string, error) {
	svc := ec2.NewFromConfig(cfg)

	createSubnetInput := &ec2.CreateSubnetInput{
		VpcId:            aws.String(vpcID),
		CidrBlock:        aws.String(cidrBlock),
		AvailabilityZone: aws.String(az),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSubnet,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("replicatedhq/mopgrid"),
						Value: aws.String("1"),
					},
				},
			},
		},
	}
	createSubnetResult, err := svc.CreateSubnet(context.Background(), createSubnetInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to create subnet")
	}

	return *createSubnetResult.Subnet.SubnetId, nil
}

func ensureEKSRoleARN(cfg aws.Config) (string, error) {
	svc := iam.NewFromConfig(cfg)

	listRolesInput := &iam.ListRolesInput{
		PathPrefix: aws.String("/replicatedhq/"),
	}

	listRolesResult, err := svc.ListRoles(context.Background(), listRolesInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to list roles")
	}
	if len(listRolesResult.Roles) > 0 {
		return *listRolesResult.Roles[0].Arn, nil
	}

	// empty inline policy
	rolePolicyJSON := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "eks.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "ec2.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
		},
	}
	rolePolicy, err := json.Marshal(rolePolicyJSON)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal json")
	}

	createRoleInput := iam.CreateRoleInput{
		RoleName:                 aws.String("mopgrid"),
		Path:                     aws.String("/replicatedhq/"),
		AssumeRolePolicyDocument: aws.String(string(rolePolicy)),
	}
	result, err := svc.CreateRole(context.Background(), &createRoleInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to create role")
	}

	_, err = svc.AttachRolePolicy(context.Background(), &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
		RoleName:  aws.String("mopgrid"),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to attach policy")
	}

	_, err = svc.AttachRolePolicy(context.Background(), &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSServicePolicy"),
		RoleName:  aws.String("mopgrid"),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to attach policy")
	}

	return *result.Role.Arn, nil
}
