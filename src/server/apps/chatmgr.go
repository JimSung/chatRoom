package apps

import "container/list"

type ChatMgr struct {
	UserMap             map[string]*ClientInfo
	ChatList            *list.List
	HotWardsList        map[string]int
	HotWardsRefreshTime int64
}

func NewChatMgr() *ChatMgr {
	return &ChatMgr{
		UserMap:      make(map[string]*ClientInfo),
		ChatList:     list.New(),
		HotWardsList: make(map[string]int),
	}
}

func (c *ChatMgr) Add(chat *Message) {
	if c.ChatList.Len() >= 100 {
		e := c.ChatList.Front()
		if e != nil {
			c.ChatList.Remove(e)
		}
	}
	c.ChatList.PushBack(chat)
	c.HotWardsList[chat.Content]++
}

func (c *ChatMgr) GetHistory() []*Message {
	var allChat []*Message
	for e := c.ChatList.Front(); e != nil; e = e.Next() {
		c, ok := e.Value.(*Message)
		if !ok {
			continue
		}
		allChat = append(allChat, c)
	}
	return allChat
}
