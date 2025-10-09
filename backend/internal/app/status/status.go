package status

import "fmt"

type Status struct {
	Id    int
	Name  string
	color string
}

func GetDefault() Status {
	return Status{
		Id:    1,
		Name:  "New",
		color: "008000",
	}
}

func (status Status) GetColor() string {
	return status.color
}

func (status Status) SetColor(color string) {

}

func (status Status) String() string {
	return fmt.Sprintf("status: %s, %s", status.Name, status.color)
}
