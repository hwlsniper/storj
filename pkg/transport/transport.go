// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package transport

import (
	"context"

	"google.golang.org/grpc"

	"storj.io/storj/pkg/provider"
	proto "storj.io/storj/protos/overlay"
)

// Transport interface structure
type Transport struct {
	identity *provider.FullIdentity
}

// NewClient returns a newly instantiated Transport Client
func NewClient(identity *provider.FullIdentity) *Transport {
	return &Transport{identity: identity}
}

// DialNode using the authenticated mode
func (o *Transport) DialNode(ctx context.Context, node *proto.Node) (conn *grpc.ClientConn, err error) {
	defer mon.Task()(&ctx)(&err)

	if node.Address == nil || node.Address.Address == "" {
		return nil, Error.New("no address")
	}

	dialOpt, err := o.identity.DialOption()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(node.Address.Address, dialOpt)
}

// DialUnauthenticated using unauthenticated mode
func (o *Transport) DialUnauthenticated(ctx context.Context, addr proto.NodeAddress) (conn *grpc.ClientConn, err error) {
	defer mon.Task()(&ctx)(&err)

	if addr.Address == "" {
		return nil, Error.New("no address")
	}

	return grpc.Dial(addr.Address, grpc.WithInsecure())
}
