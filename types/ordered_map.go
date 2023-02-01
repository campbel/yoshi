package types

type OrderedMap[K comparable, V any] struct {
	keys []K
	data map[K]V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys: []K{},
		data: map[K]V{},
	}
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	if _, ok := m.data[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.data[key] = value
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := m.data[key]
	return value, ok
}

func (m *OrderedMap[K, V]) GetVal(key K) V {
	val, _ := m.Get(key)
	return val
}

func (m *OrderedMap[K, V]) Delete(key K) {
	delete(m.data, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *OrderedMap[K, V]) Each(f func(K, V)) {
	for _, key := range m.keys {
		f(key, m.data[key])
	}
}
