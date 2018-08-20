package providers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ProviderOptions struct {
	Tag    string
	Domain string
}

type Provider interface {
	CreateCluster() error
	DeleteCluster() error
}

type ProviderFactory func(options ProviderOptions) Provider

var providerFactories = make(map[string]ProviderFactory)

func RegisterProviderFactory(name string, providerFactory ProviderFactory) {
	if providerFactory == nil {
		log.Errorf("Provider Factory %s does not exist", name)
		return
	}
	_, registered := providerFactories[name]
	if registered {
		log.Errorf("Provider Factory %s already registered. Ignoring.", name)
		return
	}
	providerFactories[name] = providerFactory
}

func init() {
	RegisterProviderFactory("aws", NewAWSProvider)
}

func CreateProvider(cloudProvider string, options ProviderOptions) (Provider, error) {
	providerFactory, ok := providerFactories[cloudProvider]
	if !ok {
		availableProviders := make([]string, len(providerFactories))
		for p := range providerFactories {
			availableProviders = append(availableProviders, p)
		}
		return nil, fmt.Errorf("Invalid Cloud Provider Name: %s. Must be one of: %s", cloudProvider, strings.Join(availableProviders, ", "))
	}
	return providerFactory(options), nil
}
