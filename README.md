Smart home mobile gate
-----------------------

[Project site](https://e154.github.io/smart-home/) |
[Server](https://github.com/e154/smart-home/) |
[Configurator](https://github.com/e154/smart-home-configurator/) |
[Node](https://github.com/e154/smart-home-node/) |
[Smart home Socket](https://github.com/e154/smart-home-socket/) |
[Modbus device controller](https://github.com/e154/smart-home-modbus-ctrl-v1/) |
[Mobile app](https://github.com/e154/smart-home-app/)

![status](https://img.shields.io/badge/status-beta-yellow.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

|Branch      |Status   |
|------------|---------|
|master      | [![Build Status](https://travis-ci.org/e154/smart-home-gate.svg?branch=master)](https://travis-ci.org/e154/smart-home-gate?branch=master)   |
|dev         | [![Build Status](https://travis-ci.org/e154/smart-home-gate.svg?branch=develop)](https://travis-ci.org/e154/smart-home-gate?branch=develop) |


Attention! The project is under active development.
---------

Gateway for connecting mobile clients to the server **smart home**. Controls access rights in accordance with the selected policy.
Determining whether a client belongs to a server occurs through an access token issued to the client.
It is possible to connect a large number of servers and client applications to the server at the same time.

<img height="auto" src="doc/smart-home-gate-network.svg" alt="smart home gate network">

<img height="auto" src="doc/screenshot1.png" alt="smart home gate">

### Installation for development

```bash
git clone https://github.com/e154/smart-home-gate $GOPATH/src/github.com/e154/smart-home-gate

cd $GOPATH/src/github.com/e154/smart-home-gate

go mod vendor

go build
```

editing configuration files

```bash
cp conf/config.dev.json conf/config.json

```

run server

```bash
./smart-home-gate
```

It's all

### LICENSE

[GPLv3 Public License](https://github.com/e154/smart-home-gate/blob/master/LICENSE)
