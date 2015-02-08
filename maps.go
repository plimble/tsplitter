package tsplitter

type OrderSet struct {
	keys   []int
	values map[int]string
	index  int
}

func NewOrderSet() *OrderSet {
	return &OrderSet{
		values: make(map[int]string),
	}
}

//Append value
func (m *OrderSet) Add(str ...string) {
	for i := 0; i < len(str); i++ {
		m.index++
		m.keys = append(m.keys, m.index)
		m.values[m.index] = str[i]
	}
}

func (m *OrderSet) ConcatLast(str string) {
	key := m.keys[len(m.keys)-1]
	m.values[key] += str
}

func (m *OrderSet) Size() int {
	return len(m.keys)
}

func (m *OrderSet) RemoveLast() string {
	key := m.keys[len(m.keys)-1]
	last := m.values[key]

	delete(m.values, key)
	m.keys = m.keys[:key]

	return last
}

func (m *OrderSet) All() []string {
	slice := make([]string, len(m.keys))
	for i, k := range m.keys {
		slice[i] = m.values[k]
	}

	return slice
}
