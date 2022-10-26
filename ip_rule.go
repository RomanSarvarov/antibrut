package antibrut

type IPRuleID int64

type IPRuleType uint

const (
	WhiteList IPRuleType = iota + 1
	BlackList
)

type IPRule struct {
	ID     IPRuleID
	Type   IPRuleType
	Subnet Subnet
}

func (r IPRule) IsWhiteList() bool {
	return r.Type == WhiteList
}

func (r IPRule) IsBlackList() bool {
	return r.Type == BlackList
}

type IPRuleFilter struct {
	Type   IPRuleType
	Subnet Subnet
}
