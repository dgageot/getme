# Getme - Downloading dependencies with a cache

## Build

```
go build
```

## Usage

```
./getme Download https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip
./getme Copy https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp/docker.zip
./getme Unzip https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp
./getme UnzipSingleFile https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip docker/docker.exe /tmp/docker-windows.exe
```

## Pruning

Every file retrieved by `getme` ends up being saved to disk in `~/.getme`. On a CI server, this folder can grow forever.
`getme` provides a way to mitigate this issue by allowing one to "Prune" the cache from files that were not recently
accessed.

Each time a url is accessed from or written to cache, this url is appened to `~/.getme/.recent`. `getme Prune` deletes
every url not present in that file and then deletes the `.recent` file.

On a CI system, this is helpful because, the cache directory would typically be retrieved, used, pruned and then saved.