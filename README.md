# cohesity-net-tools
Web Based Network tools that are typically executed on a Cohesity node to verify connectivity

## About
Cohesity-net-tools is written in HTML, Javascript, and Go. It's a small web application designed to facilitate network testing from a Cohesity node to other parts of your network.

This web app tries to use native implementations of tools wherever possible so that the requirements are minimal.

## Ping
Count is 3.
Timeout is 5s.

## Port Test
Timeout is 3s.

Great for if you think there's a firewall in the way. 

## DNS Lookup
Does a forward lookup of a name and returns A records or AAAA records based on the node's local resolver. 

There is no timeout here although that's probably a good idea.

Another todo here is to maybe lookup reverses as well or automatically do reverse lookups to see if the IPs match the name.

This section may expand in the future to allow for other record types like CNAME, MX, TXT, PTR SRV and so on.

## SSH
The hostname field asks for a port but the default port, 22, is assumed if no value is provided.

As of version 1.0 of Cohesity-net-tools, we only test password authentication but keybased auth is the next most important feature to be added.

Additionally, the command that's run by default is "/usr/bin/id" and we should probably allow more options in the future.

## Traceroute
This feature will probably take the longest amount of time to 

## Tools that might be worth adding
- NTP. Query an NTP server to see if it's working.
- HTTP(S). Get a specific page, etc. or Login via htaccess.

## Ideas
To submit an idea for consideration, please use your support account and submit an idea.

Enjoy.
