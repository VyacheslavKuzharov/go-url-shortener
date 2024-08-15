package config

import (
	"os"
)

type StorageType string

const (
	InMemory           StorageType = "mem"
	InFile             StorageType = "file"
	Postgres           StorageType = "pg"
	DefaultFileStorage             = "internal/storage/infile/urls.txt"
)

type StorageCfg struct {
	Kind     StorageType
	File     FileStorage
	Postgres PgStorage
}

func NewStorageCfg(flagFile *FileStorage, flagDB *PgStorage) *StorageCfg {
	scfg := &StorageCfg{
		Kind: InMemory,
		File: FileStorage{
			Path: DefaultFileStorage,
		},
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		scfg.Kind = InFile
		scfg.File.Path = os.Getenv("FILE_STORAGE_PATH")
	} else if flagFile.Path != "" {
		scfg.Kind = InFile
		scfg.File = *flagFile
	}

	if os.Getenv("DATABASE_DSN") != "" {
		scfg.Kind = Postgres
		scfg.Postgres.ConnectURL = os.Getenv("DATABASE_DSN")
	} else if flagDB.ConnectURL != "" {
		scfg.Kind = Postgres
		scfg.Postgres = *flagDB
	}

	return scfg
}

// FileStorage - describes File Storage connection
type FileStorage struct {
	Path string
}

func (f *FileStorage) String() string {
	return f.Path
}

func (f *FileStorage) Set(s string) error {
	f.Path = s
	return nil
}

// PgStorage - describes Postgres database connection
type PgStorage struct {
	ConnectURL string
}

func (pg *PgStorage) String() string {
	return pg.ConnectURL
}

func (pg *PgStorage) Set(s string) error {
	pg.ConnectURL = s
	return nil
}
