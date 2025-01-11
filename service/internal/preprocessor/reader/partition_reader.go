package partition_reader

import (
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
)

type PartitionReader[T any] struct {
	fs     *filesystemmanager.FileSystemManager
	reader parquet_reader.ParquetReader[T]
}
