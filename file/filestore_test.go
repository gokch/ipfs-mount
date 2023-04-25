package file

import (
	"fmt"
	"io"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func TestStoreNewGet(t *testing.T) {
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Overwrite("test/abc/d/e.txt", NewWriterFromBytes(data1, cid.Cid{}))
	require.NoError(t, err)

	reader, err := fs.Get("test/abc/d/e.txt")
	require.NoError(t, err)

	data2, err := io.ReadAll(reader)
	require.Equal(t, data1, data2)
}

func TestStoreIterate(t *testing.T) {
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	// add datas
	data := []byte("test1")
	data2 := []byte("test2")
	data3 := []byte("test3")
	err := fs.Overwrite("a/data1.txt", NewWriterFromBytes(data, cid.Cid{}))
	require.NoError(t, err)
	err = fs.Overwrite("a/b/c/data2.txt", NewWriterFromBytes(data2, cid.Cid{}))
	require.NoError(t, err)
	err = fs.Overwrite("a/b/c/d/e/data3.txt", NewWriterFromBytes(data3, cid.Cid{}))
	require.NoError(t, err)

	// iterate
	err = fs.Iterate("", func(fpath string, reader *Reader) {
		out, err := io.ReadAll(reader)
		require.NoError(t, err)
		fmt.Println(reader.AbsPath(), out)
	})
	require.NoError(t, err)
}