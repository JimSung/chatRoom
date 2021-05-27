package apps

var (
	chatServerMgr *ChatServerMgr
	chatMgr       *ChatMgr
	cmdMgr        *CmdMgr
	filterMgr     *FilterMgr
)

func OnInit() {
	chatMgr = NewChatMgr()
	cmdMgr = NewCmdMgr()
	filterMgr = NewFilterMgr()
}
