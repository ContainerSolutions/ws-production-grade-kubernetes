package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // gorm
)

type gormRecord struct {
	gorm.Model
	Key   string `gorm:"not_null;unique"`
	Value string
}

// SQLDataStore is a datastore that uses an SQL server
type SQLDataStore struct {
	driver     string
	connection string
	db         *gorm.DB
}

// NewSQLDatastore returns a new SQLDataStore
func NewSQLDatastore() *SQLDataStore {
	return &SQLDataStore{}
}

// Init the SQL datastore, connect to the database
func (d *SQLDataStore) Init(parameters ...interface{}) error {
	var err error
	d.driver = parameters[0].(string)
	d.connection = parameters[1].(string)
	d.db, err = gorm.Open(d.driver, d.connection)
	if err == nil {
		d.db.AutoMigrate(&gormRecord{})
	}
	return err
}

// Get all records from the datastore
func (d *SQLDataStore) Get() []Record {
	var gormRecords []gormRecord
	records := make([]Record, 0)
	d.db.Find(&gormRecords)
	for _, gormRecord := range gormRecords {
		records = append(records, Record{gormRecord.Key, gormRecord.Value})
	}
	return records
}

// Add the Record record to the datastore
func (d *SQLDataStore) Add(record Record) {
	if len(record.Key) == 0 {
		return
	}
	d.db.Where(gormRecord{Key: record.Key}).
		Assign(gormRecord{Key: record.Key, Value: record.Value}).
		FirstOrCreate(&gormRecord{Key: record.Key, Value: record.Value})
}

// Rem the Record record from the datastore
func (d *SQLDataStore) Rem(record Record) {
	if len(record.Key) == 0 {
		return
	}
	d.db.Unscoped().Where("Key = ?", record.Key).Delete(&gormRecord{})
}

//Record contains a Key and a Value
type Record struct {
	Key   string
	Value string
}

//Datastore is a datastore, man
type Datastore interface {
	Init(...interface{}) error
	Get() []Record
	Add(Record)
	Rem(Record)
}

//SliceDataStore is a Datastore that uses a slice as data store
type SliceDataStore struct {
	slice []Record
}

//NewSliceDataStore initializes a new slice based data store
func NewSliceDataStore() *SliceDataStore {
	return &SliceDataStore{}
}

//Init initializes the SliceDataStore with initial capacity initialCapacity
func (d *SliceDataStore) Init(parameters ...interface{}) error {
	d.slice = make([]Record, parameters[0].(int))
	return nil
}

//Get all the elements in the SliceDataStore d
func (d *SliceDataStore) Get() []Record {
	return d.slice
}

//Add an element to the SliceDataStore d
func (d *SliceDataStore) Add(record Record) {
	for i, r := range d.slice {
		if r.Key == record.Key {
			d.slice[i].Value = record.Value
			return
		}
	}

	d.slice = append(d.slice, record)
}

//Rem ove an element from the SliceDataStore d
func (d *SliceDataStore) Rem(record Record) {
	for i, r := range d.slice {
		if r.Key == record.Key {
			d.slice[i] = d.slice[len(d.slice)-1]
			d.slice = d.slice[:len(d.slice)-1]
		}
	}
}

//Size of the SliceDatastore d
func (d *SliceDataStore) Size() int {
	return len(d.slice)
}