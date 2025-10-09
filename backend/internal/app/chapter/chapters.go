package chapter

type Chapters struct {
	collection []*Chapter
}

func (c *Chapters) Add(chapter *Chapter) {
	c.collection = append(c.collection, chapter)
}

func (t *Chapters) Len() int {
	return len(t.collection)
}

func (t *Chapters) ChaptersStream() <-chan *Chapter {
	ch := make(chan *Chapter)
	go func() {
		defer close(ch)
		for _, chapter := range t.collection {
			ch <- chapter
		}
	}()
	return ch
}
