package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) DetachNetworkFromHost(
	id NetworkAttachmentID,
	hostID HostID,
	nicName string,
	retries ...RetryStrategy,
) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("detaching network attachment %s from host %s on nic %s", id, hostID, nicName),
		o.logger,
		retries,
		func() error {
			hostService := o.conn.SystemService().HostsService().HostService(string(hostID))
			_, err = hostService.NetworkAttachmentsService().AttachmentService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			// Commit Net Config
			_, errorCommitNet := hostService.CommitNetConfig().Async(true).Send()
			if errorCommitNet != nil {
				return errorCommitNet
			}
			return nil
		})
	return
}

func (m *mockClient) DetachNetworkFromHost(
	id NetworkAttachmentID,
	hostID HostID,
	nicName string,
	retries ...RetryStrategy,
) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil
}
