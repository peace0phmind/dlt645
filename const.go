package dlt645

// @EnumConfig(noCamel)
//go:generate ag

const (
	MaxReadLen             = 200 // 读数据的最大数据长度
	MaxWriteLen            = 50  // 写数据的最大数据长度
	DefaultResponseTimeout = 500 // 500ms
	MaxDeviceNameLen       = 10  // 最大设备名长度
)

/*
P is the protocol version， cmdLen is read data command len

	@Enum(cmdLen int){
		V1997(14)
		V2007(16)
	}
*/
type P int

/*
CCode dlt 645的C码的D4-D0的编码

	@Enum(old byte) {
		BRC(0x08) = 0x08 // 广播校时
		RD (0x01) = 0X11 // 读数据
		RDM(0x02) = 0x12 // 读后续数据
		RDA(0xFF) = 0x13 // 读通信地址
		WR (0x04) = 0x14 // 写数据
		WRA(0x0A) = 0x15 // 写通信地址
		DJ (0xFF) = 0x16 // 冻结命令
		BR (0x0C) = 0x17 // 更改通信速率
		PD (0x0F) = 0x18 // 修改密码
		XL (0x10) = 0x19 // 最大需量清零
		DB (0xFF) = 0x1A // 电表清零
		MSG(0xFF) = 0x1B // 事件清零
		RR (0x03) = 0xFF // 重读数据
	}
*/
type CCode byte

/*
ErrorCode

	@Enum(old byte) {
		RATE (0x40) = 0x40 // 费率数超
		DAY  (0x20) = 0x20 // 日时段数超
		YEAR (0x10) = 0x10 // 年时区数超
		BR   (0x08) = 0x08 // 通信速率不能更改
		PD   (0x04) = 0x04 // 密码错误/未授权
		DATA (0x02) = 0x02 // 无请求数据
		OTHER(0x01) = 0x01 // 其他错误
	}
*/
type ErrorCode byte

