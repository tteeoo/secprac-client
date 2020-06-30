# ![secprac-name.png](https://hosted.theohenson.com/secprac-name.png)
![Go](https://github.com/blueberry-jam/secprac-client/workflows/Go/badge.svg) ![license](https://img.shields.io/github/license/blueberry-jam/secprac-client)

`secprac` is a platform to create cyber security practice games for Linux systems similar to the <a href="https://www.uscyberpatriot.org/">Cyber Patriot competition</a>.

This is the client to run on practice machines.

See the web server at <a href="https://github.com/blueberry-jam/secprac-web">blueberry-jam/secprac-web</a>.

Tested on Ubuntu 18.04, Ubuntu 16.04, and Arch Linux.

## Installation

Easily install via scipt with:

```bash
# curl https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/install.sh | sh
```

This script downloads the `secprac-client` binary, icons, and makes the needed directories.

Otherwise, a binary is provided on the latest release page, or build from source with `go build`.

## Usage

```bash
# secprac-client <user> <server url>
```

e.g. `sudo secprac-client john http://192.168.0.3`

The user is the user that notifications will be sent to and that scripts may target. In other words, it is the user that the game is supposed to be played from.

