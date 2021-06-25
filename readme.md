# Minimal GOPROXY  (MPC)

This is a minimal (extendable) go module proxy protocol implement.

# How to use

1. define some `Resolver`
2. define some `CheckSumResolver` or use `CheckSumResolverNotSupportInstance`
3. use the sample code below:

```go

package main

import (
	"github.com/ZenLiuCN/mpc"
	"net/http"
)

type handler int

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mpc.GoProxyHandler(writer, request)
}

func main() {
	if err := mpc.RegisterResolver("localCacheResolver",0, LocalCacheResolverFactory)
	err != nil{
		panic(err)
	}
	if err := mpc.RegisterCheckSumResolver(0, mpc.CheckSumResolverNotSupportInstance); err != nil {
		panic(err)
	}
	mpc.InitialHandler("/")
	http.ListenAndServe(":80", handler(0))
}

```

# Licence

`AGPL v3`