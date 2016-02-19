# MetaMonster
![Build button](https://travis-ci.org/BrianNewsom/MetaMonster.svg?branch=master)

A Golang project allowing retrieval of HTML metadata of content driven sites.

## Example
![Example Usage && Output](http://i.imgur.com/lQeB0eV.png)

## Installation
First, [install golang if you haven't already](https://golang.org/doc/install).  Then run 
```sh
  go install github.com/briannewsom/metamonster
```

metamonster will now live at ```$GOPATH/bin/metamonster```.  I recommend adding ```$GOPATH/bin``` to your PATH for simplicity.

## Usage
```sh
  metamonster -h
  Usage of metamonster:
  	-format="plaintext": Output data format. Options - [json,plaintext]
	-url="": URL from which to retrieve metadata
```

Or, as an example
```sh
  metamonster -url="https://medium.com/@sarah_k_mock/meat-is-dead-long-live-meat-a86a7cfe7ecf" -format=json
```

## Using MetaMonster in my Go project
```go
import (
	"github.com/briannewsom/metamonster/fetcher"
	"github.com/briannewsom/metamonster/models/metadata"
)

func main(){
  url := "https://medium.com/@sarah_k_mock/meat-is-dead-long-live-meat-a86a7cfe7ecf"

  m, err := fetcher.GetInfoForUrl(url)

  metadata.PrintMetadata(*m)
}
```

## Development

### Testing
```
	go test ./...
```
