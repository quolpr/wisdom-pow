## Task

Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.
 • TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
 • The choice of the POW algorithm should be explained.
 • After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
 • Docker file should be provided both for the server and for the client that solves the POW challenge

## Running

Client:
```bash
make run-client
```

Server:
```bash
make run-server
```

## POW algorithm

I chose the SHA-256 Hashcash algorithm because it's simple to implement. The client will need to use significant resources to create a hash, but it will be easy and cheap for the server to check if it's correct. The client must find a number that, when used in the hash, starts with N zeros.

