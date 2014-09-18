# Twitter-A-Go-Go

Minimal Twitter command line client written in Go. Currently, the only function of this script is to retrieve a user's latest tweet when their Twitter screen name is specified.

## Usage

Request the latest Tweet for user 'iheartgo', using configuration file `conf.json` to store Twitter API credentials:

```bash
go run main.go -u iheartgo -c conf.json
```

Arguments:

- `-u`: Twitter screen name of the desired (non-protected) user. Default: Jwpe
- `-c`: Path to configuration file from script location, e.g. Default: `config.json`
