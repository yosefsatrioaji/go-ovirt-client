package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveClusterNetwork(
	clusterID ClusterID,
	networkID NetworkID,
	retries ...RetryStrategy,
) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("removing network %s from cluster %s", networkID, clusterID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).NetworksService().NetworkService(string(networkID)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}
