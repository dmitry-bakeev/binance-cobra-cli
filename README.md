# binance-cobra-cli

This is the solution of a [test](https://whattofarm.notion.site/6e7a3804ee0c43e3a2278a4165f18729) task for the position of Go backend developer

Requirements:

- git

- go 1.17.2 or above

How to run:

1. Clone this repo

```bash
# If use https
git clone https://github.com/dmitry-bakeev/binance-cobra-cli.git

# if use ssh
git clone git@github.com:dmitry-bakebinance-cobra-cli.git
```

2. Install dependences

```bash
go mod download
```

3. Run server

```bash
go run . server
```

4. In another window of terminal run client command

```bash
# Example
go run . rate --pair=ETH-USDT
```
