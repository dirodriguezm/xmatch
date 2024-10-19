package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Database interface {
	query()
}

type DatabaseImp1 struct{}

func (d *DatabaseImp1) query() {}

type DatabaseImp2 struct{}

func (d *DatabaseImp2) query() {}

func TestContainer(t *testing.T) {
	c := NewContainer()

	c.register("db1", func() Database {
		return &DatabaseImp1{}
	})
	c.register("db2", func() Database {
		return &DatabaseImp2{}
	})

	var db1, db2 Database
	err := c.get(&db1, "db1")
	require.Nil(t, err)
	err = c.get(&db2, "db2")
	require.Nil(t, err)

	require.NotNil(t, db1)
	require.NotNil(t, db2)

	require.IsType(t, &DatabaseImp1{}, db1)
	require.IsType(t, &DatabaseImp2{}, db2)
}

func TestRegister(t *testing.T) {
	c := NewContainer()
	Register(c, "db1", func() Database {
		return &DatabaseImp1{}
	})
	Register(c, "db2", func() Database {
		return &DatabaseImp2{}
	})

	var db1, db2 Database
	err := c.get(&db1, "db1")
	require.Nil(t, err)
	err = c.get(&db2, "db2")
	require.Nil(t, err)

	require.NotNil(t, db1)
	require.NotNil(t, db2)

	require.IsType(t, &DatabaseImp1{}, db1)
	require.IsType(t, &DatabaseImp2{}, db2)
}

func TestGet(t *testing.T) {
	c := NewContainer()

	c.register("db1", func() Database {
		return &DatabaseImp1{}
	})
	c.register("db2", func() Database {
		return &DatabaseImp2{}
	})

	db1, err := Get[Database](c, "db1")
	require.NoError(t, err)
	db2, err := Get[Database](c, "db2")
	require.NoError(t, err)
	require.NotNil(t, db1)
	require.NotNil(t, db2)

	require.IsType(t, &DatabaseImp1{}, db1)
	require.IsType(t, &DatabaseImp2{}, db2)
}

func TestRegisterExistingKey(t *testing.T) {
	c := NewContainer()
	err := Register(c, "db1", func() Database {
		return &DatabaseImp1{}
	})
	require.Nil(t, err)
	err = Register(c, "db1", func() Database {
		return &DatabaseImp1{}
	})
	require.Error(t, err)
}

func TestRegisterWrongArgument(t *testing.T) {
	c := NewContainer()
	err1 := c.register("something", "something")
	require.Error(t, err1)
	err2 := c.register("something", func(arg1 string) {})
	require.Error(t, err2)
	err3 := c.register("something", func() (int, int) { return 1, 1 })
	require.Error(t, err3)
}

func TestGetNotFound(t *testing.T) {
	c := NewContainer()
	resolved, err := Get[any](c, "something")
	require.Zero(t, resolved)
	require.Error(t, err)
}

func TestGetNoPointer(t *testing.T) {
	c := NewContainer()
	err1 := c.get("a", "a")
	require.Error(t, err1)
	err2 := c.get(nil, "a")
	require.Error(t, err2)

	err := Register(c, "db", func() Database {
		return &DatabaseImp1{}
	})
	require.Nil(t, err)

	type NotDatabase interface {
		notdb()
	}
	var notDb NotDatabase
	err3 := c.get(&notDb, "db")
	require.Error(t, err3)
}

func TestSingleton(t *testing.T) {
	c := NewContainer()

	err := Register(c, "db", func() Database {
		return &DatabaseImp1{}
	})
	require.NoError(t, err)

	db, err := Get[Database](c, "db")
	require.NoError(t, err)
	require.IsType(t, &DatabaseImp1{}, db)

	db2, err := Get[Database](c, "db")
	require.NoError(t, err)
	require.IsType(t, &DatabaseImp1{}, db2)

	require.True(t, db == db2)
}
