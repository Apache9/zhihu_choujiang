package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
)

type Member struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Data struct {
	Member *Member `json:"member"`
	Action string  `json:"action_type"`
}

type Paging struct {
	IsEnd bool `json:"is_end"`
}

type JsonResp struct {
	Paging *Paging `json:"paging"`
	Data   []*Data `json:"data"`
}

const limit = 100

func main() {
	id := os.Args[1]
	var members []*Member
	for offset := 0; ; offset += limit {
		url := fmt.Sprintf("https://api.zhihu.com/pins/%s/actions?limit=%d&offset=%d", id, limit, offset)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			log.Fatalf("unexpected status code: %d", resp.StatusCode)
		}
		var jsonResp JsonResp
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&jsonResp); err != nil {
			log.Fatal(err)
		}
		if len(jsonResp.Data) > 0 {
			for _, data := range jsonResp.Data {
				log.Printf("name: %s, action_type %s", data.Member.Name, data.Action)
				if data.Action == "repin" {
					members = append(members, data.Member)
				}
			}
		}
		if jsonResp.Paging.IsEnd {
			break
		}
		_ = resp.Body.Close()
	}
	count := len(members)
	log.Printf("total %d users repined", count)
	if index, err := rand.Int(rand.Reader, big.NewInt(int64(count))); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("the chosen one is: %#v", members[index.Int64()])
	}
}
