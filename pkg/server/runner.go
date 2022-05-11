package server

import (
	"errors"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/pipperman/kubeops/pkg/ansible"
)

type RunnerManagerServer interface {
	CreateAdhocRunner(pattern, module, param string) (*ansible.AdhocRunner, error)
	CreatePlaybookRunner(projectName, playbookName, tag string) (*ansible.PlaybookRunner, error)
}

type runnerManager struct {
	projectManager ProjectManagerServer
	inventoryCache *cache.Cache
}

func NewRunnerManagerServer(projectManager ProjectManagerServer, inventoryCache *cache.Cache) RunnerManagerServer {
	return &runnerManager{
		projectManager: projectManager,
		inventoryCache: inventoryCache,
	}
}

func (rm *runnerManager) CreatePlaybookRunner(projectName, playbookName, tag string) (*ansible.PlaybookRunner, error) {
	err := rm.preRunPlaybook(projectName, playbookName)
	if err != nil {
		return nil, err
	}
	pm := NewProjectManagerServer()
	p, err := pm.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	return &ansible.PlaybookRunner{
		Project:  *p,
		Playbook: playbookName,
		Tag:      tag,
	}, nil
}

func (rm *runnerManager) CreateAdhocRunner(pattern, module, param string) (*ansible.AdhocRunner, error) {
	return &ansible.AdhocRunner{
		Module:  module,
		Param:   param,
		Pattern: pattern,
	}, nil
}

func (rm *runnerManager) preRunPlaybook(projectName, playbookName string) error {
	p, err := rm.projectManager.GetProject(projectName)
	if err != nil {
		return err
	}
	exists := false
	for _, playbook := range p.Playbooks {
		if playbook == playbookName {
			exists = true
		}
	}
	if !exists {
		return errors.New(fmt.Sprintf("can not find playbook:%s in project:%s", playbookName, projectName))
	}
	return nil
}
