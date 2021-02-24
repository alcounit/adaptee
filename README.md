![GitHub release (latest by date)](https://img.shields.io/github/v/release/alcounit/adaptee)
![Docker Pulls](https://img.shields.io/docker/pulls/alcounit/adaptee)
![GitHub](https://img.shields.io/github/license/alcounit/adaptee)
# adaptee
Selenoid-ui adaptor for selenosis

### Available flags
```
[user@host]$ ./adaptee --help
adaptee is a adaptor sidecar for selenoid ui

Usage:
  adaptee [flags]

Flags:
      --port string                          adaptee port (default ":4444")
      --selenosis-url string                 selenosis url (default "http://selenosis:4444")
      --graceful-shutdown-timeout duration   time in seconds  gracefull shutdown timeout (default 30s)
  -h, --help                                 help for adaptee
```

### Available endpoints
| Protocol | Endpoint                    |
|--------- |---------------------------- |
| HTTP  | /status           |
| WS    | /vnc/{sessionId}  |
| WS    | /logs/{sessionId} |
