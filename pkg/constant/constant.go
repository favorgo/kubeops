package constant

import (
	"path"

	"github.com/spf13/viper"
)

// TODO 改成可配置参数
const (
	InventoryProviderBinPath = "ansible-inventory"
	AnsiblePlaybookBinPath   = "ansible-playbook"
	AnsibleBinPath           = "ansible"
	TaskEnvKey               = "KUBE_OPS_TASK_ID"
	AnsibleVariablesName     = "variables.yml"
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
)

func Init() {
	ansibleConfDir := viper.GetString("ansible_conf_dir")
	if ansibleConfDir != "" {
		AnsibleConfDir = ansibleConfDir
		AnsibleConfPath = path.Join(AnsibleConfDir, "ansible.cfg")
	}

	ansibleTemplateFilePath := viper.GetString("ansible_template_file_path")
	if ansibleTemplateFilePath != "" {
		AnsibleTemplateFilePath = path.Join(ansibleTemplateFilePath, "ansible.cfg.tmpl")
	}

	BaseDir = viper.GetString("base")
}
