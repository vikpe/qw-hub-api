# QWS

> Web API serving QuakeWorld server info

## Usage

```sh
qws [-master INTERVAL] [-server INTERVAL] [-active INTERVAL] [-port PORT]
```

| arg      | type  | description                   | default | 
|----------|-------|-------------------------------|---------|
| `port`   | `int` | HTTP listen port              | `3000`  |
| `master` | `int` | Master server update interval | `600`   |
| `server` | `int` | Server update interval        | `30`    |
| `active` | `int` | Active server update interval | `3`     |

## API endpoints

| URL                     | description             |  
|-------------------------|-------------------------|
| `/v2/servers/<address>` | Server details          |  
| `/v2/servers/mvdsv`     | Mvdsv servers           |  
| `/v2/servers/qwfwd`     | Qwfwd servers (proxies) |  
| `/v2/servers/qtv`       | QTV servers             |  

### Query params

| URL                 | description                                    |
|---------------------|------------------------------------------------|
| `has_player=xantom` | Servers where xantom is connected as player    |
| `has_client=xantom` | Servers where xantom is connected              |

## Config

### Master servers

The QuakeWorld master servers to query for servers.

**Example**
`master_servers.json`

```json
[
  "master.quakeworld.nu:27000",
  "master.quakeservers.net:27000",
  "qwmaster.ocrana.de:27000",
  "qwmaster.fodquake.net:27000"
]
```

## Build

```sh
$ go build
```

## Credits

* eb
* Tuna
* XantoM

## See also

* [streambot](https://github.com/vikpe/qw-streambot)
* [masterstat](https://github.com/vikpe/masterstat)
* [masterstat-cli](https://github.com/vikpe/masterstat-cli)
* [serverstat](https://github.com/vikpe/serverstat)
* [serverstat-cli](https://github.com/vikpe/serverstat-cli)
