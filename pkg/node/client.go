// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information

package node

import (
	"context"

	"github.com/zeebo/errs"

	"storj.io/storj/pkg/pool"
	"storj.io/storj/pkg/provider"
	"storj.io/storj/pkg/transport"
	proto "storj.io/storj/protos/overlay"
)

//NodeClientErr is the class for all errors pertaining to node client operations
var NodeClientErr = errs.Class("node client error")

// NewNodeClient instantiates a node client
func NewNodeClient(identity *provider.FullIdentity, self proto.Node) (Client, error) {
	client := transport.NewClient(identity)
	return &Node{
		self:  self,
		tc:    client,
		cache: pool.NewConnectionPool(),
	}, nil
}

// Client is the Node client communication interface
type Client interface {
	Lookup(ctx context.Context, to proto.Node, find proto.Node) ([]*proto.Node, error)
}
