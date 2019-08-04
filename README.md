
Simple code to talk to the OpenVPN management ports of multiple OpenVPN<br>
processes. This code will open sockets, send "kill" command and agreggate<br>
the number of disconnected clients.<br>
 <br>
A telnet session to a single OpenVPN process looks like this:<br>
<br>
```
[fkooman@vpn ~]$ telnet localhost 11940
Trying ::1...
telnet: connect to address ::1: Connection refused
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
>INFO:OpenVPN Management Interface Version 1 -- type 'help' for more info
kill 07d1ccc455a21c2d5ac6068d4af727ca
SUCCESS: common name '07d1ccc455a21c2d5ac6068d4af727ca' found, 1 client(s) killed
kill foo
ERROR: common name 'foo' not found
quit
Connection closed by foreign host.
[fkooman@vpn ~]$ 
```
The point here is to be able to (concurrently) connect to many OpenVPN<br> 
processes. The example below has only two. Extra functionality later will<br>
be also the use of the "status" command to see which clients are connected<br>
and aggregate that as well.<br>
<br>
Eventually this will need to become a daemon that supports TLS and abstracts<br>
the multiple OpenVPN processes away from the daemon caller...<br>
<br>