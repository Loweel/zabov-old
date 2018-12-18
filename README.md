# zabov

Tiny replacement for piHole DNS filter

Still Work in progress. 

Idea is to produce a very simple, no-web-interface , IP DNS blocker.

Data must be downloaded from URLs of blacklist mantainers.

Blacklists comes in different kinds.

One is the "singlelist", where we find a single column , full of domains:

<pre>
domain1.com
domain2.com
domain3.com
</pre>

The second is the "doublelist", where there is an IP, usually localhost or 0.0.0.0 and then the domain:

<pre>
127.0.0.1 domain1.com
127.0.0.1 domain2.com
127.0.0.1 domain3.com
</pre>

This is why configuration file has two separated items.

The config file should look like:

<pre>
{
    "zabov": {  
        "port":"5453", 
        "proto":"udp", 
        "ipaddr":"127.0.0.1",
        "upstream":"8.8.8.8:53,1.1.1.1:53"  ,
        "singlefilters":"https://mirror1.malwaredomains.com/files/justdomains, https://tspprs.com/dl/cl1" ,
        "doublefilters":"http://sysctl.org/cameleon/hosts,https://www.malwaredomainlist.com/hostslist/hosts.txt,https://adaway.org/hosts.txt", 
        "blackholeip":"127.0.0.1"
    }

}
</pre>

Where:

-port is the port number. Usually is 53
-proto is the protocol. Choices are "udp", "tcp", "tcp/udp"
-upstream: upstream DNS where to forward the DNS query. Comma separated list of IP:PORT
-singlefilters: comma separated list of download URLs, for blacklists following the "singlefilter" schema.
-doublefilters: comma separated list of download URLs, for blacklists following the "doublefilter" schema.
-blackholeip: IP address to return when the IP is banned. This is because you may want to avoid MX issues, mail loops on localhost, or you have a web server running on localhost




