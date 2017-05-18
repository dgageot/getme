# Getme - Downloading dependencies with a cache

https://travis-ci.org/dgageot/getme.svg?branch=master

## Build status

[![Build Status](https://travis-ci.org/dgageot/getme.png?branch=master)](https://travis-ci.org/dgageot/getme)

## Build

```
go build
```

## Usage

```
./getme Download https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip
./getme Copy https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp/docker.zip
./getme Extract https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip /tmp
./getme Extract https://test.docker.com/builds/Windows/x86_64/docker-17.05.0-ce-rc1.zip docker/docker.exe /tmp/docker-windows.exe
```
