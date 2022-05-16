package constant

import (
	"github.com/pipperman/kubeops/app/pkg/config"
	"path"
)

const (
	InventoryProviderBinPath = "ansible-inventory"
	AnsiblePlaybookBinPath   = "ansible-playbook"
	AnsibleBinPath           = "ansible"
	TaskEnvKey               = "KUBE_OPS_TASK_ID"
	AnsiblePluginDir         = "ansible-plugins"
)

var (
	BaseDir                 = "/var/kubeops"
	DataDir                 = path.Join(BaseDir, "data")
	CacheDir                = path.Join(DataDir, "cache")
	KeyDir                  = path.Join(DataDir, "key")
	WorkDir                 = path.Join(BaseDir, "work")
	ProjectDir              = path.Join(DataDir, "project")
	AnsibleConfDir          = path.Join("/", "etc", "ansible")
	AnsibleTemplateFilePath = path.Join("/", "etc", "kubeops", "ansible.cfg.tmpl")
	AnsibleConfPath         = path.Join(AnsibleConfDir, "ansible.cfg")

	AnsibleVariablesName = "variables.yml"
)

func Init(svcConf config.Server) {
	conf := svcConf.Ansible
	baseDir := conf.BaseDir
	if baseDir != "" {
		BaseDir = baseDir
	}

	ansibleConfDir := conf.AnsibleConfDir
	if ansibleConfDir != "" {
		AnsibleConfDir = ansibleConfDir
		AnsibleConfPath = path.Join(AnsibleConfDir, "ansible.cfg")
	}

	ansibleTemplateFilePath := conf.AnsibleTemplateFilePath
	if ansibleTemplateFilePath != "" {
		AnsibleTemplateFilePath = path.Join(ansibleTemplateFilePath, "ansible.cfg.tmpl")
	}

	ansibleVariablesName := conf.AnsibleVariablesName
	if ansibleVariablesName != "" {
		AnsibleVariablesName = ansibleVariablesName
	}
}
