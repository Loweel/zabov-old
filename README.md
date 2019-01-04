# zabov

Tiny replacement for piHole DNS filter

Still Work in progress, usable.

Idea is to produce a very simple, no-web-interface , IP DNS blocker.

# INSTALL

Zabov requires golang 1.11 or later.

<pre>
git clone https://git.keinpfusch.net/LowEel/zabov.git
cd zabov
go get
go build
</pre>

Then, edit config.json: please notice config.json must be in the same folder of the executable you run.


Just a few words about "singlefilters" and "doublefilters":

Data must be downloaded from URLs of blacklist mantainers.

There are two kinds of blacklists:

One is the "singlefilter", where we find a single column , full of domains:

<pre>
domain1.com
domain2.com
domain3.com
</pre>

The second is the "doublefilter" (a file in "hosts" format, to be precise), where there is an IP, usually localhost or 0.0.0.0 and then the domain:

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
        "port":"53", 
        "proto":"udp", 
        "ipaddr":"127.0.0.1",
        "upstream":"8.8.8.8:53,1.1.1.1:53"  ,
        "cachettl": "4",
        "killfilettl": "12",
        "singlefilters":"https://mirror1.malwaredomains.com/files/justdomains, https://tspprs.com/dl/cl1" ,
        "doublefilters":"http://sysctl.org/cameleon/hosts,https://www.malwaredomainlist.com/hostslist/hosts.txt,https://adaway.org/hosts.txt", 
        "blackholeip":"127.0.0.1",
        "hostsfile":"./hosts.txt"
    }

}
</pre>

Where:

- port is the port number. Usually is 53, you can change for docker, if you like
- proto is the protocol. Choices are "udp", "tcp", "tcp/udp"
- ipaddr is the port to listen to. Maybe empty, (which will result in listening to 0.0.0.0) to avoid issues with docker.
- upstream: upstream DNS where to forward the DNS query. Comma separated list of IP:PORT
- cachettl: amount of time the cache is kept (in hours)
- killfilettl: refresh time for killfiles
- singlefilters: comma separated list of download URLs, for blacklists following the "singlefilter" schema.
- doublefilters: comma separated list of download URLs, for blacklists following the "doublefilter" schema.
- blackholeip: IP address to return when the IP is banned. This is because you may want to avoid MX issues, mail loops on localhost, or you have a web server running on localhost
- hostsfile: path where you keep your local blacklistfile : this is in the format "singlefilter", meaning one domain per line, unlike hosts file.

# TODO:

- ~~caching~~
- monitoring port


