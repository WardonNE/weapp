package weapp

type Model struct {
	*Database `inject:"db:default"`
	Component `inject:"component"`
}

func (m *Model) Init() {
}
