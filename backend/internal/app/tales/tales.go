package tales

type Tales struct {
	collection []*Tale
}

func (t *Tales) Add(tale *Tale) {
	t.collection = append(t.collection, tale)
}

func (t *Tales) Len() int {
	return len(t.collection)
}

func (t *Tales) TaleStream() <-chan *Tale {
	ch := make(chan *Tale)
	go func() {
		defer close(ch)
		for _, tale := range t.collection {
			ch <- tale
		}
	}()
	return ch
}
