## Minter Mass Check Redeemer

### Setup

```bash
mkdir -p $GOPATH/src/github.com/danil-lashin
cd $GOPATH/src/github.com/danil-lashin
git clone https://github.com/danil-lashin/minter-mass-check-redeemer.git
cd minter-mass-check-redeemer

dep ensure
```


### Running 

```bash
go run main.go {PRIVATE_KEY} {ITERATION} {NODE ADDRESS}
```