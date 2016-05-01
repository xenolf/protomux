A simple PoC for multiplexing multiple protocols over a single port

```go
listener, err := protomux.Listen(":443")
if err != nil {
  return err
}

for {
  var conn protomux.Conn
  conn, err = listener.Accept()
  if err != nil {
    return err
  }
  
  switch conn.Protocol() {
  case protomux.HTTP:
    // do something
  case protomux.TLS:
    // too something encrypted
  case protomux.None:
    // forward TCP
  }
}
```
