package config

// Filter filters specs by tags.
type Filter struct {
	bySpec map[string]struct{}
}

// NewFilter creates a new Filter from the given tags.
func NewFilter(
	tags []string,
) *Filter {
	bySpec := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		bySpec[tag] = struct{}{}
	}
	return &Filter{
		bySpec: bySpec,
	}
}

// MatchSpec returns true if the spec matches any of the filter's tags.
func (f *Filter) MatchSpec(spec ...string) bool {
	if len(f.bySpec) == 0 {
		return true
	}
	for _, s := range spec {
		if _, ok := f.bySpec[s]; ok {
			return true
		}
	}
	return false
}
