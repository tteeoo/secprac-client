# ![secprac-name.png](https://directory.theohenson.com/file/img/secprac-name.png)
![Go build status](https://github.com/blueberry-jam/secprac-client/workflows/Go/badge.svg) ![License (MIT)](https://img.shields.io/github/license/blueberry-jam/secprac-client)

`secprac` is a platform to create cyber security practice games for Linux systems similar to the [Cyber Patriot competition](https://www.uscyberpatriot.org/).

This is the client to run on practice virtual machines.

See the web server at <a href="https://github.com/blueberry-jam/secprac-web">blueberry-jam/secprac-web</a>.

Tested on Ubuntu 18.04, Ubuntu 16.04, Debian 10, and Arch Linux.

## Installation

Easily install via scipt by running:

```
curl https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/install.sh | sudo sh
```

Note: on some Linux distributions `curl` is not installed by default. For example, Ubuntu. To install curl (which is needed for the command above) on Ubuntu, run

```
sudo apt install curl
```

This script downloads the program, icons, and makes the needed directories.

Alternatively, an executalbe is provided on the latest GitHub release page, or build the client from source with `go build`.

## Usage

```
# secprac-client <user> <server url>
```

e.g. `sudo secprac-client john http://192.168.0.3`

The user is the user that notifications will be sent to and that scripts may target. In other words, it is the user that the game is supposed to be played from.

## Setting up a game (and more extensive documentation)

See the [web server repo's wiki](https://github.com/blueberry-jam/secprac-web/wiki).

## License

All files are licensed under the permissive MIT License.
