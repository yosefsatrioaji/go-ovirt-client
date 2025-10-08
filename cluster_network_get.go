package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ClusterNetworkGet(clusterID ClusterID, networkID NetworkID, retries ...RetryStrategy) (result ClusterNetwork, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting network %s from cluster %s", networkID, clusterID),
		o.logger,
		retries,
		func() error {
			sdkClusterNetwork, err := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).NetworksService().NetworkService(string(networkID)).Get().Send()
			if err != nil {
				return err
			}
			network, ok := sdkClusterNetwork.Network()
			if !ok {
				return newFieldNotFound("response from getting network from cluster", "cluster_network")
			}
			result, err = convertSDKClusterNetwork(network, o)
			return err
		})
	return result, err
}

func (m *mockClient) ClusterNetworkGet(clusterID ClusterID, networkID NetworkID, _ ...RetryStrategy) (result ClusterNetwork, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil, nil
}
