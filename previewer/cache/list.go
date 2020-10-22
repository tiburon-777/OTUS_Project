package cache

type ListInterface interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый Item
	Back() *ListItem                   // последний Item
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type List struct {
	Info ListItem
	len  int
}

func NewList() *List {
	return &List{len: 0}
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Front() *ListItem {
	return l.Info.Next
}

func (l *List) Back() *ListItem {
	return l.Info.Prev
}

func (l *List) PushFront(v interface{}) *ListItem {
	e := &ListItem{Value: v}
	if l.len != 0 {
		e.Prev = l.Info.Next
		l.Info.Next.Next = e
	} else {
		l.Info.Prev = e
	}
	l.Info.Next = e
	l.len++
	return e
}

func (l *List) PushBack(v interface{}) *ListItem {
	e := &ListItem{Value: v}
	if l.len != 0 {
		e.Next = l.Info.Prev
		l.Info.Prev.Prev = e
	} else {
		l.Info.Next = e
	}
	l.Info.Prev = e
	l.len++
	return e
}

func (l *List) Remove(i *ListItem) {
	if i.Prev == nil {
		i.Prev = &ListItem{}
		l.Info.Prev = i.Next
	}
	if i.Next == nil {
		i.Next = &ListItem{}
		l.Info.Next = i.Prev
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	if l.Len() > 1 {
		l.Info.Next.Next = nil
		l.Info.Prev.Prev = nil
	} else {
		l.Info = ListItem{}
	}
	l.len--
}

func (l *List) MoveToFront(i *ListItem) {
	if l.Info.Next == i {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	} else {
		i.Next.Prev = i.Prev
		l.Info.Prev = i.Next
	}
	i.Prev = l.Front()
	l.Front().Next = i
	i.Next = nil
	l.Info.Next = i
}
