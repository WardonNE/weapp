package weapp

type IMigration interface {
	Up() error
	Down() error
	Database() string
	Version() string
}

type Migration struct {
}

func (m *Migration) Database() string {
	return "default"
}

func (m *Migration) Up() error {
	return nil
}

func (m *Migration) Down() error {
	return nil
}
