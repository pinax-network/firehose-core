package mindreader

import (
	"testing"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func testBlock(num uint64, id string) *pbbstream.Block {
	return &pbbstream.Block{Number: num, Id: id, LibNum: num - 1}
}

func TestCursorTracker_AdvancesOverStoredPrefixOnly(t *testing.T) {
	dir := t.TempDir()
	tracker, err := NewCursorTracker(dir, zap.NewNop())
	require.NoError(t, err)

	for i, id := range []string{"a", "b", "c"} {
		tracker.Queued(id, testBlock(uint64(10+i), id))
	}

	// "b" and "c" stored before "a": the cursor waits for "a"
	tracker.Uploaded("c")
	tracker.Uploaded("b")
	tracker.Flush()
	cursor, err := ReadReaderCursor(tracker.Path())
	require.NoError(t, err)
	assert.Nil(t, cursor)

	tracker.Uploaded("a")
	tracker.Flush()
	cursor, err = ReadReaderCursor(tracker.Path())
	require.NoError(t, err)
	require.NotNil(t, cursor)
	assert.Equal(t, CursorBlockRef{Num: 12, ID: "c"}, cursor.Block)
	assert.Equal(t, uint64(11), cursor.LibNum)
}

func TestCursorTracker_Values(t *testing.T) {
	dir := t.TempDir()
	tracker, err := NewCursorTracker(dir, zap.NewNop())
	require.NoError(t, err)

	// a value before its block is queued, one while queued, one after it was stored
	tracker.SetValue(10, "a", "va")
	tracker.Queued("a", testBlock(10, "a"))
	tracker.Queued("b", testBlock(11, "b"))
	tracker.SetValue(11, "b", "vb")
	tracker.Queued("c", testBlock(12, "c"))

	tracker.Uploaded("a")
	tracker.Flush()
	cursor, _ := ReadReaderCursor(tracker.Path())
	assert.Equal(t, &CursorValue{Block: CursorBlockRef{Num: 10, ID: "a"}, Value: "va"}, cursor.Value)

	tracker.Uploaded("b")
	tracker.Uploaded("c")
	tracker.Flush()
	cursor, _ = ReadReaderCursor(tracker.Path())
	assert.Equal(t, CursorBlockRef{Num: 12, ID: "c"}, cursor.Block)
	assert.Equal(t, &CursorValue{Block: CursorBlockRef{Num: 11, ID: "b"}, Value: "vb"}, cursor.Value, "c has no value yet")

	tracker.SetValue(12, "c", "vc")
	tracker.Flush()
	cursor, _ = ReadReaderCursor(tracker.Path())
	assert.Equal(t, &CursorValue{Block: CursorBlockRef{Num: 12, ID: "c"}, Value: "vc"}, cursor.Value)

	// a restart keeps the cursor
	tracker, err = NewCursorTracker(dir, zap.NewNop())
	require.NoError(t, err)
	tracker.Queued("d", testBlock(13, "d"))
	tracker.Uploaded("d")
	tracker.Flush()
	cursor, _ = ReadReaderCursor(tracker.Path())
	assert.Equal(t, CursorBlockRef{Num: 13, ID: "d"}, cursor.Block)
	assert.Equal(t, "vc", cursor.Value.Value)
}

func TestCursorTracker_BlockPrintedTwice(t *testing.T) {
	tracker, err := NewCursorTracker(t.TempDir(), zap.NewNop())
	require.NoError(t, err)

	tracker.Queued("a", testBlock(10, "a"))
	tracker.Queued("a", testBlock(10, "a"))
	tracker.Queued("b", testBlock(11, "b"))
	tracker.Uploaded("a")
	tracker.Uploaded("b")

	tracker.Flush()
	cursor, _ := ReadReaderCursor(tracker.Path())
	assert.Equal(t, CursorBlockRef{Num: 11, ID: "b"}, cursor.Block)
}
