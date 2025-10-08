package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveNetwork(id NetworkID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("removing network %s", id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().NetworksService().NetworkService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}
