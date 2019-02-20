package main

const (
	ConfTemplate =
`
{
  "chainId": "0x00",
  "boot": true,
	"chain":{
		"blockprize": "0x1158e460913d00000"
	},
  "default_rep": "10",
	"log": {
		"logLevel": 3
	},
  "logConfig": {
    "logLevel": 3
  },
  "p2p": {
  	"port": 55555,
    "bootNodes":  [
    ]
  },
  "consensus": {
    "consensusMode":"bft"
  }
}
`
)
