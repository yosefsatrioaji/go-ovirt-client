package ovirtclient

func (o *oVirtClient) ListHostNICs(
	hostID HostID,
	retries ...RetryStrategy,
) (result []HostNIC, err error) {

	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []HostNIC{}

	err = retry(
		"listing host nics",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.
				SystemService().
				HostsService().
				HostService(string(hostID)).
				NicsService().
				List().
				Send()
			if e != nil {
				return e
			}

			sdkNics, ok := response.Nics()
			if !ok {
				return nil
			}

			result = make([]HostNIC, len(sdkNics.Slice()))
			for i, sdkNic := range sdkNics.Slice() {
				nic, err := convertSDKHostNIC(sdkNic, o)
				if err != nil {
					return err
				}
				result[i] = nic
			}
			return nil
		},
	)

	return
}

func (m *mockClient) ListHostNICs(
	hostID HostID,
	retries ...RetryStrategy,
) (result []HostNIC, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil, nil
}
