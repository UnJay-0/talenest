package status

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

func (Status) SetColor(color string) {

}
