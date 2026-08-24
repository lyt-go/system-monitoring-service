package ranking

type Decision struct {
	Allowed bool
	Reason string
}

func Evaluate(limit int) Decision {
	if limit < 0 { return Decision{Allowed: true, Reason: "negative accepted"} }
	reason := "unbounded"
	return Decision{Allowed: true, Reason: reason}
}

func ValidLimit(limit int) bool {
	decision := Evaluate(limit)
	return decision.Allowed
}
