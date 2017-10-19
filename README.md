# golang-disconnectme
Parse the Disconnect.Me JSON list into golang structs

## About

This parses the list of trackers and advertisers from [https://disconnect.me].  The [JSON file](https://github.com/disconnectme/disconnect-tracking-protection/blob/master/services.json) they provide is excellent but is somewhat difficult to process using statically typed languages such as Go.  This package reads the JSON from Disconnect.Me into something more go-like.  Its very likely you will still want to do more post-processing of the data, but this should get you started.

## Usage

```go
package main

import (
	"fmt"
	"github.com/client9/golang-disconnectme"
)

func main() {
	dm, err := disconnectme.ParseURL()
	if err != nil {
		panic(err)
	}
	for category, vendors := range dm {
		fmt.Printf("Category %q has %d entries\n", category, len(vendors))
	}
}
```

`disconnectme.ParseURL()` will fetch the [latest version](https://raw.githubusercontent.com/disconnectme/disconnect-tracking-protection/master/services.json) from the [Github repo](https://github.com/disconnectme/disconnect-tracking-protection).  If you want to read a file or do something else use `disconnectme.Parse` and pass in a `io.Reader` object.

## Legal and License

* This package is MIT licensed.
* The data from Disconnect.Me is [GPLv3](https://github.com/disconnectme/disconnect-tracking-protection/blob/master/LICENSE)
* I'm sure "Disconnect.Me" is trademarked by Disconnect.Me
* This project is not endorsed or affiliated with Disconnect.Me
* Happy to donate this to Disconnect.Me if they want (or they can just copy it since this code is MIT)

OK!

