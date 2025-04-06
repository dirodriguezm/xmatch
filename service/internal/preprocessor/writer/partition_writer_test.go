package partition_writer

import (
	"log/slog"
	"os"
	"path"
	"testing"

	xwaveWriter "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestRow struct {
	Id      *string `parquet:"name=id, type=BYTE_ARRAY"`
	Column1 *int    `parquet:"name=column1, type=INT64"`
	Column2 *string `parquet:"name=column2, type=BYTE_ARRAY"`
}

func NewTestRow(id string, column1 int, column2 string) TestRow {
	return TestRow{
		Id:      &id,
		Column1: &column1,
		Column2: &column2,
	}
}

func (r TestRow) GetId() string {
	return *r.Id
}

type IdMapSuite struct {
	suite.Suite
	writer *PartitionWriter[TestRow]
}

func (s *IdMapSuite) SetupTest() {
	s.writer = &PartitionWriter[TestRow]{
		BaseWriter: &xwaveWriter.BaseWriter[TestRow]{},
		fs:         &filesystemmanager.FileSystemManager{},
	}
}

func (s *IdMapSuite) TestEmptyList() {
	rows := make([]TestRow, 0)

	result := s.writer.idMap(rows)
	require.Empty(s.T(), result)
}

func (s *IdMapSuite) TestWithRows() {
	rows := []TestRow{
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id2", 2, "value2"),
	}

	result := s.writer.idMap(rows)
	require.Len(s.T(), result, 2)
	require.Equal(s.T(), []int{0, 1}, result["id1"])
	require.Equal(s.T(), []int{2}, result["id2"])
}

func TestIdMapSuite(t *testing.T) {
	suite.Run(t, new(IdMapSuite))
}

type updateCurrentWriterSuite struct {
	suite.Suite
	writer  *PartitionWriter[TestRow]
	tempDir string
}

func (s *updateCurrentWriterSuite) SetupTest() {
	s.writer = &PartitionWriter[TestRow]{
		BaseWriter: &xwaveWriter.BaseWriter[TestRow]{},
		fs:         &filesystemmanager.FileSystemManager{},
		dirMap:     make(map[string]int),
	}
	s.tempDir = s.T().TempDir()
}

func (s *updateCurrentWriterSuite) TestWithFirstFile() {
	err := s.writer.updateCurrentWriter(s.tempDir)
	require.NoError(s.T(), err)
	require.FileExists(s.T(), path.Join(s.tempDir, "001.parquet"))
	require.NotNil(s.T(), s.writer.currentWriter)
}

func (s *updateCurrentWriterSuite) TestWithSecondFile() {
	tmpFile, _ := os.CreateTemp(s.tempDir, "test")
	s.writer.currentFile = tmpFile
	s.writer.currentWriter, _ = s.writer.createParquetWriter(s.writer.currentFile)
	s.writer.dirMap[s.tempDir] = 1

	err := s.writer.updateCurrentWriter(s.tempDir)
	require.NoError(s.T(), err)
	require.Error(s.T(), tmpFile.Close())
	require.FileExists(s.T(), path.Join(s.tempDir, "002.parquet"))
	require.NotNil(s.T(), s.writer.currentWriter)
}

func TestUpdateCurrentWriterSuite(t *testing.T) {
	suite.Run(t, new(updateCurrentWriterSuite))
}

type writeSuite struct {
	suite.Suite
	writer  *PartitionWriter[TestRow]
	tempDir string
}

func (s *writeSuite) SetupTest() {
	s.tempDir = s.T().TempDir()
	s.writer = &PartitionWriter[TestRow]{
		BaseWriter: &xwaveWriter.BaseWriter[TestRow]{},
		fs: &filesystemmanager.FileSystemManager{
			BaseDir: s.tempDir,
			Handler: partition.PartitionHandler{
				NumPartitions:   16,
				PartitionLevels: 2,
			},
		},
		maxFileSize: 300 * 1024 * 1024,
		dirMap:      make(map[string]int),
	}
	setupTestLogger(s.T())
}

func (s *writeSuite) TestWrite() {
	rows := []TestRow{
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id2", 2, "value2"),
	}

	err := s.writer.write(rows)
	require.NoError(s.T(), err)

	// now read the written files, but first we need to find them
	assignedDir, err := s.writer.fs.GetDirectory("id1")
	require.NoError(s.T(), err)
	require.FileExists(s.T(), path.Join(assignedDir, "001.parquet"))

	assignedDir, err = s.writer.fs.GetDirectory("id2")
	require.NoError(s.T(), err)
	require.FileExists(s.T(), path.Join(assignedDir, "001.parquet"))
}

func (s *writeSuite) TestWrite_when_file_is_full() {
	rows := []TestRow{
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id1", 1, "value1"),
		NewTestRow("id2", 2, "value2"),
	}

	s.writer.maxFileSize = 1
	err := s.writer.write(rows)
	require.NoError(s.T(), err)

	// now read the written files, but first we need to find them
	assignedDir, err := s.writer.fs.GetDirectory("id1")
	require.NoError(s.T(), err)
	require.FileExists(s.T(), path.Join(assignedDir, "001.parquet"), "id1")
	require.FileExists(s.T(), path.Join(assignedDir, "002.parquet"), "id1, second row")

	assignedDir, err = s.writer.fs.GetDirectory("id2")
	require.NoError(s.T(), err)
	require.FileExists(s.T(), path.Join(assignedDir, "001.parquet"), "id2")
}

func TestWriteSuite(t *testing.T) {
	suite.Run(t, new(writeSuite))
}

func setupTestLogger(t *testing.T) {
	t.Helper()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
