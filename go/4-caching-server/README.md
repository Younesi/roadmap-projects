# Caching Proxy

A solution for the [Caching Proxy](https://roadmap.sh/projects/caching-server) challenge from [roadmap.sh](https://roadmap.sh/).

## How to run
Clone the repository and run the following command:

```bash
git clone https://github.com/https://github.com/younesi/roadmap-projects/.git
cd roadmap-projects/go/4-caching-server
```

Run the following command to build and run the project:

```bash
go build -o caching-server
```

### Start caching proxy server for a host and on any port you want!
```bash
./caching-server --port 3000 --origin http://dummyjson.com
```

Then you will see that requests are being cached, e.g. :
```bash
 4-caching-server % ./caching-server --port 3000 --origin http://dummyjson.com  
Server is running on port : 3000 
Enter 'clear-cache' to clear the cache, or 'quit' to exit: 
 Cache MISS : GET:/ 

 Cache MISS : GET:/products 

 Cache HIT : GET:/ 

 Cache HIT : GET:/products 
 
clear-cache
cache is cleared.

 Cache MISS : GET:/ 

 Cache MISS : GET:/products
```

##
You can the total cache of the proxy server.
```bash
clear-cache
```