/*
DIC data identification code. the old is 1997 code, the val is 2007 code

	@Enum(old uint, oldFormat string, oldSize int, format string, size int, unit string) {
		// 电能量数据标识
		TotalActiveEnergy             (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00000000 // 组合有功总电能
		PositiveTotalActiveEnergy     (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00010000 // 正向有功总电能
		NegativeTotalActiveEnergy     (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00020000 // 反向有功总电能
		TotalReactiveEnergy1          (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00030000 // 组合无功1总电能
		TotalReactiveEnergy2          (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00040000 // 组合无功2总电能
		FirstQuadrantReactiveEnergy   (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00050000 // 第一象限无功电能
		SecondQuadrantReactiveEnergy  (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00060000 // 第二象限无功电能
		ThirdQuadrantReactiveEnergy   (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00070000 // 第三象限无功电能
		FourthQuadrantReactiveEnergy  (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00080000 // 第四象限无功电能
		PositiveTotalApparentEnergy   (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "KVAh")	= 0x00090000 // 正向视在总电能
		NegativeTotalApparentEnergy   (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "KVAh")	= 0x000A0000 // 反向视在总电能
		AssociatedTotalElectricEnergy (0xFFFFFFFF, "", 0, "XXXXXX.XX", 4, "KVh")	= 0x00800000 // 关联总电能

		// 变量数据标识
		PhaseAVoltage 		(0x0000B611, "XXX", 2, "XXX.X", 2, "V")			= 0x02010100 // A相电压
		PhaseBVoltage 		(0x0000B612, "XXX", 2, "XXX.X", 2, "V")			= 0x02010200 // B相电压
		PhaseCVoltage 		(0x0000B613, "XXX", 2, "XXX.X", 2, "V")			= 0x02010300 // C相电压
		Voltage       		(0xFFFFFFFF, "", 0, "XXX.X", 2, "V")			= 0x0201FF00 // 电压数据块
		PhaseACurrent 		(0x0000B621, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020100 // A相电流
		PhaseBCurrent 		(0x0000B622, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020200 // B相电流
		PhaseCCurrent 		(0x0000B623, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020300 // C相电流
		Current       		(0xFFFFFFFF, "", 0, "XXX.XXX", 3, "A")			= 0x0202FF00 // 电流数据块
		TotalActivePower  	(0x0000B630, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030000 // 总有功功率
		PhaseAActivePower 	(0x0000B631, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030100 // A相有功功率
		PhaseBActivePower 	(0x0000B632, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030200 // B相有功功率
		PhaseCActivePower 	(0x0000B633, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030300 // C相有功功率
		ActivePower       	(0xFFFFFFFF, "", 0, "XX.XXXX", 3, "kW")			= 0x0203FF00 // 有功功率数据块
		TotalReactivePower  (0x0000B640, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040000 // 总无功功率
		PhaseAReactivePower (0x0000B641, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040100 // A相无功功率
		PhaseBReactivePower (0x0000B642, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040200 // B相无功功率
		PhaseCReactivePower (0x0000B643, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040300 // C相无功功率
		ReactivePower       (0xFFFFFFFF, "", 0, "XX.XXXX", 3, "kvar")		= 0x0204FF00 // 无功功率数据块
		TotalApparentPower  (0x0000B660, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050000 // 总视在功率
		PhaseAApparentPower (0x0000B661, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050100 // A相视在功率
		PhaseBApparentPower (0x0000B662, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050200 // B相视在功率
		PhaseCApparentPower (0x0000B663, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050300 // C相视在功率
		ApparentPower       (0xFFFFFFFF, "", 0, "XX.XXXX", 3, "kVA")		= 0x0205FF00 // 视在功率数据块
		TotalPowerFactor  	(0xFFFFFFFF, "", 0, "X.XXX", 2, "")				= 0x02060000 // 总功率因素
		PhaseAPowerFactor 	(0xFFFFFFFF, "", 0, "X.XXX", 2, "")				= 0x02060100 // A相功率因素
		PhaseBPowerFactor 	(0xFFFFFFFF, "", 0, "X.XXX", 2, "")				= 0x02060200 // B相功率因素
		PhaseCPowerFactor 	(0xFFFFFFFF, "", 0, "X.XXX", 2, "")				= 0x02060300 // C相功率因素
		PowerFactor       	(0xFFFFFFFF, "", 0, "X.XXX", 2, "")				= 0x0206FF00 // 功率因素数据块
		ABLineVoltage 		(0x0000B691, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0100 // AB线电压
		BCLineVoltage 		(0x0000B692, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0200 // BC线电压
		CALineVoltage 		(0x0000B693, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0300 // CA线电压
		LineVoltage   		(0xFFFFFFFF, "", 0, "XXX.X", 2, "V")			= 0x020CFF00 // 线电压数据块
		Frequency 			(0xFFFFFFFF, "", 0, "XX.XX", 2, "Hz")			= 0x02800002 // 频率

		// 事件记录数据标识
		TotalOverCurrentCount   (0xFFFFFFFF, "", 0, "XXXXXX, XXXXXX", 6, "次,分")	= 0x030C0000 // 过流总次数，总时间
		TotalMeterResetCount 	(0xFFFFFFFF, "", 0, "XXXXXX", 3, "次")				= 0x03300100 // 电表清零总次数
		MeterResetRecord     	(0xFFFFFFFF, "", 0, "", 0, "")						= 0x03300101 // 电表清零记录, 这个返回的是一个对象的结构体

		// 参变量数据标识
		DateTime            (0xFFFFFFFF, "", 0, "YYMMDDWW", 4, "年月日星期")  = 0x04000101 // 年月日星期
		Time                (0xFFFFFFFF, "", 0, "hhmmss", 3, "时分秒")		= 0x04000102 // 时分秒
		AssetManagementCode (0xFFFFFFFF, "", 0, "N", 32, "")				= 0x04000403 // 资产管理编码
		ActiveConstant		(0xFFFFFFFF, "", 0, "XXXXXX", 3, "imp/kWh")		= 0x04000409 // 电表有功常数
		ReactiveConstant	(0xFFFFFFFF, "", 0, "XXXXXX", 3, "imp/kvarh")	= 0x0400040A // 电表无功常数
	}
*/
type DIC uint
