# GoCacheable

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

### Cache function result

Caching a function result is pretty simple. You just need to call it inside Cacheable method.

Suppose you want to cache the function: *func example() string*

```
var outValue string

err := cacheableManager.Cacheable(
    "moduleName", 
    "test_object", 
    func() (interface{}, error) {
        return example(), nil
    },
    &outValue
)
```

When using the code above, if the call is already cached, it returns the cached value. Otherwise, it call the function, caches the result and returns the result into **&outValue**.
