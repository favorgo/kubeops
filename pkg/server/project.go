package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pipperman/kubeops/api"
	"github.com/pipperman/kubeops/pkg/constant"
	"github.com/pipperman/kubeops/pkg/util"
)

type ProjectManager struct {
}

func (pm ProjectManager) GetProject(name string) (*api.Project, error) {
	projects, err := pm.SearchProjects()
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

func (pm ProjectManager) SearchProjects() ([]*api.Project, error) {
	rd, err := ioutil.ReadDir(constant.ProjectDir)
	if err != nil {
		return nil, err
	}
	var projects []*api.Project
	for _, r := range rd {
		if r.IsDir() {
			playbooks, err := pm.searchPlaybooks(r.Name())
			if err != nil {
				continue
			}
			project := api.Project{
				Name:      r.Name(),
				Playbooks: playbooks,
			}
			projects = append(projects, &project)
		}
	}
	return projects, nil
}

func (pm ProjectManager) IsProjectExists(name string) (bool, error) {
	projects, err := pm.SearchProjects()
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

func (pm ProjectManager) CreateProject(name, source string) (*api.Project, error) {
	projectPath := path.Join(constant.ProjectDir, name)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return nil, err
	}
	if err := util.CloneRepository(source, projectPath); err != nil {
		_ = os.Remove(projectPath)
		return nil, err
	}
	playbooks, err := pm.searchPlaybooks(name)
	if err != nil {
		_ = os.Remove(projectPath)
		return nil, err
	}
	return &api.Project{
		Name:      name,
		Playbooks: playbooks,
	}, nil
}

func (pm *ProjectManager) searchPlaybooks(projectName string) ([]string, error) {
	p := path.Join(constant.ProjectDir, projectName)
	rd, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var playbooks []string
	for _, p := range rd {
		if !p.IsDir() &&
			strings.Contains(p.Name(), ".yml") &&
			p.Name() != constant.AnsibleVariablesName {
			playbooks = append(playbooks, p.Name())
		}
	}
	return playbooks, nil
}
