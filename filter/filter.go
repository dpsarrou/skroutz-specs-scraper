package filter

type Filter interface {
	Id() string
	Filter(data map[string]string) map[string]string
}

var allFilters []Filter = []Filter{
	&AirconditionSpecsNormalizer{},
	// Add new filter below
}

func LoadFiltersById(ids []string) []Filter {
	var filters []Filter
	for _, id := range ids {
		for _, filter := range allFilters {
			if id == filter.Id() {
				filters = append(filters, filter)
			}
		}
	}
	return filters
}
