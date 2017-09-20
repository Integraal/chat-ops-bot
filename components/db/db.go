package db

var database = &Database{
	events: make(map[string]*Event),
}

// TODO: real database storage

type Database struct {
	events map[string]*Event
}

func Get() *Database {
	return database
}
