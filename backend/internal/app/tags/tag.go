package tags

import "fmt"

type Tag struct {
	Id   int
	Name string
}

func (tag Tag) String() string {
	return fmt.Sprintf("Tag %d: %s", tag.Id, tag.Name)
}
