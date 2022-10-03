wow-server
==========

Example server for shares quotes from "Words of Wisdom" (www.wow4u.com) with proof-of-work ddos protection and client 
for interact with it.

How to run
==========
With go compiler:

```shell
# Server:
go run ./cmd/server

# Client:
go run ./cmd/client
```

With docker:

```shell
docker build -t wow-server .

# server:
docker run -it --rm -p 8081:8081 --entrypoint server wow-server

# client:
docker run -it --rm --entrypoint client wow-server --server 172.17.0.1:8081
```

Protocol
========
Protocol is based on tcp with request-response interaction model.

```
Client                Server
1. send request

          request
           (GET)
       ------------->
                      
                      2. generate token
                       
       ddos protection
         challenge
  (token + Complexity level)
       <-------------

3. Solve challenge       
       
     challenge solve resp
          (nonce)
       -------------->
           
                      4. Prepare payload
       
       payload
       <-------------
       
5. End
```  
  

Proof of work algorithm
=======================

For ddos protection chose sha256-based HashCash algorithm, because it is simple and widely covered on the Internet.

To ensure the uniqueness of the token, the following information is embedded in it:

- Addr of client (ip + port)
- Current timestamp in millis
- Random integer value
