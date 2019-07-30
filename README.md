# GoCacheable
[![Build Status](https://travis-ci.org/josemiguelmelo/gocacheable.svg?branch=master)](https://travis-ci.org/josemiguelmelo/gocacheable)
[![Coverage Status](https://coveralls.io/repos/github/josemiguelmelo/gocacheable/badge.svg)](https://coveralls.io/github/josemiguelmelo/gocacheable)
[![Go Report Card](https://goreportcard.com/badge/github.com/josemiguelmelo/gocacheable)](https://goreportcard.com/report/github.com/josemiguelmelo/gocacheable)


Golang function calls result cache made easy.

## Install

To use this utility on your project, get it by running the following command:

```
go get github.com/josemiguelmelo/gocacheable
```


## Providers

Currently, the providers available are:

1) BigCache (https://github.com/allegro/bigcache)

For more information about providers, [click here](docs/providers)

## How to use

Using cacheable is quite simple: 

### 1. Create Cacheable Manager

First of all, you must create a Cacheable Manager. This can be done by calling a single method:

```
managerIdentifier := "manager_id"
cacheableManager := gocacheable.NewCacheableManager(managerIdentifier)
```

### 2. Add new module to the manager

A module is a feature related group, with a cache provider associated. 
In order to use gocacheable it is required to have at least one module.

To add a new module, you only need to add the following line to the code:

```
storageProvider := &bcProvider.BigCacheProvider{}
timeToLive := 2
err := cacheableManager.AddModule(moduleName, storageProvider, timeToLive)
```

### 3. Cache function result

After having the cacheable manager and a module, to cache a function return you just need to call it inside Cacheable method.

Suppose you want to cache the function: *func example() string*

```
var outValue string
cacheKey := "test_object"
err := cacheableManager.Cacheable(
    moduleName, 
    cacheKey, 
    func() (interface{}, error) {
        return example(), nil
    },
    &outValue
)
```

When using the code above, if the call is already cached, it returns the cached value. Otherwise, it call the function, caches the result and returns the result into **&outValue**.
