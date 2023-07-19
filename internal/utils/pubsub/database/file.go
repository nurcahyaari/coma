package database

type FileDatabase struct {
	path string
}

func NewFileDatabase(path string) Databaser {
	return &FileDatabase{
		path: path,
	}
}

func (f *FileDatabase) Retrieve() (Backups, error) {
	return nil, nil
}

func (f *FileDatabase) Store(data Backup) error {
	return nil
}
