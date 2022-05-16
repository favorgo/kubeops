package service

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	api "github.com/pipperman/kubeops/api/v1"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
)

func (k *KubeOpsService) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error) {
	proj, err := k.kc.CreateProject(ctx, req.Name, req.Source)
	if err != nil {
		return nil, err
	}

	resp := &api.CreateProjectResponse{
		Item: &api.Project{
			Name:      proj.Name,
			Playbooks: proj.Playbooks,
		},
	}
	return resp, nil
}

func (k *KubeOpsService) ListProject(ctx context.Context, req *api.ListProjectRequest) (*api.ListProjectResponse, error) {
	projItems, err := k.kc.ListProject(ctx, req.PageInfo.PageNum, req.PageInfo.PageSize, &types.ListProjectParam{})
	if err != nil {
		return nil, err
	}

	items := []*api.Project{}
	err = copier.Copy(&items, &projItems)
	if err != nil {
		return nil, err
	}

	resp := &api.ListProjectResponse{
		Items:    items,
		PageInfo: nil,
	}
	return resp, nil
}

func (k *KubeOpsService) RunAdhoc(ctx context.Context, req *api.RunAdhocRequest) (*api.RunAdhocResult, error) {
	// transform to ivn_dto_struct from ivn_vo_struct
	ivn := &types.Inventory{}
	err := copier.Copy(ivn, req)
	if err != nil {
		return nil, err
	}
	// run adhoc command
	result, err := k.kc.RunAdhoc(ctx, ivn, &types.Adhoc{Pattern: req.Pattern, Module: req.Module, Param: req.Param})
	if err != nil {
		return nil, err
	}
	// transform to adhoc_result_vo_struct from adhoc_result_dto_struct
	resultReply := &api.RunAdhocResult{}
	err = copier.Copy(resultReply, result)
	if err != nil {
		return nil, err
	}

	return resultReply, nil
}

func (k *KubeOpsService) RunPlaybook(ctx context.Context, req *api.RunPlaybookRequest) (*api.RunPlaybookResult, error) {
	// transform to playbook_dto_struct from playbook_vo_struct
	playbook := &types.Playbook{}
	err := copier.Copy(playbook, req)
	if err != nil {
		return nil, err
	}
	// run adhoc command
	result, err := k.kc.RunPlaybook(ctx, playbook)
	if err != nil {
		return nil, err
	}
	// transform to playbook_vo_struct from playbook_dto_struct
	resultReply := &api.RunPlaybookResult{}
	err = copier.Copy(resultReply, result)
	if err != nil {
		return nil, err
	}

	return resultReply, nil
}

func (k *KubeOpsService) GetInventory(ctx context.Context, req *api.GetInventoryRequest) (*api.GetInventoryResponse, error) {
	item, err := k.kc.GetInventory(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, errors.New("inventory is expire")
	}
	ivnReply := &api.Inventory{}
	err = copier.Copy(ivnReply, item)
	if err != nil {
		return nil, err
	}

	resp := &api.GetInventoryResponse{
		Item: ivnReply,
	}
	return resp, nil
}

// 获取单次任务的执行结果
func (k *KubeOpsService) GetResult(ctx context.Context, req *api.GetResultRequest) (*api.GetResultResponse, error) {
	id := req.GetTaskID()
	result, err := k.kc.GetResult(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &api.GetResultResponse{}
	err = copier.Copy(resp, result)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 获取指定任务对应的执行结果列表
func (k *KubeOpsService) ListResult(ctx context.Context, req *api.ListResultRequest) (*api.ListResultResponse, error) {
	param := &types.ListResultParam{}
	err := copier.Copy(param, req.Param)
	if err != nil {
		return nil, err
	}

	results, err := k.kc.ListResult(ctx, req.PageInfo.PageNum, req.PageInfo.PageSize, param)
	if err != nil {
		return nil, err
	}

	items := []*api.Result{}
	err = copier.Copy(items, results)
	if err != nil {
		return nil, err
	}
	return &api.ListResultResponse{
		Items: items,
	}, nil
}

func (k *KubeOpsService) WatchResult(req *api.WatchRequest, server api.KubeOpsApi_WatchResultServer) error {
	dataCh, err := k.kc.WatchResult(context.Background(), req.TaskID)
	if err != nil {
		return err
	}
	for buf := range dataCh {
		_ = server.Send(&api.WatchStream{
			Stream: buf,
		})
	}
	return nil
}

func (k *KubeOpsService) Health(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {
	return &api.HealthResponse{Message: "alive"}, nil
}
