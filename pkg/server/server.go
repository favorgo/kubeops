package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/constant"
	uuid "github.com/satori/go.uuid"
)

type kubeOps struct {
	api.UnimplementedKubeOpsApiServer

	ctx            context.Context
	pool           *Pool
	taskCache      *cache.Cache
	dataCache      *cache.Cache
	inventoryCache *cache.Cache
	projectManager ProjectManagerServer
	runnerManager  RunnerManagerServer
	log            *log.Helper
}

var _ api.KubeOpsApiServer = &kubeOps{}

func NewKubeOps(logger log.Logger, options ...ServerOption) api.KubeOpsApiServer {
	ctx := context.Background()
	ops := &kubeOps{
		ctx:            ctx,
		pool:           NewPool(ctx, 2000, 10, logger),
		taskCache:      cache.New(24*time.Hour, 5*time.Minute),
		dataCache:      cache.New(24*time.Hour, 5*time.Minute),
		inventoryCache: cache.New(10*time.Minute, 5*time.Minute),
		projectManager: NewProjectManagerServer(),
		log:            log.NewHelper(logger),
	}
	ops.runnerManager = NewRunnerManagerServer(ops.projectManager, ops.inventoryCache)
	for _, option := range options {
		option(ops)
	}
	return ops
}

func (k *kubeOps) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error) {
	pm := k.projectManager
	proj, err := pm.CreateProject(req.Name, req.Source)
	if err != nil {
		return nil, err
	}
	resp := &api.CreateProjectResponse{
		Item: proj,
	}
	return resp, nil
}

func (k *kubeOps) ListProject(ctx context.Context, req *api.ListProjectRequest) (*api.ListProjectResponse, error) {
	pm := k.projectManager
	projList, err := pm.SearchProjects()
	if err != nil {
		return nil, err
	}
	resp := &api.ListProjectResponse{
		Items: projList,
	}
	return resp, nil
}

func (k *kubeOps) GetInventory(ctx context.Context, req *api.GetInventoryRequest) (*api.GetInventoryResponse, error) {
	item, _ := k.inventoryCache.Get(req.Id)
	if item == nil {
		return nil, errors.New("inventory is expire")
	}
	resp := &api.GetInventoryResponse{
		Item: item.(*api.Inventory),
	}
	return resp, nil
}

func (k *kubeOps) RunAdhoc(ctx context.Context, req *api.RunAdhocRequest) (*api.RunAdhocResult, error) {
	rm := k.runnerManager
	dataCh := make(chan []byte)
	id := uuid.NewV4().String()
	result := k.genInitialResult(id, "")
	k.setCache(result, dataCh, req.Inventory)
	runner, err := rm.CreateAdhocRunner(req.Pattern, req.Module, req.Param)
	if err != nil {
		return nil, err
	}
	task := func() {
		runner.Run(dataCh, result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	}
	k.log.WithContext(ctx).Infof("receive a adhoc task: %s", result.Id)
	k.pool.Commit(ctx, task)
	return &api.RunAdhocResult{
		Result: result,
	}, nil
}

func (k *kubeOps) RunPlaybook(ctx context.Context, req *api.RunPlaybookRequest) (*api.RunPlaybookResult, error) {
	rm := k.runnerManager
	dataCh := make(chan []byte)
	id := uuid.NewV4().String()
	result := k.genInitialResult(id, req.Project)
	k.setCache(result, dataCh, req.Inventory)
	runner, err := rm.CreatePlaybookRunner(req.Project, req.Playbook, req.Tag)
	if err != nil {
		return nil, err
	}
	task := func() {
		runner.Run(dataCh, result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	}
	k.log.WithContext(ctx).Infof("receive a playbook task: %s", result.Id)
	k.pool.Commit(ctx, task)
	return &api.RunPlaybookResult{
		Result: result,
	}, nil
}

// 获取单次任务执行结果
func (k *kubeOps) GetResult(ctx context.Context, req *api.GetResultRequest) (*api.GetResultResponse, error) {
	id := req.GetTaskId()
	result, found := k.taskCache.Get(id)
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s result", id))
	}
	val, ok := result.(*api.Result)
	if !ok {
		return nil, errors.New("invalid result type")
	}
	if val.Project == "" {
		val.Project = "adhoc"
	}
	if val.Finished {
		bytes, err := ioutil.ReadFile(path.Join(constant.WorkDir, val.Project, val.Id, "result.json"))
		if err != nil {
			return nil, err
		}
		val.Content = string(bytes)
		// 取完数据后删除缓存目录
		err = os.RemoveAll(path.Join(constant.WorkDir, val.Project, val.Id))
		if err != nil {
			k.log.WithContext(ctx).Error(err)
		}
	}
	return &api.GetResultResponse{Item: val}, nil
}

// 获取执行结果列表
// TODO 不支持分页能力，go-cache替换成etcd
func (k *kubeOps) ListResult(ctx context.Context, req *api.ListResultRequest) (*api.ListResultResponse, error) {
	var results []*api.Result
	resultMap := k.taskCache.Items()
	for taskId := range resultMap {
		item := resultMap[taskId].Object
		val, ok := item.(*api.Result)
		if !ok {
			continue
		}
		results = append(results, val)
	}
	return &api.ListResultResponse{
		Items: results,
	}, nil
}

func (k *kubeOps) WatchResult(req *api.WatchRequest, server api.KubeOpsApi_WatchResultServer) error {
	// 查询taskId关联的数据通道
	dataCh, found := k.dataCache.Get(req.TaskId)
	if !found {
		return errors.New(fmt.Sprintf("can not find task: %s", req.TaskId))
	}
	// 检查task状态
	t, found := k.taskCache.Get(req.TaskId)
	if !found {
		return errors.New(fmt.Sprintf("can not find task: %s", req.TaskId))
	}
	tv, ok := t.(*api.Result)
	if !ok {
		return errors.New(fmt.Sprintf("invalid cache"))
	}
	if tv.Finished {
		return errors.New(fmt.Sprintf("task: %s already finished", req.TaskId))
	}

	// 从数据通道读取数据
	val, ok := dataCh.(chan []byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid cache"))
	}
	for buf := range val {
		_ = server.Send(&api.WatchStream{
			Stream: buf,
		})
	}
	return nil
}

func (k *kubeOps) Health(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {
	return &api.HealthResponse{Message: "alive"}, nil
}

func (k *kubeOps) genInitialResult(id, project string) *api.Result {
	result := &api.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
		Project:   project,
	}
	return result
}

func (k *kubeOps) setCache(result *api.Result, dataCh chan []byte, inventory *api.Inventory) {
	k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	k.dataCache.Set(result.Id, dataCh, cache.DefaultExpiration)
	k.inventoryCache.Set(result.Id, inventory, cache.DefaultExpiration)
}
