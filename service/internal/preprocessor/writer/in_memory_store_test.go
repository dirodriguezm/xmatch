package partition_writer

import (
	"io"
	"log/slog"
	"path"
	"testing"

	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/suite"
)

type InMemoryStoreSuite struct {
	suite.Suite
	store  *InMemoryStore
	tmpDir string
}

func setupTestLogger(t *testing.T, stdout io.Writer) {
	t.Helper()

	handler := slog.NewTextHandler(stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func (suite *InMemoryStoreSuite) SetupSuite(stdout io.Writer) {
	setupTestLogger(suite.T(), stdout)
}

func (suite *InMemoryStoreSuite) SetupTest() {
	fs := filesystemmanager.AFileSystemManager(suite.T()).
		WithBaseDir(suite.T().TempDir()).
		WithNumLevels(1).
		WithNumPartitions(1).
		Build()
	suite.store = NewInMemoryStore(2, &fs)
	suite.tmpDir = suite.store.fs.BaseDir
}

func NewTestRow(id string, column1 int, column2 string) TestInputSchema {
	return TestInputSchema{
		Id:      &id,
		Column1: &column1,
		Column2: &column2,
	}
}

func (suite *InMemoryStoreSuite) TestWriteSuccess() {
	rows := []repository.InputSchema{
		NewTestRow("1", 1, "value1"),
		NewTestRow("2", 2, "value2"),
	}

	toFlush, err := suite.store.Write(rows)

	suite.NoError(err)
	suite.Empty(toFlush)

	dir := path.Join(suite.tmpDir, "0")
	suite.Len(suite.store.store[dir], 2)
}

func (suite *InMemoryStoreSuite) TestWritePartitionFull() {
	rows := []repository.InputSchema{
		NewTestRow("1", 1, "value1"),
		NewTestRow("1", 2, "value2"),
		NewTestRow("1", 3, "value3"),
	}

	toFlush, err := suite.store.Write(rows)

	suite.NoError(err)
	dir := path.Join(suite.tmpDir, "0")
	suite.Len(toFlush[dir], 3)
	suite.Len(suite.store.store[dir], 0)
}

func (suite *InMemoryStoreSuite) TestCanWrite() {
	suite.store.store["001"] = []repository.InputSchema{
		NewTestRow("1", 1, "value1"),
	}

	canWrite := suite.store.canWrite("001")
	suite.True(canWrite)

	suite.store.store["001"] = append(suite.store.store["001"], NewTestRow("2", 2, "value2"))
	canWrite = suite.store.canWrite("001")
	suite.False(canWrite)
}

func TestInMemoryStoreSuite(t *testing.T) {
	suite.Run(t, new(InMemoryStoreSuite))
}
