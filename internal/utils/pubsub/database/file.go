package database

type FileDatabase struct {
	path string
}

func NewFileDatabase(path string) Databaser {
	return &FileDatabase{
		path: path,
	}
}

func (f *FileDatabase) RetrieveAndDelete(topic string) (Backups, error) {
	// todo: implement golang file access
	return nil, nil
}

func (f *FileDatabase) Store(data Backup) error {
	// todo: implement golang file access
	return nil
}
