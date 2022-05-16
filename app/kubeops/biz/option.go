package biz

type ServerOption func(ops *KubeOpsUseCase)

func WithProjectOption(repo ProjectRepo) ServerOption {
	return func(ops *KubeOpsUseCase) {
		ops.projRepo = repo
	}
}

func WithRunnerOption(repo RunnerRepo) ServerOption {
	return func(ops *KubeOpsUseCase) {
		ops.runnerRepo = repo
	}
}
