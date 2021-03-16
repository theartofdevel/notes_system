package freecache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var uuid = []byte("Lorem ipsum")
var data = []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
Nulla id tincidunt urna. Proin auctor pretium ornare. Donec vitae felis est.
Sed sed venenatis ex. Nunc semper vel quam sit amet molestie.
Quisque lacus tortor, convallis eget ante tristique, interdum auctor massa.
Praesent a lacus tristique, facilisis tellus eget, malesuada urna.
Vivamus hendrerit posuere mauris, nec rhoncus turpis suscipit et. Proin at risus ac odio laoreet imperdiet.
Phasellus ex nulla, sagittis sed tempus maximus, sagittis at dolor.
Aliquam erat volutpat. In tincidunt eros quis pharetra efficitur. 
Phasellus facilisis sagittis porta. Curabitur sed dui non ligula malesuada aliquet eget sed lectus. 
Interdum et malesuada fames ac ante ipsum primis in faucibus. Nunc eget auctor felis. 
Vestibulum vel metus eu velit molestie hendrerit. Sed pharetra vel arcu et efficitur. 
Maecenas enim mauris, efficitur venenatis libero sit amet, finibus vulputate est. 
Sed augue ex, viverra a pretium malesuada, elementum vitae aenean. `)

var repo = NewCacheRepo(104857600) // 100MB

// Test scenario:
// 1. Set some data
// 2. Check entry count == 1
// 3. Get data
// 4. Check hit count == 1
// 5. Get data with invalid key
// 6. Check miss count == 1
// 7. Delete data using valid key
// 8. Delete data using invalid key (must fail)
func TestRepository(t *testing.T) {
	err := repo.Set(uuid, data, -1)
	assert.NoError(t, err, "failed to set data in cache")

	entryCount := repo.EntryCount()
	assert.Equal(t, entryCount, int64(1))

	entry, err := repo.Get(uuid)
	if assert.NoError(t, err, "failed to get entry from cache") {
		assert.Equal(t, entry, data)
	}

	hitCount := repo.HitCount()
	assert.Equal(t, hitCount, int64(1))

	entry, err = repo.Get([]byte("invalid key"))
	if assert.Error(t, err, "failed to get entry from cache") {
		assert.Nil(t, entry)
	}

	missCount := repo.HitCount()
	assert.Equal(t, missCount, int64(1))

	affected := repo.Del(uuid)
	assert.Equal(t, affected, true, "failed to delete entry from cache")

	affected = repo.Del([]byte("invalid key"))
	assert.Equal(t, affected, false)
}

func TestRepositoryLoad(t *testing.T) {
	entryNum := 10000
	keys := [][]byte{}
	for i := 0; i < entryNum; i++ {
		keys = append(keys, []byte(strconv.FormatInt(int64(i), 10)))
	}

	var err error

	for i := 0; i < entryNum; i++ {
		err = repo.Set(keys[i], data, -1)
		if err != nil {
			t.Errorf("failed to set entry to cache: %v", err)
		}
	}

	var entry []byte

	for i := 0; i < entryNum; i++ {
		entry, err = repo.Get(keys[i])
		if assert.NoError(t, err) {
			assert.Equal(t, entry, data)
		}
	}

	var affected bool

	for i := 0; i < entryNum; i++ {
		affected = repo.Del(keys[i])
		assert.Equal(t, affected, true)
	}
}
