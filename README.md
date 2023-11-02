# Running Server

```
make server
```

this will start the server on port 4000,and will wait for clients to syncronize states among them.

# Running Clients

start clients by running:

```
make client
```

## Example

```
❯ make server
TRAC[0000] [INBOX] started                               pid=local/server
starting HTTP server on port 4000
TRAC[0000] [PROCESS] started                             pid=local/server
server !
```

on shell 1: run `make client`
on shell 2: run `make client`

server:

```
...
new client trying to connect
TRAC[0009] [INBOX] started pid=local/server/session_3552661241410906520
rev: {}
TRAC[0009] [PROCESS] started pid=local/server/session_3552661241410906520
client with sid 3552661241410906520 and pid local/server/session_3552661241410906520 just connected
recieved login msg
new client trying to connect
TRAC[0028] [INBOX] started pid=local/server/session_5777754122544690289
rev: {}
TRAC[0028] [PROCESS] started pid=local/server/session_5777754122544690289
client with sid 5777754122544690289 and pid local/server/session_5777754122544690289 just connected
recieved login msg
sending this state &{100 {365 957} 5777754122544690289}
sending this state &{100 {811 724} 3552661241410906520}
sending this state &{100 {889 512} 5777754122544690289}
sending this state &{100 {512 906} 3552661241410906520}
sending this state &{100 {441 658} 5777754122544690289}
```

client 1:

```
...
need to update the state of player {100 {365 957} 5777754122544690289}
need to update the state of player {100 {889 512} 5777754122544690289}
...
```

client 2:

```
...
need to update the state of player {100 {811 724} 3552661241410906520}
need to update the state of player {100 {512 906} 3552661241410906520}
...
```
