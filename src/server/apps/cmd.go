package apps

import (
	"fmt"
	"time"
)

type CmdMgr struct {
	CmdHandler map[string]func(*ClientInfo) string
}

func NewCmdMgr() *CmdMgr {
	mgr := &CmdMgr{
		CmdHandler: make(map[string]func(*ClientInfo) string),
	}
	mgr.CmdHandler["/popular"] = GetHotWords
	mgr.CmdHandler["/stats"] = GetOnlineTime
	return mgr
}

func GetOnlineTime(user *ClientInfo) string {
	return fmt.Sprintf("online time:%v", time.Now().Sub(user.CreateTime).String())
}

func GetHotWords(user *ClientInfo) string {
	hotWords := ""
	hotWeight := 0
	for words, weight := range chatMgr.HotWardsList {
		if weight > hotWeight {
			hotWords = words
		}
	}
	return fmt.Sprintf("hot words:%v, use num:%v", hotWords, hotWeight)
}
