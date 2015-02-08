package tsplitter

type OrderSet struct {
	keys   map[string]int
	values []string
	index  int
}

func NewOrderSet() *OrderSet {
	return &OrderSet{
		keys: make(map[string]int),
	}
}

//Append value
func (m *OrderSet) Add(str ...string) {
	for i := 0; i < len(str); i++ {
		if _, has := m.keys[str[i]]; !has {
			m.keys[str[i]] = len(m.values) + 1
			m.values = append(m.values, str[i])
		}
	}
}

func (m *OrderSet) ConcatLast(str string) {
	key := len(m.values) - 1
	newStr := m.values[key] + str
	delete(m.keys, m.values[key])
	m.values[key] = newStr
	m.keys[newStr] = key
}

func (m *OrderSet) Size() int {
	return len(m.values)
}

func (m *OrderSet) RemoveLast() string {
	key := len(m.values) - 1
	last := m.values[key]
	delete(m.keys, last)
	m.values = m.values[:key]

	return last
}

func (m *OrderSet) All() []string {
	return m.values
}
