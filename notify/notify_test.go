package notify

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewFileWatcher(t *testing.T) {
	t.Run("create a watch", func(t *testing.T) {
		w := NewFileWatcher()
		defer w.watcher.Close()

		if w == nil {
			t.Errorf("NewFileWatcher() expect to return a valid pointer but get a nil")
		}

		tp := reflect.TypeOf(w)

		if tp.Kind() != reflect.Ptr {
			t.Errorf("NewFileWatcher() expect to return a pointer")
		}

		if strings.LastIndex(tp.String(), "FileWatcher") == -1 {
			t.Errorf("NewFileWatcher() expect to return a pointer of struct FileWatcher")
		}
	})
}

func TestFileWatcher_DecodeSignal(t *testing.T) {
	tests := []struct {
		name      string
		signal    string
		eventType EventType
		filename  string
		wantErr   bool
		match     bool
	}{
		{"create file", "0|a.txt", EventCreate, "a.txt", false, true},
		{"modify file", "1|b.txt", EventModify, "b.txt", false, true},
		{"delete file", "2|c.txt", EventDelete, "c.txt", false, true},
		{"modify file but type wrong", "1|c.txt", EventDelete, "c.txt", false, false},
		{"modify file but filename wrong", "1|d.txt", EventModify, "c.txt", false, false},
		{"wrong split", "1*c.txt", EventModify, "c.txt", true, false},
		{"wrong signal", "1|c.txt|x", EventModify, "c.txt", true, false},
		{"invalid event type", "7|c.txt", EventModify, "c.txt", true, false},
		{"invalid event type", "ok|c.txt", EventModify, "c.txt", true, false},
	}
	w := NewFileWatcher()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeader, gotFilename, err := w.DecodeSignal(tt.signal)
			if tt.wantErr != (err != nil) {
				t.Errorf("FileWatcher.DecodeSignal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (EventType(gotHeader) != tt.eventType || gotFilename != tt.filename) == tt.match {
				t.Errorf("FileWatcher.DecodeSignal() decode failed")
			}
		})
	}
	w.watcher.Close()
}

func TestFileWatcher_AddDir(t *testing.T) {
	w := NewFileWatcher()
	defer w.watcher.Close()

	t.Run("add dir successful", func(t *testing.T) {
		tp, _ := ioutil.TempDir("", "")
		defer os.RemoveAll(tp)
		err := w.AddDir(tp)
		if err != nil {
			t.Errorf("FileWatcher.AddDir() add dir failed")
		}
		var hasAdded bool
		for _, dir := range w.Dirs {
			if dir == tp {
				hasAdded = true
			}
		}
		if !hasAdded {
			t.Errorf("FileWatcher.AddDir() add dir failed")
		}

		// add the same again
		err = w.AddDir(tp)
		if err != nil {
			t.Errorf("FileWatcher.AddDir() add dir failed")
		}

		if !reflect.DeepEqual([]string{tp}, w.Dirs) {
			t.Errorf("FileWatcher.AddDir() add dir failed")
		}
	})

	t.Run("add dir failed for unexist dir", func(t *testing.T) {
		tp, _ := ioutil.TempDir("", "")
		os.RemoveAll(tp)
		err := w.AddDir(tp)
		if err == nil {
			t.Errorf("FileWatcher.AddDir() expect to get err when add unexist dir")
			return
		}
	})

	t.Run("add dir failed for a exist file dir", func(t *testing.T) {
		tp, _ := ioutil.TempDir("", "")
		fp := filepath.Join(tp, "txt.txt")
		os.Create(fp)
		defer os.RemoveAll(tp)
		err := w.AddDir(fp)
		if err == nil {
			t.Errorf("FileWatcher.AddDir() expect to get err when add a exist file")
			return
		}
	})
}

func TestFileWatcher_AddDirs(t *testing.T) {
	w := NewFileWatcher()
	defer w.watcher.Close()

	t.Run("add dirs successful", func(t *testing.T) {

		tp1, _ := ioutil.TempDir("", "")
		tp2, _ := ioutil.TempDir("", "")
		var dirs = []string{tp1, tp2, tp1}
		defer os.RemoveAll(tp1)
		defer os.RemoveAll(tp2)
		err := w.AddDirs(dirs)

		if err != nil {
			t.Errorf("FileWatcher.AddDirs() add dirs failed 1")
		}

		if !reflect.DeepEqual([]string{tp1, tp2}, w.Dirs) {
			t.Errorf("FileWatcher.AddDir() add dirs failed 2")
		}
	})
}

func TestFileWatcher_encodeSignal(t *testing.T) {
	w := NewFileWatcher()
	defer w.watcher.Close()

	tests := []struct {
		name     string
		header   EventType
		filename string
		want     string
	}{
		{"create success", EventCreate, "text.txt", "0|text.txt"},
		{"modify success", EventModify, "text.txt", "1|text.txt"},
		{"delete success", EventDelete, "text.txt", "2|text.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := w.encodeSignal(tt.header, tt.filename); got != tt.want {
				t.Errorf("FileWatcher.encodeSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}
