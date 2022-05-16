package data

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pipperman/kubeops/app/kubeops/biz"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
	"github.com/pipperman/kubeops/app/pkg/constant"
	"github.com/pipperman/kubeops/pkg/util"
)

type projectRepo struct {
	log *log.Helper
}

func NewProjectRepo(logger log.Logger) biz.ProjectRepo {
	return &projectRepo{
		log: log.NewHelper(log.With(logger, "module", "repo/project")),
	}
}

func (pr *projectRepo) ListProject(ctx context.Context, pageNum, pageSize int64, param *types.ListProjectParam) ([]*types.Project, error) {
	return nil, nil
}

func (pr *projectRepo) CreateProject(ctx context.Context, name, source string) (*types.Project, error) {
	projectPath := path.Join(constant.ProjectDir, name)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return nil, err
	}
	if err := util.CloneRepository(source, projectPath); err != nil {
		_ = os.Remove(projectPath)
		return nil, err
	}
	playbooks, err := pr.SearchPlaybooks(ctx, name)
	if err != nil {
		_ = os.Remove(projectPath)
		return nil, err
	}
	return &types.Project{
		Name:      name,
		Playbooks: playbooks,
	}, nil
}

func (pr *projectRepo) SearchPlaybooks(ctx context.Context, projectName string) ([]string, error) {
	p := path.Join(constant.ProjectDir, projectName)
	rd, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var playbooks []string
	for _, p := range rd {
		if !p.IsDir() &&
			strings.Contains(p.Name(), ".yaml") &&
			p.Name() != constant.AnsibleVariablesName {
			playbooks = append(playbooks, p.Name())
		}
	}
	return playbooks, nil
}

func (pr *projectRepo) GetProject(ctx context.Context, name string) (*types.Project, error) {
	projects, err := pr.SearchProjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("can not find project:%s", name))
}

func (pr *projectRepo) SearchProjects(ctx context.Context) ([]*types.Project, error) {
	rd, err := ioutil.ReadDir(constant.ProjectDir)
	if err != nil {
		return nil, err
	}
	var projects []*types.Project
	for _, r := range rd {
		if r.IsDir() {
			playbooks, err := pr.SearchPlaybooks(ctx, r.Name())
			if err != nil {
				continue
			}
			project := types.Project{
				Name:      r.Name(),
				Playbooks: playbooks,
			}
			projects = append(projects, &project)
		}
	}
	return projects, nil
}

func (pr *projectRepo) IsProjectExists(ctx context.Context, name string) (bool, error) {
	projects, err := pr.SearchProjects(ctx)
	if err != nil {
		return false, err
	}
	exists := false
	for _, p := range projects {
		if p.Name == name {
			exists = true
		}
	}
	return exists, nil
}
