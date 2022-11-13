# ChatRoom (Web application)
## About project
This project consists on recreating the web chatroom in a `Server-Client` architecture that can run in a server mode on 
a specified port listening for incoming connections. 

## Features:
- TCP connection between server and multiple clients (relation of 1 to many).
- A username selection to the client.
- Messages sent, must be identified by the time that was sent and the username of who sent the message, 
example : `[2020-01-20 15:48:41][client.name]:[client.message]`
- If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.

## Installation
```bash
git clone https://github.com/ZakirAvrora/OneLab-ChatRoom
```
At the root directory
```bash
docker-compose up
```

The default port of running is `8080`