Forwards failed requests to other server lol

```
go build 

./proxy-forwarder PORT scheme://host:port scheme2://alternatehost:port2 code1 ...
./proxy-forwarder.exe # if ur on windows
# example:
./proxy-forwarder.exe 5002 'http://localhost:5000' 'http://localhost:5001' 401 404
```
