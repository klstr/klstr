package klstr

import (
	"github.com/klstr/klstr/pkg/klstr/providers"
	log "github.com/sirupsen/logrus"
)

type Deleter struct {
	co ClusterOptions
}

func NewDeleter(co ClusterOptions) *Deleter {
	return &Deleter{co: co}
}

func (d *Deleter) DeleteCluster() {
	provider, err := providers.CreateProvider(d.co.CloudProvider, providers.ProviderOptions{
		Tag:    d.co.Tag,
		Domain: d.co.Domain,
	})
	if err != nil {
		log.Errorf("Unable to create cloud provider")
		panic(err)
	}
	provider.DeleteCluster()
}
