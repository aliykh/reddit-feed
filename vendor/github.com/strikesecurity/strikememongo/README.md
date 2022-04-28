`strikememongo` is a Go package that spins up a real MongoDB server, backed by in-memory
storage, for use in testing and mocking during development. It's based on
[mongodb-memory-server](https://github.com/nodkz/mongodb-memory-server) for
NodeJS and [memongo](https://github.com/benweissmann/memongo) for Go.

In general, it's better to mock out interaction with the database, so you don't
need to run a Mongo server during testing. But because most Mongo clients use
a fluent interface that's tough to mock, and sometimes you need to test the
queries themselves, it's often helpful to be able to spin up a Mongo server
quickly and easily. That's where `strikememongo` comes in!

# Project Status

Beta. Tests and CI are set up and working, but more esoteric configurations may not work. If StrikeMemongo isn't working on your platform, you might want to use `strikememongo.StartWithOptions()` and pass the correct download URL for your platform manually.

# Caveats and Notes

Currently, `strikememongo` only supports UNIX systems. CI will run on MacOS, Ubuntu Xenial, Ubuntu Trusty, and Ubuntu Precise. Other flavors of Linux may or may not work.

# Basic Usage

Spin up a server for a single test:

```go
func TestSomething(t *testing.T) {
  mongoServer, err := strikememongo.Start("4.0.5")
  if (err != nil) {
    t.Fatal(err)
  }
  defer mongoServer.Stop()

  connectAndDoStuff(mongoServer.URI(), strikememongo.RandomDatabase())
}
```

Spin up a server, shared between tests:

```go
var mongoServer strikememongo.Server;

func TestMain(m *testing.M) {
  mongoServer, err = strikememongo.Start("4.0.5")
  if (err != nil) {
    log.Fatal(err)
  }
  defer mongoServer.Stop()

  os.Exit(m.Run())
}

func TestSomething(t *testing.T) {
  connectAndDoStuff(mongoServer.URI(), strikememongo.RandomDatabase())
}
```

Spin up a replica set server:

```go
var mongoServer strikememongo.Server;

func TestMain(m *testing.M) {
  mongoServer, err = strikememongo.StartWitOptions(&strikememongo.Options{MongoVersion: "4.2.1", ShouldUseReplica: true})
  if (err != nil) {
    log.Fatal(err)
  }
  defer mongoServer.Stop()

  os.Exit(m.Run())
}

func TestSomething(t *testing.T) {
  connectAndDoStuff(mongoServer.URI(), strikememongo.RandomDatabase())
}
```

# How it works

Behind the scenes, when you run `Start()`, a few things are happening:

1. If you specified a MongoDB version number (rather than a URL or binary path),
   `strikememongo` detects your operating system and platform to determine the
   download URL for the right MongoDB binary.

2. If you specified a MongoDB version number or download URL, `strikememongo`
   downloads MongoDB to a cache location. For future runs, `strikememongo` will just
   use the copy from the cache. You only need to be connected to the internet
   the first time you run `Start()` for a particular MongoDB version.

3. `strikememongo` starts a process running the downloaded `mongod` binary. It uses
   the `ephemeralForTest` storage engine, a temporary directory for a `dbpath`,
   and a random free port number.

4. `strikememongo` also starts up a "watcher" process. This process is a simple
   portable shell script that kills the `mongod` process when the current
   process exits. This ensures that we don't leave behind `mongod` processes,
   even if your tests exit uncleanly or you don't call `Stop()`.

# Configuration

The behavior of `strikememongo` can be controlled by using
`strikememongo.StartWithOptions` instead of `strikememongo.Start`. See
[the godoc](https://godoc.org/github.com/strikesecurity/strikememongo) for all the options. Many options can also be set via environment variable.

A few common use-cases are covered here:

Note that you must use MongoDB version 3.2 or greater, because the `ephemeralForTest` storage engine was not present before 3.2.

## Set the cache path

`strikememongo` downloads a pre-compiled binary of MongoDB from https://www.mongodb.org and caches it on your local system. This path is set by (in order of preference):

- The `CachePath` passed to `memongo.StartWithOptions`
- The environment variable `MEMONGO_CACHE_PATH`
- If `XDG_CACHE_HOME` is set, `$XDG_CACHE_HOME/memongo`
- `~/.cache/memongo` on Linux, or `~/Library/Caches/memongo` on MacOS

## Override download URL

By default, `strikememongo` tries to detect the platform you're running on and download an official MongoDB release for it. If `strikememongo` doesn't yet support your platform, of you'd like to use a custom version of MongoDB, you can pass `DownloadURL` to `strikememongo.StartWithOptions` or set the environment variable `MEMONGO_DOWNLOAD_URL`.

`strikememongo`'s caching will still work with custom download URLs.

## Use a custom MongoDB binary

If you'd like to bypass `strikememongo`'s download beahvior entirely, you can pass `MongodBin` to `strikememongo.StartWithOptions`, or set the environment variable `MEMONGO_MONGOD_BIN` to the path to a `mongod` binary. `strikememongo` will use this binary instead of downloading one.

If you're running on a platform that doesn't have an official MongoDB release (such as Alpine), you'll need to use this option.

## Reduce or increase logging

By default, `strikememongo` logs at an "info" level. You may call `StartWithOptions` with `LogLevel: strikememongolog.LogLevelWarn` for fewer logs, `LogLevel: strikememongolog.LogLevelSilent` for no logs, or `LogLevel: strikememongolog.LogLevelDebug` for verbose logs (including full logs from MongoDB).

By default, `strikememongo` logs to stdout. To log somewhere else, specify a `Logger` in `StartWithOptions`.
