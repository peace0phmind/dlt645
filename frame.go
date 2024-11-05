package dlt645

type Code byte

func NewCode(cc CCode) Code {
	return Code(cc.Val())
}

func (c Code) IsResponse() bool {
	return c&(1<<7) != 0
}

func (c Code) Error() bool {
	return c&(1<<6) != 0
}

func (c Code) IsRequest() bool {
	return c&(1<<5) != 0
}

type Frame struct {
	Start      byte    `value:"0x68"` // 帧起始符
	Address    [6]byte // 地址域
	FrameStart byte    `value:"0x68"` // 帧起始符
	C          Code    // 控制码
	L          byte    // 数据域长度
	Data       []byte  // 数据域
	CS         byte    // 校验码
	End        byte    `value:"0x16"` // 帧结束符
}
