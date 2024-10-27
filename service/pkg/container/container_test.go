package container

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Config struct {
	Url string
}
type Database interface {
	query() string
}
type SQLDatabase struct {
	config Config
}

func (db SQLDatabase) query() string {
	return db.config.Url
}

func TestRegistration(t *testing.T) {
	ctr := NewContainer()
	factory := func() Database {
		return &SQLDatabase{}
	}
	ctr.Register("db", factory)
	entry, exists := ctr[reflect.TypeOf(factory).Out(0)]
	require.True(t, exists)
	bind, exists := entry["db"]
	require.True(t, exists)
	require.NotNil(t, bind)
}

func TestResolve(t *testing.T) {
	ctr := NewContainer()
	ctr.Register("db", func() Database {
		return &SQLDatabase{config: Config{Url: "test"}}
	})

	var db Database
	err := ctr.Resolve("db", &db)
	require.NoError(t, err)
	require.NotNil(t, db)
	require.Equal(t, "test", db.query())
}

func TestMultipleRegistration(t *testing.T) {
	ctr := NewContainer()
	ctr.Register("db1", func() Database {
		return &SQLDatabase{}
	})
	ctr.Register("db2", func() Database {
		return &SQLDatabase{}
	})
	var db1 Database
	err := ctr.Resolve("db1", &db1)
	require.NoError(t, err)
	require.NotNil(t, db1)
	var db2 Database
	err = ctr.Resolve("db2", &db2)
	require.NoError(t, err)
	require.NotNil(t, db2)
}

func TestDependencyBinding(t *testing.T) {
	ctr := NewContainer()

	ctr.Register("config", func() Config {
		return Config{Url: "test"}
	})
	ctr.Register("db", func(cfg Config) Database {
		return &SQLDatabase{config: cfg}
	})

	var db Database
	err := ctr.ResolveWithBinds("db", &db, []string{"config"})
	require.NoError(t, err)
	require.NotNil(t, db)
	require.Equal(t, "test", db.query())
}
