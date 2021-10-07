# GLS package tracker

Small lib and example app for GLS package status tracking (using the Hungarian GLS site).
It uses the web backend's json endpoint.

Usage:

```
$ go run example/example.go -pkg <package number>
```

It'll return the current status and the change history, but feel free to extend with any data from the struct.


The example-raw contains a PoC, original implementation.