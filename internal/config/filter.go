package config

type Filter struct {
	bySpec map[string]struct{}
}

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
