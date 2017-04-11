# Getme - Downloading dependencies with a cache

## Build

```
go build
```

## Usage

```
./getme Download https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp/docker.zip
./getme Unzip https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp
./getme UnzipSingleFile https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip docker/docker.exe /tmp/docker-windows.exe
```