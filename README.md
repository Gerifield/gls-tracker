# GLS package tracker

Small PoC app for GLS package status tracking.
It uses the web backend's json endpoint.

Usage:

```
$ go run main.go -pkg <package number>
```

It'll return the current status and the change history.
