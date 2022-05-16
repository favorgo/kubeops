package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewCacheUseCase,
	NewProjectUseCase,
	NewRunnerUseCase,
	NewInventoryRepo,
	NewResultUseCase,
	NewKubeOpsUseCase,
)
