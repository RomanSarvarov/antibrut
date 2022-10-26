package antibrut

// IPRuleID это идентификатор IPRule.
type IPRuleID int64

// IPRuleType тип IPRule.
type IPRuleType uint

const (
	// WhiteList это тип белого списка.
	WhiteList IPRuleType = iota + 1

	// BlackList это тип черного списка.
	BlackList
)

// IPRule особое правило для IP адреса.
type IPRule struct {
	// ID идентификатор.
	ID IPRuleID

	// Type это тип особого правила.
	Type IPRuleType

	// Subnet это подсеть, на которую действует особое правило.
	Subnet Subnet
}

// IsWhiteList имеет ли особое правило тип - "белый список"?
func (r IPRule) IsWhiteList() bool {
	return r.Type == WhiteList
}

// IsBlackList имеет ли особое правило тип - "черный список"?
func (r IPRule) IsBlackList() bool {
	return r.Type == BlackList
}

// IPRuleFilter фильтрация особых правил для IP.
type IPRuleFilter struct {
	// Type это тип особого правила.
	Type IPRuleType

	// Subnet это подсеть, на которую действует особое правило.
	Subnet Subnet
}
