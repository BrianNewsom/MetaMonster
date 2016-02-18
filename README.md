# MetaMonster
![Build button](https://travis-ci.org/BrianNewsom/MetaMonster.svg?branch=master)

A Golang project allowing retrieval of HTML metadata of content driven sites.

## Example
![Example Usage && Output](http://i.imgur.com/p4b5fIC.png)

## Usage
```sh
  metamonster "https://medium.com/@sarah_k_mock/meat-is-dead-long-live-meat-a86a7cfe7ecf" 
```

## Using MetaMonster in my Go project
```go
import "github.com/briannewsom/metamonster/infofetcher"

func main(){
  url := "https://medium.com/@sarah_k_mock/meat-is-dead-long-live-meat-a86a7cfe7ecf"
  
  metadata := infofetcher.GetInfoForUrl(url)

  infofetcher.PrintMetadata(*metadata)
}
```
