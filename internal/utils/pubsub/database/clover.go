package database

import (
	"github.com/ostafen/clover"
)

type CloverDatabase struct {
	name string
	db   *clover.DB
}

func NewCloverDatabase(db *clover.DB) Databaser {
	name := "x_system_storage_pubsub_bck"
	db.CreateCollection(name)
	return &CloverDatabase{
		db:   db,
		name: name,
	}
}

func (db *CloverDatabase) getLastSequenceId() (int64, error) {
	var backup Backup
	doc, err := db.db.Query(db.name).
		Sort(clover.SortOption{
			Field:     "sequence_id",
			Direction: -1}).
		FindFirst()
	if err != nil {
		return 0, err
	}

	err = doc.Unmarshal(&backup)
	if err != nil {
		return 0, err
	}

	return backup.SequenceId, nil
}

func (db *CloverDatabase) delete(topic string) error {
	return db.db.Query(db.name).Where(clover.Field("topic").Eq(topic)).Delete()
}

func (db *CloverDatabase) RetrieveAndDelete(topic string) (Backups, error) {
	var backups Backups

	docs, err := db.db.Query(db.name).
		Where(clover.Field("topic").Eq(topic)).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		backup := Backup{}
		err := doc.Unmarshal(&backup)
		if err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}

	db.delete(topic)

	return backups, nil
}

func (db *CloverDatabase) Store(data Backup) error {
	sequenceId, err := db.getLastSequenceId()
	if err != nil {
		return err
	}

	data.SequenceId = sequenceId
	dataMap, err := data.MapStringInterface()
	if err != nil {
		return err
	}

	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	_, err = db.db.InsertOne(db.name, doc)
	if err != nil {
		return err
	}

	return nil
}
