package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

var taskWithAsteriskIsCompleted = true //не пойму, почему располагая переменную в файле top.go-она видна и в top_test.go? А наоборот нет?

const (
	top = 10
)

func Top10(in string) []string {
	if in == "" {
		return []string{}
	}
	list := List{}
	var strspl []string
	if taskWithAsteriskIsCompleted {
		asteriskstr := strings.ReplaceAll(in, " - ", "")
		asteriskstr = strings.ToLower(asteriskstr)
		strspl = strings.FieldsFunc(asteriskstr, func(r rune) bool {
			return strings.ContainsRune(" ;,.!&\"\t\r\n", r)
		})
	} else {
		strspl = strings.FieldsFunc(in, func(r rune) bool {
			return strings.ContainsRune(" \t\r\n", r)
		})
	}
	/*dict := map[string]int{} //можно сделать и мапой, но решил что слайс структур красивее
	for _, word := range strspl {
		dict[word]++
	}*/
	sort.Slice(strspl, func(i, j int) bool { //сортируем исходный текст по алфавиту для упаковки в слайс (с мапой без сортировки)
		return strspl[i] < strspl[j]
	})
	err := list.AddRange(strspl)
	if err != nil {
		panic(err)
	}
	list.Sort()
	var top10 = make([]string, 0, top)
	for _, word := range list[0:top] {
		top10 = append(top10, word.key)
	}
	return top10
}
