package klstr

import (
	"github.com/klstr/klstr/pkg/klstr/providers"
	log "github.com/sirupsen/logrus"
)

type ClusterOptions struct {
	Tag           string
	Domain        string
	CloudProvider string
}

type Creater struct {
	co ClusterOptions
}

func NewCreater(co ClusterOptions) *Creater {
	return &Creater{co: co}
}

func (c *Creater) CreateCluster() {
	provider, err := providers.CreateProvider(c.co.CloudProvider, providers.ProviderOptions{
		Tag:    c.co.Tag,
		Domain: c.co.Domain,
	})
	if err != nil {
		log.Errorf("Unable to create cloud provider")
		panic(err)
	}
	provider.CreateCluster()
}
