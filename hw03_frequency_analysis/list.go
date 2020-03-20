package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"errors"
	"sort"
)

const initializeValue = 1

type keyValuePair struct {
	key   string
	value int
}
type List []keyValuePair

func (l *List) Sort() {
	sort.Sort(l)
}
func (l *List) AddRange(array []string) error {
	if len(*l) > 0 {
		return errors.New("List must be empty before adding range, but it has len:" + string(l.Len()))
	}
	sort.Slice(array, func(i, j int) bool { //сортируем исходный текст по алфавиту для упаковки в слайс (с мапой без сортировки)
		return array[i] < array[j]
	})
	var lword string
	for _, word := range array {
		if word != lword {
			*l = append(*l, keyValuePair{
				key:   word,
				value: initializeValue,
			})
			lword = word
		} else {
			(*l)[l.Len()-1].key = word
			(*l)[l.Len()-1].value++
		}
	}
	return nil
}

func (l List) Len() int {
	return len(l)
}

func (l List) Less(i, j int) bool {
	return l[i].value > l[j].value
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
