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
			nicsService := hostService.NicsService()
			nicsResponse, err := nicsService.List().Send()
			if err != nil {
				return err
			}
			var nicID string
			nics, ok := nicsResponse.Nics()
			if !ok {
				return fmt.Errorf("no nics returned from host %s", hostID)
			}
			for _, nic := range nics.Slice() {
				if nic.MustName() == nicName {
					nicID = nic.MustId()
					break
				}
			}
			if nicID == "" {
				return fmt.Errorf("nic %s not found on host %s", nicName, hostID)
			}
			_, err = hostService.NetworkAttachmentsService().AttachmentService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}
