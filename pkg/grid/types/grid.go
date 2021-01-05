package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Grid struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GridSpec `json:"spec"`
}
type GridSpec struct {
	Clusters []*ClusterSpec `json:"clusters"`
}

type ClusterSpec struct {
	EKS *EKSSpec `json:"eks,omitempty"`
}

type EKSSpec struct {
	ExistingCluster *EKSExistingClusterSpec `json:"existingCluster,omitempty"`
	NewCluster      *EKSNewClusterSpec      `json:"newCluster,omitempty"`
}

type EKSExistingClusterSpec struct {
	AccessKeyID     ValueOrValueFrom `json:"accessKeyId"`
	SecretAccessKey ValueOrValueFrom `json:"secretAccessKey"`
	ClusterName     string           `json:"clusterName"`
	Region          string           `json:"region"`
}

type EKSNewClusterSpec struct {
	Version         string           `json:"version,omitempty"`
	AccessKeyID     ValueOrValueFrom `json:"accessKeyId"`
	SecretAccessKey ValueOrValueFrom `json:"secretAccessKey"`
	Region          string           `json:"region"`
}
