package chapter

import "fmt"

type Chapter struct {
	Id        int
	Content   string
	sentiment float64
	TaleId    int
}

func (chapter *Chapter) String() string {
	return fmt.Sprintf("Chapter %d [%d]:\n%s\n", chapter.Id, chapter.TaleId, chapter.Content)
}
