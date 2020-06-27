package hw04_lru_cache //nolint:golint,stylecheck

// List ...
type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	len       int
	frontItem *listItem
	backItem  *listItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.frontItem
}

func (l *list) Back() *listItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *listItem {
	new := listItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	}
	if l.frontItem != nil {
		l.insertBefore(l.frontItem, &new)
	} else {
		l.backItem = &new
	}
	l.frontItem = &new
	l.len++

	return &new
}

func (l *list) PushBack(v interface{}) *listItem {
	new := listItem{
		Value: v,
		Prev:  l.backItem,
		Next:  nil,
	}
	if l.backItem != nil {
		l.backItem.Next = &new
	} else {
		l.frontItem = &new
	}
	l.backItem = &new
	l.len++

	return &new
}

func (l *list) Remove(i *listItem) {
	l.pop(i)
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if i.Prev == nil {
		return
	}
	l.insertBefore(l.frontItem, l.pop(i))
	l.frontItem = i
}

func (l *list) pop(i *listItem) *listItem {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.frontItem = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
	i.Next = nil
	i.Prev = nil

	return i
}

func (l *list) insertBefore(root *listItem, new *listItem) {
	new.Next = root
	if root != nil {
		root.Prev = new
	}
}

// NewList ...
func NewList() List {
	return &list{}
}
