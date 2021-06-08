package store

import (
	"testing"

	"github.com/nityanandagohain/log-kv-store/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	key := utils.RandomKey()
	val := utils.RandomValue()
	err := testStore.Put(key, val)
	require.NoError(t, err)

	//

	storedVal, err := testStore.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, storedVal)

}

func BenchmarkWrite(b *testing.B) {
	for n := 0; n < b.N; n++ {
		key := utils.RandomKey()
		val := utils.RandomValue()
		err := testStore.Put(key, val)
		require.NoError(b, err)
		storedVal, err := testStore.Get(key)
		require.NoError(b, err)
		require.Equal(b, val, storedVal)
	}
}
