package partition_writer

import (
	"fmt"
	"os"

	xwaveWriter "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	parquetWriter "github.com/xitongsys/parquet-go/writer"
)

type RowWithId interface {
	GetId() string
}

type PartitionWriter[T RowWithId] struct {
	*xwaveWriter.BaseWriter[T]
	fs            *filesystemmanager.FileSystemManager
	maxFileSize   int
	currentWriter *parquetWriter.ParquetWriter
	currentFile   *os.File
	dirMap        map[string]int
	currentDir    string
}

func (w *PartitionWriter[T]) write(rows []T) error {
	idMap := w.idMap(rows)

	for id := range idMap {
		dir, err := w.fs.GetDirectory(id)
		if err != nil {
			return fmt.Errorf("Error while getting directory for id %s\n%w", id, err)
		}

		dirChanged := w.currentDir != dir
		if dirChanged {
			w.currentDir = dir
			_, err := w.fs.CreatePartition(id)
			if err != nil {
				return fmt.Errorf("Error while creating partition for id %s\n%w", id, err)
			}
		}

		if w.currentFile == nil || w.fs.GetSizeOfFile(w.currentFile) > int64(w.maxFileSize) || dirChanged {
			err := w.updateCurrentWriter(w.currentDir)
			if err != nil {
				return fmt.Errorf("Error while updating current writer and file\n%w", err)
			}
		}

		for _, rowIdx := range idMap[id] {
			if err := w.currentWriter.Write(rows[rowIdx]); err != nil {
				return fmt.Errorf("ParquetWriter could not write object %v\n%w", rows[rowIdx], err)
			}
			if w.fs.GetSizeOfFile(w.currentFile) > int64(w.maxFileSize) {
				err := w.updateCurrentWriter(w.currentDir)
				if err != nil {
					return fmt.Errorf("Error updating current writer and file due to file size\n%w", err)
				}
			}
		}
	}

	return nil
}

// Create a map of object ids and the indexes of their rows
func (w *PartitionWriter[T]) idMap(rows []T) map[string][]int {
	idMap := make(map[string][]int)

	for i, row := range rows {
		id := getId(row)
		if _, ok := idMap[id]; !ok {
			idMap[id] = make([]int, 0)
		}
		idMap[id] = append(idMap[id], i)
	}

	return idMap
}

func (w *PartitionWriter[T]) updateCurrentWriter(dir string) error {
	if err := w.closeResources(); err != nil {
		return fmt.Errorf(
			"Error closing resources while changing to a new file\n%w",
			err,
		)
	}

	updateCurrentDirFileIndex(w.dirMap, dir)

	var err error
	w.currentFile, err = w.fs.CreateNewFile(dir, w.dirMap[dir])
	if err != nil {
		return fmt.Errorf("Error while creating file\n%w", err)
	}

	w.currentWriter, err = w.createParquetWriter(w.currentFile)
	if err != nil {
		return fmt.Errorf("ParquetGo could not create writer from writer\n%w", err)
	}

	return nil
}

func (w *PartitionWriter[T]) createParquetWriter(file *os.File) (*parquetWriter.ParquetWriter, error) {
	schema := new(T)
	writer, err := parquetWriter.NewParquetWriterFromWriter(file, schema, 1)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *PartitionWriter[T]) closeResources() error {
	if w.currentWriter == nil || w.currentFile == nil {
		return nil
	}

	if err := w.currentWriter.WriteStop(); err != nil {
		return fmt.Errorf("ParquetWriter could not stop. %w", err)
	}

	if err := w.currentFile.Close(); err != nil {
		return fmt.Errorf("ParquetWriter could not close parquet file %w", err)
	}

	return nil
}

func updateCurrentDirFileIndex(dirMap map[string]int, dir string) {
	if _, ok := dirMap[dir]; !ok {
		dirMap[dir] = 1
	} else {
		dirMap[dir]++
	}
}

func getId[T RowWithId](row T) string {
	return row.GetId()
}
