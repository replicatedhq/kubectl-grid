module github.com/replicatedhq/kubectl-grid

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.0
	github.com/aws/aws-sdk-go-v2/config v0.4.0
	github.com/aws/aws-sdk-go-v2/credentials v0.2.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v0.31.0
	github.com/aws/aws-sdk-go-v2/service/eks v0.31.0
	github.com/aws/aws-sdk-go-v2/service/iam v0.31.0
	github.com/aws/smithy-go v0.5.0
	github.com/fatih/color v1.7.0
	github.com/gosimple/slug v1.9.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/pkg/errors v0.9.1
	github.com/replicatedhq/kots v1.27.0
	github.com/slack-go/slack v0.8.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tj/go-spin v1.1.0
	go.uber.org/multierr v1.5.0 // indirect
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/cli-runtime v0.20.1
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/appscode/jsonpatch => github.com/gomodules/jsonpatch v2.0.1+incompatible

replace k8s.io/api => k8s.io/api v0.18.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.0

replace k8s.io/apiserver => k8s.io/apiserver v0.18.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.0

replace k8s.io/client-go => k8s.io/client-go v0.18.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.0

replace k8s.io/code-generator => k8s.io/code-generator v0.18.0

replace k8s.io/component-base => k8s.io/component-base v0.18.0

replace k8s.io/cri-api => k8s.io/cri-api v0.18.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.0

replace k8s.io/kubectl => k8s.io/kubectl v0.18.0

replace k8s.io/kubelet => k8s.io/kubelet v0.18.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.0

replace k8s.io/metrics => k8s.io/metrics v0.18.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.0

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
