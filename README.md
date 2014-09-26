go-fswatch
===

[![GoDoc](https://godoc.org/github.com/joshrotenberg/go-fswatch?status.png)](https://godoc.org/github.com/joshrotenberg/go-fswatch)


Poll a file or directory for changes in Go.

Halt! Unless you explicitly need a polling fs watching solution, you should
probably look at [fsnotify](https://github.com/go-fsnotify/fsnotify) ... and
that may [eventually](https://github.com/go-fsnotify/fsnotify/issues/9) support
polling, making this library obsolete.

This library came about when I realized I needed to support the watching of
similar directory structures that, in one case, reside on the local disk, and
in other cases, reside on an NFS mounted filesystem. There are at least a few
other polling filesystem watching libraries in Go, but fsnotify has a nice API
and fit well with my use case otherwise, so I decided to be a copycat and make
one with the (almost) exact same API.
