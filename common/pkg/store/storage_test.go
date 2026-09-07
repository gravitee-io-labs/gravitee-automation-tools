package store_test

import (
	"context"
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/store"
	"github.com/stretchr/testify/assert"
)

type testObject struct {
	key  string
	data string
}

func (k testObject) Identity() string {
	return k.key
}

var underTest *store.Store[testObject]

func before() {
	underTest = store.NewStoreWithData[testObject](context.Background(),
		testObject{
			key:  "foo",
			data: "bar",
		})
}

func TestGet(t *testing.T) {
	before()
	x, ok := underTest.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", x.data)
	_, ok = underTest.Get("unknown")
	assert.False(t, ok)
}

func TestPut(t *testing.T) {
	before()
	underTest.Put(testObject{"new", "hey!"})
	x, _ := underTest.Get("new")
	assert.Equal(t, "hey!", x.data)

	underTest.Put(testObject{"foo", "bar updated"})
	x, _ = underTest.Get("foo")

	assert.Equal(t, "bar updated", x.data)
	x, _ = underTest.Get("new")
	assert.Equal(t, "hey!", x.data)
}

func TestDelete(t *testing.T) {
	before()
	_, ok := underTest.Get("foo")
	assert.True(t, ok)
	underTest.Delete(testObject{"foo", "bar"})
	_, ok = underTest.Get("foo")
	assert.False(t, ok)
	assert.NotPanics(t, func() {
		underTest.Delete(testObject{"foo", "bar"})
	})
}

func TestDeleteByKey(t *testing.T) {
	before()
	_, ok := underTest.Get("foo")
	assert.True(t, ok)
	underTest.DeleteByKey("foo")
	_, ok = underTest.Get("foo")
	assert.False(t, ok)
	assert.NotPanics(t, func() {
		underTest.DeleteByKey("foo")
	})
}

func TestGetAll(t *testing.T) {
	before()
	underTest.Put(testObject{"new", "hey!"})
	underTest.Put(testObject{"hello", "world"})
	all := underTest.GetAll()
	assert.Len(t, all, 3)
	assert.Equal(t, all[0], testObject{"foo", "bar"})
	assert.Equal(t, all[1], testObject{"new", "hey!"})
	assert.Equal(t, all[2], testObject{"hello", "world"})

	underTest.DeleteByKey("new")

	assert.Len(t, all, 3)
	assert.Equal(t, all[0], testObject{"foo", "bar"})
	assert.Equal(t, all[1], testObject{"new", "hey!"})

}

func TestGetPage(t *testing.T) {
	before()
	underTest.Put(testObject{"new", "hey!"})
	underTest.Put(testObject{"hello", "world"})
	all := underTest.GetPage(1, 10)
	assert.Len(t, all, 3)
	assert.Contains(t, all,
		testObject{"foo", "bar"},
		testObject{"new", "hey!"},
		testObject{"hello", "world"})
	page := underTest.GetPage(1, 2)
	assert.Len(t, page, 2)
	assert.Contains(t, page,
		testObject{"foo", "bar"},
		testObject{"new", "hey!"})

	underTest.Put(testObject{"forth", "4"})
	underTest.Put(testObject{"fifth", "5"})
	underTest.Put(testObject{"sixth", "6"})
	underTest.Put(testObject{"seventh", "7"})
	underTest.Put(testObject{"eighth", "8"})
	underTest.Put(testObject{"ninth", "9"})
	underTest.Put(testObject{"tenth", "10"})
	underTest.Put(testObject{"eleventh", "11"})

	page = underTest.GetPage(-1, -1)
	assert.Len(t, page, 10)
	assert.Equal(t, page[9], testObject{"tenth", "10"})
	assert.NotContains(t, page, testObject{"eleventh", "11"})

	page = underTest.GetPage(2, 10)
	assert.Len(t, page, 1)
	assert.Equal(t, page[0], testObject{"eleventh", "11"})

	page = underTest.GetPage(3, 3)
	assert.Len(t, page, 3)
	assert.Equal(t, page[0], testObject{"seventh", "7"})
	assert.Equal(t, page[1], testObject{"eighth", "8"})
	assert.Equal(t, page[2], testObject{"ninth", "9"})
}

func TestEmptyStore(t *testing.T) {
	underLocalTest := store.NewStore[testObject](context.Background())
	_, ok := underLocalTest.Get("foo")
	assert.False(t, ok)

	assert.NotPanics(t, func() {
		underLocalTest.DeleteByKey("foo")
	})
	assert.NotPanics(t, func() {
		underLocalTest.Delete(testObject{"foo", "bar"})
	})
	assert.Len(t, underLocalTest.GetAll(), 0)
	assert.Len(t, underLocalTest.GetPage(1, 10), 0)

}
