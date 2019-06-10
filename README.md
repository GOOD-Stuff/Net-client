# Net-client
Utility for work with UDP/TCP packets, rewritten from old C++ project.

## What it could?
* send UDP/TCP packets (repeatedly and once);
* receive UDP/TCP packets;
* send byte/ASCII data;
* read data from file (and send it).


## How to launch
```
go build -o net_client main.go
./net_client -h # show help page
./net_client -ip <ip address> -pt <port of connection> -t <type of connection: udp/tcp>
```
