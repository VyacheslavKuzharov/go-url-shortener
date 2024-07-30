package config

import "os"

type StorageType string

const (
	InMemory           StorageType = "mem"
	InFile             StorageType = "file"
	DefaultFileStorage             = "internal/storage/infile/urls.txt"
)

type StorageCfg struct {
	Kind StorageType
	File FileStorage
}

func NewStorageCfg(flagFile *FileStorage) *StorageCfg {
	scfg := &StorageCfg{
		Kind: InFile,
		File: FileStorage{
			Path: DefaultFileStorage,
		},
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		scfg.File.Path = os.Getenv("FILE_STORAGE_PATH")
	} else if flagFile.Path != "" {
		scfg.File = *flagFile
	}

	return scfg
}

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
