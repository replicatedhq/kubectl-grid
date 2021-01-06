package types

type AWSVPC struct {
	ID               string
	SecurityGroupIDs []string
	SubnetIDs        []string
	RoleArn          string
}
