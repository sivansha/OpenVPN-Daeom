# VPN Daemon

We want to create a simple daemon that can be queried by the Let's Connect!
portal, hereafter called portal.

## Why?

Currently we have many OpenVPN processes to manage. The portal connects to 
every OpenVPN process using the OpenVPN management (TCP) socket. This works 
fine if the OpenVPN processes and the portal run on the same machine. If both 
the portal and OpenVPN processes run on different hosts this is less than ideal 
for security, performance and reliability reasons. 

## How?

What we want to build is a simple daemon that runs on the same node as the 
OpenVPN processes and is reachable over a local file system socket as well as
a TCP socket protected by TLS. The daemon will then take care of contacting 
the OpenVPN processes and execute the commands.

Currently there are two commands used to talk to OpenVPN: `status` and `kill` 
where `status` returns a list of connected clients, and `kill` disconnects a
client.

By default Let's Connect! has two OpenVPN processes, so the daemon will need
to talk to both OpenVPN processes. The portal can just talk to the daemon and
issues a command there. The results will be merged by the daemon. 

In addition: we can create a (much) cleaner API then the one used by OpenVPN 
and abstract the CSV format of the `status` command in something more modern,
e.g. JSON or maybe even protobuf.

## Before

Current situation:

                   .----------------.
                   | Portal         |
          .--------|                |------.
          |        '----------------'      |
          |                                |
          |                                |
          |                                |
          |                                |
          |Local/Remote TCP Socket         |Local/Remote TCP Socket
          |                                |
          v                                v
    .----------------.               .----------------.
    | OpenVPN 1      |               | OpenVPN 2      |
    |                |               |                |
    '----------------'               '----------------'

## After

                  .----------------.
                  | Portal         |
                  |                |
                  '----------------'
                           |
                           | Local/Remote TCP/TLS Socket
                           v
                  .----------------.
                  | Daemon         |
          .-------|                |-------.
          |       '----------------'       |
          |                                |
          |Local Socket                    |Local Socket
          |                                |
          v                                v
    .----------------.               .----------------.
    | OpenVPN 1      |               | OpenVPN 2      |
    |                |               |                |
    '----------------'               '----------------'

## Benefits

The daemon will be written in Go, which can handle connections to OpenVPN
concurrently, it doesn't have to do one after the other thus potentially 
increasing performance.

We can use TLS with a daemon. Go makes this easy to do it securely.

The parsing of the OpenVPN "legacy" protocol and merging of the 
information can be done by the daemon.

We can also begin to envision implementing other VPN protocols when we have
a control daemon, e.g. Wireguard. The daemon would need to have additional 
commands then, i.e. `setup` and `teardown`.

## Steps

1. Create a socket client that can talk to OpenVPN management port
2. Implement `kill`
3. Implement connecting to multiple OpenVPN processes concurrent
4. Implement daemon and listen on TCP socket and handle commands from daemon
5. Aggregate feedback from the various OpenVPN managements ports
6. Implement `status`
7. Implement a way to periodically kill all client connections where the 
   certificate expired, is this even possible without referring to the CA?
8. Allow the portal to talk to the daemon(s) and indicate which clients are
   allowed to connect until what time (@ cert creation either through client
   API or manually through the portal)
   - then there is no need for the node to talk to the portal @connect
9. How to sync the list of clients between portal/nodes when a new node is 
   introduced, or they lose "sync"? 

## Daemon API

### Status

This command will tell the portal when it received its last update from the 
portal. The portal can then "replay" the `ADD` / `REMOVE` commands necessary to 
get the daemon in the correct state again.

    STATUS

Example:

    STATUS

Response:

    2019-08-04T10:27:44+00:00

### Add

`ADD` will allow this CN to connect, until the provided `EXPIRY`. Once `EXPIRY`
is reached, the `REMOVE` command will be run automatically.

    ADD <CN> <EXPIRY>

Example:

    ADD 07d1ccc455a21c2d5ac6068d4af727ca 2019-08-04T10:27:44+00:00

Response:

    OK

### Remove 

`REMOVE` will no longer allow the mentioned CN to connect, and also disconnect 
any clients connected with this CN. This is intended for (admin) initiated 
client removals, e.g. in case of abuse.

    REMOVE <CN>

Example:

    REMOVE 07d1ccc455a21c2d5ac6068d4af727ca
    
Response:

    OK

### List

This will list all currently connected clients.

    LIST

    <CN> <EXPIRY> <IPv4> <IPv6>

Example:

    LIST

Response:

    07d1ccc455a21c2d5ac6068d4af727ca 2019-08-04T10:27:44+00:00 10.42.42.2 fd00:4242:4242:4242::1000
