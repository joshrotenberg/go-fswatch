package watch

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	path, err := ioutil.TempDir("", "watch-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	w, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	w.Add(path)

}

func TestWatcherTick(t *testing.T) {
	path, err := ioutil.TempDir("", "watch-test")
	if err != nil {
		t.Fatal("unable to create temp directory for watcher tick test", err)
	}
	defer os.RemoveAll(path)

	ioutil.WriteFile(path+"/doof", []byte("doof"), 0644)
	ioutil.WriteFile(path+"/bar", []byte("bar"), 0644)

	w, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	w.Add(path)

	time.Sleep(time.Millisecond * 500)

	ioutil.WriteFile(path+"/shoof", []byte("shoof"), 0644)
	createEvent := <-w.Events
	if createEvent.Op&Create != Create {
		t.Fatal("expected create event")
	}
	time.Sleep(time.Millisecond * 500)

	ioutil.WriteFile(path+"/doof", []byte("buff"), 0644)
	writeEvent := <-w.Events
	if writeEvent.Op&Write != Write {
		t.Fatal("expected write event")
	}
	time.Sleep(time.Millisecond * 500)

	os.Chmod(path+"/shoof", 0777)
	chmodEvent := <-w.Events
	if chmodEvent.Op&Chmod != Chmod {
		t.Fatal("expected chmod event")
	}
	time.Sleep(time.Millisecond * 500)

	os.Remove(path + "/bar")
	removeEvent := <-w.Events
	if removeEvent.Op&Remove != Remove {
		t.Fatal("expected remove event")
	}

}

func TestWatcherRemovedDirectory(t *testing.T) {
	path, err := ioutil.TempDir("", "watch-remove-test")
	if err != nil {
		t.Fatal("unale to creat temp directory for watcher removal test", err)
	}

	w, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	w.Add(path)

	os.RemoveAll(path)
	errorEvent := <-w.Errors
	if !strings.Contains(errorEvent.Error(), "no such file or directory") {
		t.Fatal("expected no such file or directory error")
	}
}

func ExampleWatcher() {
	watcher, err := NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watcher.Add("/tmp/foo")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("event: %+v", event)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	<-done
}
