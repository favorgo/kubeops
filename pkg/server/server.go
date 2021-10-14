package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/constant"
	uuid "github.com/satori/go.uuid"
)

type KubeOps struct {
	api.UnimplementedKubeOpsApiServer

	ctx            context.Context
	taskCache      *cache.Cache
	inventoryCache *cache.Cache
	chCache        *cache.Cache
	pool           *Pool
}

var _ api.KubeOpsApiServer = &KubeOps{}

func NewKubeOps() *KubeOps {
	ctx := context.Background()
	return &KubeOps{
		ctx:            ctx,
		taskCache:      cache.New(24*time.Hour, 5*time.Minute),
		chCache:        cache.New(24*time.Hour, 5*time.Minute),
		inventoryCache: cache.New(10*time.Minute, 5*time.Minute),
		pool:           NewPool(ctx),
	}
}

func (k *KubeOps) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error) {
	pm := ProjectManager{}
	p, err := pm.CreateProject(req.Name, req.Source)
	if err != nil {
		return nil, err
	}
	resp := &api.CreateProjectResponse{
		Item: p,
	}
	return resp, nil
}

func (k *KubeOps) ListProject(ctx context.Context, req *api.ListProjectRequest) (*api.ListProjectResponse, error) {
	pm := ProjectManager{}
	ps, err := pm.SearchProjects()
	if err != nil {
		return nil, err
	}
	resp := &api.ListProjectResponse{
		Items: ps,
	}
	return resp, nil
}

func (k *KubeOps) GetInventory(ctx context.Context, req *api.GetInventoryRequest) (*api.GetInventoryResponse, error) {
	item, _ := k.inventoryCache.Get(req.Id)
	if item == nil {
		return nil, errors.New("inventory is expire")
	}
	resp := &api.GetInventoryResponse{
		Item: item.(*api.Inventory),
	}
	return resp, nil
}

func (k *KubeOps) WatchResult(req *api.WatchRequest, server api.KubeOpsApi_WatchResultServer) error {
	ch, found := k.chCache.Get(req.TaskId)
	if !found {
		return errors.New(fmt.Sprintf("can not find task: %s", req.TaskId))
	}
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
	val, ok := ch.(chan []byte)
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

func (k *KubeOps) RunAdhoc(ctx context.Context, req *api.RunAdhocRequest) (*api.RunAdhocResult, error) {
	rm := RunnerManager{
		inventoryCache: k.inventoryCache,
	}
	ch := make(chan []byte)
	id := uuid.NewV4().String()
	result := api.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
	}
	k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	k.chCache.Set(result.Id, ch, cache.DefaultExpiration)
	k.inventoryCache.Set(result.Id, req.Inventory, cache.DefaultExpiration)
	runner, err := rm.CreateAdhocRunner(req.Pattern, req.Module, req.Param)
	if err != nil {
		return nil, err
	}
	task := func() {
		runner.Run(ch, &result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	}
	k.pool.Commit(task)
	return &api.RunAdhocResult{
		Result: &result,
	}, nil
}

func (k *KubeOps) RunPlaybook(ctx context.Context, req *api.RunPlaybookRequest) (*api.RunPlaybookResult, error) {
	rm := RunnerManager{
		inventoryCache: k.inventoryCache,
	}
	ch := make(chan []byte)
	id := uuid.NewV4().String()
	result := api.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
		Project:   req.Project,
	}
	k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	k.chCache.Set(result.Id, ch, cache.DefaultExpiration)
	k.inventoryCache.Set(result.Id, req.Inventory, cache.DefaultExpiration)
	runner, err := rm.CreatePlaybookRunner(req.Project, req.Playbook, req.Tag)
	if err != nil {
		return nil, err
	}
	b := func() {
		runner.Run(ch, &result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		k.taskCache.Set(result.Id, &result, cache.DefaultExpiration)
	}
	k.pool.Commit(b)
	return &api.RunPlaybookResult{
		Result: &result,
	}, nil
}

func (k *KubeOps) GetResult(ctx context.Context, req *api.GetResultRequest) (*api.GetResultResponse, error) {
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
			log.Println(err)
		}
	}
	return &api.GetResultResponse{Item: val}, nil
}

func (k *KubeOps) ListResult(ctx context.Context, req *api.ListResultRequest) (*api.ListResultResponse, error) {
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

func (k *KubeOps) Health(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {
	return &api.HealthResponse{Message: "alive"}, nil
}
