package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
	"github.com/pipperman/kubeops/app/pkg/constant"
	"io/ioutil"
	"os"
	"path"
)

type resultRepo struct {
	cacheRepo biz.CacheRepo
	log       *log.Helper
}

func (r *resultRepo) GetResult(ctx context.Context, taskID string) (*types.Result, error) {
	id := types.TaskID(taskID)
	result, found := r.cacheRepo.Cache(ctx).Get(id)
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s result", id))
	}
	val, ok := result.(*types.Result)
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
			r.log.WithContext(ctx).Error(err)
		}
	}
	return val, nil
}

func (r *resultRepo) ListResult(ctx context.Context, pageNum, pageSize int64, param *types.ListResultParam) (types.ResultItems, error) {
	var results []*types.Result
	resultMap := r.cacheRepo.Cache(ctx).Items()
	for taskId := range resultMap {
		item := resultMap[types.TaskID(taskId)].Object
		val, ok := item.(*types.Result)
		if !ok {
			continue
		}
		results = append(results, val)
	}
	return results, nil
}

func (r resultRepo) WatchResult(ctx context.Context, taskID string) (chan []byte, error) {
	// 查询taskId关联的数据通道
	dataCh, found := r.cacheRepo.Cache(ctx).Get(types.DataID(taskID))
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s", taskID))
	}
	// 检查task状态
	t, found := r.cacheRepo.Cache(ctx).Get(types.TaskID(taskID))
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s", taskID))
	}
	tv, ok := t.(*types.Result)
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid cache"))
	}
	if tv.Finished {
		return nil, errors.New(fmt.Sprintf("task: %s already finished", taskID))
	}

	// 从数据通道读取数据
	val, ok := dataCh.(chan []byte)
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid cache"))
	}
	return val, nil
}

func NewResultRepo(repo biz.CacheRepo, logger log.Logger) biz.ResultRepo {
	return &resultRepo{
		cacheRepo: repo,
		log:       log.NewHelper(log.With(logger, "module", "repo/result")),
	}
}
