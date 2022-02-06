package imagecache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
	}

	if l.head == nil {
		l.head = newListItem
		l.tail = newListItem
	} else {
		current := l.head
		newListItem.Next = current
		newListItem.Prev = nil
		current.Prev = newListItem
		l.head = newListItem
	}

	l.len++

	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
	}

	if l.head == nil {
		l.head = newListItem
		l.tail = newListItem
	} else {
		current := l.tail
		current.Next = newListItem
		newListItem.Prev = current
		newListItem.Next = nil
		l.tail = newListItem
	}

	l.len++

	return newListItem
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	if l.Back() == i {
		l.tail = i.Prev
		i.Prev.Next = nil
		return
	}
	if l.Front() == i {
		l.head = i.Next
		i.Next.Prev = nil
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len == 1 || l.head == i {
		return
	}
	if l.Back() == i {
		i.Prev.Next = nil
		i.Next = l.head
		l.tail = i.Prev
		l.head.Prev = i
		l.head = i
		i.Prev = nil
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}
