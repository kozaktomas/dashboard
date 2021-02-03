package utils

// paragraph
type Paragraph struct {
	Text string
}

func (p Paragraph) GetText() string {
	return p.Text
}

// empty line
type Break struct {

}

func (p Break) GetText() string {
	return "\n"
}