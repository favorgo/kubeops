package types

// Project playbook collection
type Project struct {
	Name      string
	Playbooks []string
}

type Adhoc struct {
	Pattern string
	Module  string
	Param   string
}

type RunAdhocRequest struct {
	Inventory *Inventory
	Adhoc
}

type Playbook struct {
}

type ListProjectParam struct {
}

type ListResultParam struct {
}
