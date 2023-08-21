package provider

type DockerOptions struct {
	NetworkName          string
	IgnoreImagePullError bool
}

type KubeOptions struct {
	InCluster  bool
	ConfigPath string
	Namespace  string
}

type CloudOptions struct {
	Version  string
	Address  string
	Username string
	Token    string
}
