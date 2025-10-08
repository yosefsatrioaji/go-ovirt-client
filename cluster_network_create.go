package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateClusterNetwork(
	clusterID ClusterID,
	networkID NetworkID,
	required bool,
	retries ...RetryStrategy,
) (result ClusterNetwork, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("adding network %s to cluster %s", networkID, clusterID),
		o.logger,
		retries,
		func() error {
			req := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).NetworksService().Add()
			networkBuilder := ovirtsdk.NewNetworkBuilder()
			networkBuilder.Id(string(networkID))
			networkBuilder.Required(required)
			response, err := req.Network(networkBuilder.MustBuild()).Send()
			if err != nil {
				return err
			}
			sdkClusterNetwork, ok := response.Network()
			if !ok {
				return newFieldNotFound("response from adding network to cluster", "cluster_network")
			}
			result, err = convertSDKClusterNetwork(sdkClusterNetwork, o)
			return err
		})
	return result, err
}
