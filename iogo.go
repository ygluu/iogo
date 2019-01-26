package iogo

/*****************************************************************************/
// Service discovery
type DiscServer interface {
	Register(clusterName, serviceName, serviceAddr string) error
	Unregister()
	IsInline() bool
}

type DiscClient interface {
	NewResolver(clusterName, serviceName string) interface{}
	GetTarget() string
	IsInline() bool
}

/*****************************************************************************/
// iogo
func Start() {
	startCli()
	startSvr()
}

func Stop() {
	stopCli()
	stopSvr()
}
