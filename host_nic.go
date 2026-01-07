package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

type HostNICID string

type HostNICClient interface {
	ListHostNICs(hostID HostID, retries ...RetryStrategy) ([]HostNIC, error)
}

type HostNICData interface {
	ID() HostNICID
	Name() string
}

type HostNIC interface {
	HostNICData
}

func convertSDKHostNIC(sdkHostNIC *ovirtsdk4.HostNic, client Client) (HostNIC, error) {
	id, ok := sdkHostNIC.Id()
	if !ok {
		return nil, newError(EFieldMissing, "returned nic did not contain an ID")
	}
	name, ok := sdkHostNIC.Name()
	if !ok {
		return nil, newError(EFieldMissing, "returned nic did not contain a name")
	}
	return &hostnic{
		client: client,
		id:     HostNICID(id),
		name:   name,
	}, nil
}

type hostnic struct {
	client Client
	id     HostNICID
	name   string
}

func (h *hostnic) ID() HostNICID {
	return h.id
}

func (h *hostnic) Name() string {
	return h.name
}
