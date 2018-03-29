package notify

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/fsnotify"
)

// EventType event type
type EventType int

// ErrEventType err placehold
const ErrEventType EventType = -1

const (
	// EventCreate event of create
	EventCreate EventType = iota
	// EventModify event of modify
	EventModify
	// EventDelete event of remove
	EventDelete
)

const signalSplit = "|"

// FileWatcher watch file
type FileWatcher struct {
	ChangedFile chan string
	Dirs        []string
	watcher     *fsnotify.Watcher
	closeChan   chan bool
}

// NewFileWatcher create a FileWatcher without any args
func NewFileWatcher() *FileWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("new file watcher failed")
	}

	fw := &FileWatcher{
		watcher:     watcher,
		ChangedFile: make(chan string, 10),
		closeChan:   make(chan bool),
	}

	fw.StartWatch()
	return fw
}

// StartWatch start watch dirs
func (fw *FileWatcher) StartWatch() error {
	watcher := fw.watcher

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if err := fw.sendSignalByNofifyEvent(event); err != nil {
					log.Fatalf("message send failed: %v", err)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			case <-fw.closeChan:
				log.Println("close file watcher....")
				close(fw.ChangedFile)
				close(fw.closeChan)
				goto END_FOR
			}
		}
	END_FOR:
		return
	}()

	return nil
}

// Close close the watch
func (fw *FileWatcher) Close() {
	fw.closeChan <- true
	fw.watcher.Close()
}

// sendSignalByNofifyEvent send signal to channel
func (fw *FileWatcher) sendSignalByNofifyEvent(event *fsnotify.FileEvent) error {
	var header = ErrEventType
	if event.IsModify() {
		header = EventModify
	} else if event.IsCreate() {
		header = EventCreate
	} else if event.IsDelete() {
		header = EventDelete
	}

	if header == ErrEventType {
		return errors.New("unexpect notify event type")
	}

	signal := fw.encodeSignal(header, event.Name)

	fw.ChangedFile <- signal
	return nil
}

// encodeSignal encode signal sended
func (fw *FileWatcher) encodeSignal(header EventType, filename string) string {
	return fmt.Sprintf("%d%s%s", int(header), signalSplit, filename)
}

// DecodeSignal decode notify signal
func DecodeSignal(signal string) (header EventType, filename string, err error) {
	arr := strings.Split(signal, signalSplit)
	if len(arr) != 2 {
		return ErrEventType, "", errors.New("invalid signal")
	}
	et, err := strconv.Atoi(arr[0])
	if err != nil {
		return ErrEventType, "", errors.New("invalid signal header")
	}

	t := EventType(et)
	switch t {
	case EventCreate, EventDelete, EventModify:
		return t, arr[1], nil
	default:
		return ErrEventType, "", errors.New("invalid signal header")
	}
}

// AddDir add a new dir
func (fw *FileWatcher) AddDir(dir string) error {
	for _, d := range fw.Dirs {
		if dir == d {
			return nil
		}
	}
	return fw.addDir(dir)
}

// AddDirs add dirs
func (fw *FileWatcher) AddDirs(dirs []string) error {
	for _, dir := range dirs {
		var has bool
		for _, d := range fw.Dirs {
			if dir == d {
				has = true
			}
		}
		if !has {
			fw.addDir(dir)
		}
	}
	return nil
}

// addDir check, watch and add the dir
func (fw *FileWatcher) addDir(dir string) error {
	// check
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return errors.New("dir is not exist")
	}
	if !info.IsDir() {
		return errors.New("not a directory")
	}

	// watch dir
	if err := fw.watcher.Watch(dir); err != nil {
		return fmt.Errorf("watch dir: %s failed, err: %v", dir, err)
	}
	fw.Dirs = append(fw.Dirs, dir)
	return nil
}
