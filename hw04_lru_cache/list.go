package hw04_lru_cache //nolint:golint,stylecheck
import (
	"github.com/pkg/errors"
)

var _ List = (*list)(nil)

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem) error          // удалить элемент
	MoveToFront(i *listItem) error     // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Prev  *listItem   // следующий элемент
	Next  *listItem   // предыдущий элемент

}

type list struct {
	len   int
	back  *listItem
	front *listItem
}

func isNilItem(i *listItem) error {
	if i == nil {
		return errors.New("item is nil")
	}
	return nil
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	item := new(listItem)
	item.Value = v
	item.Next = nil
	if l.front == nil { //пустой список
		item.Prev = nil
		l.back = item
	} else {
		item.Prev = l.front
		l.front.Next = item
	}
	l.front = item
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *listItem {
	item := new(listItem)
	item.Value = v
	item.Prev = nil
	if l.back == nil { //пустой список
		item.Next = nil
		l.front = item
	} else {
		item.Next = l.back
		l.back.Prev = item
	}
	l.back = item
	l.len++
	return item
}

func (l *list) Remove(i *listItem) error {
	if innerError := isNilItem(i); innerError != nil {
		return errors.Wrap(innerError, "can't remove item")
	}
	prev := i.Prev
	next := i.Next
	if prev != nil {
		prev.Next = next
	} else { //удаляем back элемент
		l.back = i.Next
	}
	if next != nil {
		next.Prev = prev
	} else { //удаляем front элемент
		l.front = i.Prev
	}
	l.len--
	return nil
}

func (l *list) MoveToFront(i *listItem) error {
	if err := l.Remove(i); err != nil {
		return errors.Wrap(err, "can't move to front item")
	}
	l.PushFront(i.Value)
	return nil
}

func NewList() List {
	return &list{}
}
