## About The Project

Renterd-remote is an open source Golang library enabling users to access their renterd node from the Internet in total security.

## Getting Started

### Pre requisites

Make sure you have installed all of the following prerequisites on your development machine:

- Git - [Download & Install Git](https://git-scm.com/downloads). OSX and Linux machines typically have this already installed.
- Golang - [Download & Install GoLang](https://go.dev/doc/install) and the Go Framework.

### Installation and Running example

```console
> go build .
> go run .
```

Now open your browser at `http://localhost:8080`

### Create linux service 

To start renterd-remote with linux systems, follow these steps:
- You must run first the app to setup it and to create the configuration file named .env
```console
> go run .
```

- Issue the command: 
```console
> go build .
```

- Create the service with the following command :
```console
nano /etc/systemd/system/renterd-remote.service
```

- Save the following instructions in the file
```console
[Unit]
Description=Renterd-Remote service
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
WorkingDirectory=/<renterd-remote_folder_location/
ExecStart=/<renterd-remote_folder_location>/renterd-remote

[Install]
WantedBy=multi-user.target
```

- Then active the service :
```console
> systemctl enable renterd-remote
```

- Start the service
```console
> systemctl start  renterd-remote
```