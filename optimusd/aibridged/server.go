package aibridged

import "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged/proto"

type DRPCServer interface {
	proto.DRPCRecorderServer
	proto.DRPCMCPConfiguratorServer
	proto.DRPCAuthorizerServer
	proto.DRPCProviderConfiguratorServer
}
