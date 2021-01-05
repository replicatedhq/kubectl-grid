package grid

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
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
