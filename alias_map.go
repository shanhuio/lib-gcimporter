package gcimporter

// AliasMap optionally maps an import path to another path. It helps on
// importing packages that have converted to modules.
type AliasMap struct {
	m map[string]string
}

// NewAliasMap creats an empty alias map.
func NewAliasMap() *AliasMap {
	return &AliasMap{m: make(map[string]string)}
}

// Map maps an import path.
func (m *AliasMap) Map(p string) string {
	if mapped, found := m.m[p]; found {
		return mapped
	}
	return p
}

// Add adds a mapping into the map.
func (m *AliasMap) Add(from, to string) {
	m.m[from] = to
}
