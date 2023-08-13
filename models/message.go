package models

type Message struct {
	MessageId  string // 由UserId和 TargetId 生成
	UserId     string // 发送者
	TargetId   string // 接收者
	Type       int    // 发送类型: 群聊,私聊,广播等
	Media      int    // 文字,图片,音频
	Content    string // 消息内容
	CreateTime uint64 // 创建时间
	ReadTime   uint64 // 读取时间
	Pic        string // 图片
	Url        string // 链接
	Desc       string // 描述
	Amount     int    // 其它统计
}
