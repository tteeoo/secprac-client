# ![secprac-name.png](https://dir.theohenson.com/file/img/secprac-name.png)
![Go build status](https://github.com/tteeoo/secprac-client/workflows/Go/badge.svg) ![License (MIT)](https://img.shields.io/github/license/tteeoo/secprac-client)

`secprac` is a platform to create cyber security practice games for Linux systems, similar to the [CyberPatriot](https://www.uscyberpatriot.org/) competition.

This is the client to run on practice virtual machines.

See the web server at [tteeoo/secprac-web](https://github.com/tteeoo/secprac-web).

## Installation

Easily install via script by running:

```
wget https://raw.githubusercontent.com/tteeoo/secprac-client/master/install.sh && sudo sh install.sh
```

## Usage

To start the client, run:

```
# secprac-start <user> <server url>
```

Replace `<user>` with the username of the main non-root user, and `<server url>` with the URL of the server, including `http://` and without a trailing slash.

E.g. `sudo secprac-start john http://192.168.0.3`.

The above command will use a script that will attempt to start the client as a systemd service.

To run the client directly (not recommended) run:

```
# secprac-client <user> <server url>
```

## Setting up a game (and more extensive documentation)

See the [web server repository's wiki](https://github.com/tteeoo/secprac-web/wiki).

## License

All files are licensed under the MIT License, except for `data/FiraSans-Light.ttf`, whose license is located at `data/LICENSE-Fira`.
