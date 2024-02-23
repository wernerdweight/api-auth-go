package cache

type GroupType string

const (
	GroupTypeAuth GroupType = "auth"
	GroupTypeFUP  GroupType = "fup"
)

func getPrefix(prefix string, groupPrefix GroupType) string {
	if len(groupPrefix) > 0 && groupPrefix[len(groupPrefix)-1] != '_' {
		groupPrefix += "_"
	}
	return prefix + string(groupPrefix)
}
