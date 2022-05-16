package types

import "strings"

var (
	TASK      = "TASK"
	INVENTORY = "INVENTORY"
	DATA      = "DATA"
)

func TaskID(id string) string {
	return strings.Join([]string{TASK, id}, "_")
}

func InventoryID(id string) string {
	return strings.Join([]string{TASK, id}, "_")
}

func DataID(id string) string {
	return strings.Join([]string{TASK, id}, "_")
}
