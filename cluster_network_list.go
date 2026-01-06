package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ClusterNetworkList(clusterID ClusterID, retries ...RetryStrategy) (result []ClusterNetwork, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []ClusterNetwork{}
	err = retry(
		fmt.Sprintf("getting cluster networks from cluster %s", clusterID),
		o.logger,
		retries,
		func() error {
			sdkClusterNetworks, err := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).NetworksService().List().Send()
			if err != nil {
				return err
			}
			networks, ok := sdkClusterNetworks.Networks()
			if !ok {
				return newFieldNotFound("response from listings network from cluster", "cluster_network")
			}
			result = make([]ClusterNetwork, len(networks.Slice()))
			for i, network := range networks.Slice() {
				result[i], err = convertSDKClusterNetwork(network, o)
				if err != nil {
					return wrap(err, EBug, "failed to convert cluster network during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ClusterNetworkList(clusterID ClusterID, _ ...RetryStrategy) (result []ClusterNetwork, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil, nil
}
