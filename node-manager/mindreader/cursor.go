package mindreader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"go.uber.org/zap"
)

// ReaderCursorFileName is the name of the reader cursor file, in the reader node's working
// directory.
const ReaderCursorFileName = "reader-cursor.json"

// ReaderCursor is the last block the reader node stored in its one-block store, every block it
// read before that one being stored too. A reader restarts from it instead of from what it
// printed last: blocks it printed after the cursor may never have been stored.
type ReaderCursor struct {
	Block  CursorBlockRef `json:"block"`
	LibNum uint64         `json:"lib_num"`

	// Value is the reader's own cursor ("FIRE CURSOR" line) for the last stored block that had
	// one, at or before Block.
	Value *CursorValue `json:"value,omitempty"`
}

type CursorBlockRef struct {
	Num uint64 `json:"num"`
	ID  string `json:"id"`
}

type CursorValue struct {
	Block CursorBlockRef `json:"block"`
	Value string         `json:"value"`
}

func ReadReaderCursor(path string) (*ReaderCursor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading reader cursor: %w", err)
	}

	cursor := &ReaderCursor{}
	if err := json.Unmarshal(data, cursor); err != nil {
		return nil, fmt.Errorf("decoding reader cursor %q: %w", path, err)
	}
	return cursor, nil
}

// CursorTracker follows blocks from the archiver to the one-block store and keeps the reader
// cursor file at the last block stored without a gap before it. Blocks are queued in the
// order the reader printed them, and uploads complete in any order.
type CursorTracker struct {
	path   string
	logger *zap.Logger

	mu      sync.Mutex
	queue   []*trackedBlock
	byFile  map[string][]*trackedBlock // a block printed twice is queued twice under one file
	values  map[string]string          // reader cursor values waiting for their block, by block id
	current *ReaderCursor
	// recent are the last stored blocks, so a cursor value arriving after its block was
	// stored still applies
	recent map[string]uint64
}

type trackedBlock struct {
	ref      CursorBlockRef
	libNum   uint64
	value    *string
	uploaded bool
}

const recentStoredBlocks = 1000

func NewCursorTracker(workingDirectory string, logger *zap.Logger) (*CursorTracker, error) {
	path := filepath.Join(workingDirectory, ReaderCursorFileName)
	current, err := ReadReaderCursor(path)
	if err != nil {
		return nil, err
	}

	logger.Info("reader cursor loaded", zap.String("path", path), zap.Any("cursor", current))
	return &CursorTracker{
		path:    path,
		logger:  logger,
		byFile:  map[string][]*trackedBlock{},
		values:  map[string]string{},
		current: current,
		recent:  map[string]uint64{},
	}, nil
}

func (t *CursorTracker) Path() string {
	return t.path
}

// Queued records a block written to the local one-block store under file name file.
func (t *CursorTracker) Queued(file string, block *pbbstream.Block) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tracked := &trackedBlock{
		ref:    CursorBlockRef{Num: block.Number, ID: block.Id},
		libNum: block.LibNum,
	}
	if value, ok := t.values[block.Id]; ok {
		tracked.value = &value
		delete(t.values, block.Id)
	}

	t.queue = append(t.queue, tracked)
	t.byFile[file] = append(t.byFile[file], tracked)
}

// Uploaded records that file is in the one-block store. Files the tracker didn't queue (left
// in the local store by a previous run) are ignored.
func (t *CursorTracker) Uploaded(file string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tracked, ok := t.byFile[file]
	if !ok {
		return
	}
	for _, block := range tracked {
		block.uploaded = true
	}
	delete(t.byFile, file)

	advanced := false
	for len(t.queue) > 0 && t.queue[0].uploaded {
		head := t.queue[0]
		t.queue = t.queue[1:]

		next := &ReaderCursor{Block: head.ref, LibNum: head.libNum}
		if t.current != nil {
			next.Value = t.current.Value
		}
		if head.value != nil {
			next.Value = &CursorValue{Block: head.ref, Value: *head.value}
		}
		t.current = next
		t.remember(head.ref)
		advanced = true
	}

	if advanced {
		t.write()
	}
}

// SetValue records the reader's own cursor for a block ("FIRE CURSOR" line). It is written
// once that block is stored.
func (t *CursorTracker) SetValue(num uint64, id string, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, tracked := range t.queue {
		if tracked.ref.ID == id {
			tracked.value = &value
			return
		}
	}

	// already stored: it becomes the cursor's value unless a later block's value is there
	if storedNum, ok := t.recent[id]; ok && storedNum == num {
		if t.current != nil && (t.current.Value == nil || t.current.Value.Block.Num <= num) {
			t.current.Value = &CursorValue{Block: CursorBlockRef{Num: num, ID: id}, Value: value}
			t.write()
		}
		return
	}

	t.values[id] = value
}

func (t *CursorTracker) remember(ref CursorBlockRef) {
	if len(t.recent) >= recentStoredBlocks {
		for id, num := range t.recent {
			if num+recentStoredBlocks < ref.Num {
				delete(t.recent, id)
			}
		}
		if len(t.recent) >= 2*recentStoredBlocks {
			t.recent = map[string]uint64{}
		}
	}
	t.recent[ref.ID] = ref.Num
}

// write replaces the cursor file atomically and syncs it, so a crash leaves either the
// previous cursor or this one.
func (t *CursorTracker) write() {
	if err := writeFileSynced(t.path, t.current); err != nil {
		t.logger.Error("writing reader cursor", zap.String("path", t.path), zap.Error(err))
	}
}

func writeFileSynced(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}

	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
