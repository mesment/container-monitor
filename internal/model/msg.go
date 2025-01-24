package model

type MsgReq struct {
	Subject string     `json:"subject"`  //主题
	MsgType string     `json:"msg_type"` // 消息类型 email sms
	Content MsgContent `json:"content"`  // 消息内容
	To      string     `json:"to"`       // 接受人
}

type MsgContent struct {
	Text string `json:"text"`
}
