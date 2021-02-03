package integrations

type Item struct {
	Id   string // item ID
	Text string // text of list item
	Url  string // open in browser
	Copy string // copy text to clipboard
}

type ItemDetail struct {
	Title string
	Parts []Renderable
}

type Renderable interface {
	GetText() string
}

type Integration interface {
	GetName() string
	GetItems() []Item
	GetDetail(Item) ItemDetail
}
