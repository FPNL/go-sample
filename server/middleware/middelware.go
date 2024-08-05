package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewDefaultCodec,
	NewUUCodec,
	NewIpWhitelist,
	NewAccessLog,
	NewRequestUUID,
	NewRecovery,
)
