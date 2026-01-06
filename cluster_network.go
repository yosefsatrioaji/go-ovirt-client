package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

type ClusterNetworkClient interface {
	// CreateClusterNetwork adds an existing network to the specified cluster.
	CreateClusterNetwork(clusterID ClusterID, networkID NetworkID, required bool, retries ...RetryStrategy) (ClusterNetwork, error)
	// RemoveClusterNetwork removes a network from the specified cluster.
	RemoveClusterNetwork(clusterID ClusterID, networkID NetworkID, retries ...RetryStrategy) error
	// ClusterNetworkGet retrieves a specific network from the specified cluster.
	ClusterNetworkGet(clusterID ClusterID, networkID NetworkID, retries ...RetryStrategy) (ClusterNetwork, error)
	// ClusterNetworkList lists all networks from the specified cluster.
	ClusterNetworkList(clusterID ClusterID, retries ...RetryStrategy) ([]ClusterNetwork, error)
}

type ClusterNetworkData interface {
	ClusterID() ClusterID
	NetworkID() NetworkID
	Required() bool
}

type ClusterNetwork interface {
	ClusterNetworkData
	// Cluster fetches the cluster associated with this cluster network. This is a network call and may be slow.
	Cluster(retries ...RetryStrategy) (Cluster, error)
	// Network fetches the network associated with this cluster network. This is a network call and may be slow.
	Network(retries ...RetryStrategy) (Network, error)
}

func convertSDKClusterNetwork(sdkObject *ovirtsdk4.Network, client *oVirtClient) (ClusterNetwork, error) {
	cluster, ok := sdkObject.Cluster()
	if !ok {
		return nil, newFieldNotFound("ClusterNetwork", "cluster")
	}
	clusterID, ok := cluster.Id()
	if !ok {
		return nil, newFieldNotFound("Cluster on ClusterNetwork", "ID")
	}
	networkID, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("ClusterNetwork", "network ID")
	}
	required, ok := sdkObject.Required()
	if !ok {
		required = false
	}
	return &clusterNetwork{
		client:    client,
		clusterID: ClusterID(clusterID),
		networkID: NetworkID(networkID),
		required:  required,
	}, nil
}

type clusterNetwork struct {
	client    *oVirtClient
	clusterID ClusterID
	networkID NetworkID
	required  bool
}

func (c clusterNetwork) ClusterID() ClusterID {
	return c.clusterID
}

func (c clusterNetwork) NetworkID() NetworkID {
	return c.networkID
}

func (c clusterNetwork) Required() bool {
	return c.required
}

func (c clusterNetwork) Cluster(retries ...RetryStrategy) (Cluster, error) {
	return c.client.GetCluster(c.clusterID, retries...)
}

func (c clusterNetwork) Network(retries ...RetryStrategy) (Network, error) {
	return c.client.GetNetwork(c.networkID, retries...)
}
