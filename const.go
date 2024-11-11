package dlt645

import (
	"encoding/binary"
	"fmt"
)

// @EnumConfig(noCamel)
//go:generate ag --package-mode=true

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
C dlt 645的C码的D4-D0的编码

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
type C byte

func (c C) Value(protocol P) byte {
	if protocol == PV2007 {
		if c.Val() == 0xFF {
			panic(fmt.Errorf("2007 not support control code: %s", c.Name()))
		}
		return c.Val()
	} else {
		if c.Old() == 0xFF {
			panic(fmt.Errorf("1997 not support control code: %s", c.Name()))
		}
		return c.Old()
	}
}

/*
ErrorCode

	@Enum(msg string) {
		RATE ("费率数超") 		= 0x40 // 费率数超
		DAY  ("日时段数超") 		= 0x20 // 日时段数超
		YEAR ("年时区数超") 		= 0x10 // 年时区数超
		BR   ("通信速率不能更改") 	= 0x08 // 通信速率不能更改
		PD   ("密码错误/未授权") 	= 0x04 // 密码错误/未授权
		DATA ("无请求数据") 		= 0x02 // 无请求数据
		OTHER("其他错误") 		= 0x01 // 其他错误
	}
*/
type ErrorCode byte

/*
DIC data identification code. the old is 1997 code, the val is 2007 code

	@EnumConfig(noCase, Values)
	@Enum(old uint16, oldFormat string, oldSize int, newFormat string, newSize int, unit string) {
		// 电能量数据标识
		TotalActiveEnergy             (0xFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00000000 // 组合有功总电能
		PositiveTotalActiveEnergy     (0xFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00010000 // 正向有功总电能
		NegativeTotalActiveEnergy     (0xFFFF, "", 0, "XXXXXX.XX", 4, "kWh")	= 0x00020000 // 反向有功总电能
		TotalReactiveEnergy1          (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00030000 // 组合无功1总电能
		TotalReactiveEnergy2          (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00040000 // 组合无功2总电能
		FirstQuadrantReactiveEnergy   (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00050000 // 第一象限无功电能
		SecondQuadrantReactiveEnergy  (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00060000 // 第二象限无功电能
		ThirdQuadrantReactiveEnergy   (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00070000 // 第三象限无功电能
		FourthQuadrantReactiveEnergy  (0xFFFF, "", 0, "XXXXXX.XX", 4, "kvarh")	= 0x00080000 // 第四象限无功电能
		PositiveTotalApparentEnergy   (0xFFFF, "", 0, "XXXXXX.XX", 4, "KVAh")	= 0x00090000 // 正向视在总电能
		NegativeTotalApparentEnergy   (0xFFFF, "", 0, "XXXXXX.XX", 4, "KVAh")	= 0x000A0000 // 反向视在总电能
		AssociatedTotalElectricEnergy (0xFFFF, "", 0, "XXXXXX.XX", 4, "KVh")	= 0x00800000 // 关联总电能

		// 变量数据标识
		PhaseAVoltage 		(0xB611, "XXX", 2, "XXX.X", 2, "V")			= 0x02010100 // A相电压
		PhaseBVoltage 		(0xB612, "XXX", 2, "XXX.X", 2, "V")			= 0x02010200 // B相电压
		PhaseCVoltage 		(0xB613, "XXX", 2, "XXX.X", 2, "V")			= 0x02010300 // C相电压
		Voltage       		(0xFFFF, "", 0, "XXX.X", 2, "V")			= 0x0201FF00 // 电压数据块
		PhaseACurrent 		(0xB621, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020100 // A相电流
		PhaseBCurrent 		(0xB622, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020200 // B相电流
		PhaseCCurrent 		(0xB623, "XX.XX", 2, "XXX.XXX", 3, "A")		= 0x02020300 // C相电流
		Current       		(0xFFFF, "", 0, "XXX.XXX", 3, "A")			= 0x0202FF00 // 电流数据块
		TotalActivePower  	(0xB630, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030000 // 总有功功率
		PhaseAActivePower 	(0xB631, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030100 // A相有功功率
		PhaseBActivePower 	(0xB632, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030200 // B相有功功率
		PhaseCActivePower 	(0xB633, "XX.XXXX", 3, "XX.XXXX", 3, "kW")	= 0x02030300 // C相有功功率
		ActivePower       	(0xFFFF, "", 0, "XX.XXXX", 3, "kW")			= 0x0203FF00 // 有功功率数据块
		TotalReactivePower  (0xB640, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040000 // 总无功功率
		PhaseAReactivePower (0xB641, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040100 // A相无功功率
		PhaseBReactivePower (0xB642, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040200 // B相无功功率
		PhaseCReactivePower (0xB643, "", 0, "XX.XXXX", 3, "kvar")		= 0x02040300 // C相无功功率
		ReactivePower       (0xFFFF, "", 0, "XX.XXXX", 3, "kvar")		= 0x0204FF00 // 无功功率数据块
		TotalApparentPower  (0xB660, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050000 // 总视在功率
		PhaseAApparentPower (0xB661, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050100 // A相视在功率
		PhaseBApparentPower (0xB662, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050200 // B相视在功率
		PhaseCApparentPower (0xB663, "", 0, "XX.XXXX", 3, "kVA")		= 0x02050300 // C相视在功率
		ApparentPower       (0xFFFF, "", 0, "XX.XXXX", 3, "kVA")		= 0x0205FF00 // 视在功率数据块
		TotalPowerFactor  	(0xFFFF, "", 0, "X.XXX", 2, "")				= 0x02060000 // 总功率因素
		PhaseAPowerFactor 	(0xFFFF, "", 0, "X.XXX", 2, "")				= 0x02060100 // A相功率因素
		PhaseBPowerFactor 	(0xFFFF, "", 0, "X.XXX", 2, "")				= 0x02060200 // B相功率因素
		PhaseCPowerFactor 	(0xFFFF, "", 0, "X.XXX", 2, "")				= 0x02060300 // C相功率因素
		PowerFactor       	(0xFFFF, "", 0, "X.XXX", 2, "")				= 0x0206FF00 // 功率因素数据块
		ABLineVoltage 		(0xB691, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0100 // AB线电压
		BCLineVoltage 		(0xB692, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0200 // BC线电压
		CALineVoltage 		(0xB693, "XXX", 2, "XXX.X", 2, "V")			= 0x020C0300 // CA线电压
		LineVoltage   		(0xFFFF, "", 0, "XXX.X", 2, "V")			= 0x020CFF00 // 线电压数据块
		Frequency 			(0xFFFF, "", 0, "XX.XX", 2, "Hz")			= 0x02800002 // 频率

		// 事件记录数据标识
		TotalOverCurrentCount   (0xFFFF, "", 0, "XXXXXX, XXXXXX", 6, "次,分")	= 0x030C0000 // 过流总次数，总时间
		TotalMeterResetCount 	(0xFFFF, "", 0, "XXXXXX", 3, "次")				= 0x03300100 // 电表清零总次数
		MeterResetRecord     	(0xFFFF, "", 0, "", 0, "")						= 0x03300101 // 电表清零记录, 这个返回的是一个对象的结构体

		// 参变量数据标识
		DateTime            (0xFFFF, "", 0, "YYMMDDWW", 4, "年月日星期")  = 0x04000101 // 年月日星期
		Time                (0xFFFF, "", 0, "hhmmss", 3, "时分秒")		= 0x04000102 // 时分秒
		AssetManagementCode (0xFFFF, "", 0, "N", 32, "")				= 0x04000403 // 资产管理编码
		ActiveConstant		(0xFFFF, "", 0, "XXXXXX", 3, "imp/kWh")		= 0x04000409 // 电表有功常数
		ReactiveConstant	(0xFFFF, "", 0, "XXXXXX", 3, "imp/kvarh")	= 0x0400040A // 电表无功常数
	}
*/
type DIC uint32

func (dic DIC) Code(protocol P) (ret []byte) {
	if protocol == PV2007 {
		ret = binary.LittleEndian.AppendUint32(ret, dic.Val())
	} else {
		if dic.OldSize() == 0 {
			panic(fmt.Errorf("1997 unsupport %s code", dic.Name()))
		}
		ret = binary.LittleEndian.AppendUint16(ret, dic.Old())
	}

	return ret
}

func (dic DIC) Format(protocol P) string {
	if protocol == PV2007 {
		return dic.NewFormat()
	} else {
		if dic.OldSize() == 0 {
			panic(fmt.Errorf("1997 unsupport %s format", dic.Name()))
		}
		return dic.OldFormat()
	}
}

func (dic DIC) Size(protocol P) int {
	if protocol == PV2007 {
		return dic.NewSize()
	} else {
		if dic.OldSize() == 0 {
			panic(fmt.Errorf("1997 unsupport %s size", dic.Name()))
		}
		return dic.OldSize()
	}
}

func getDICs(dic DIC, bitSize int) (ret []DIC) {
	prefix := dic.Val() >> bitSize
	for _, v := range DICValues() {
		if v.Val()>>bitSize == prefix && v != dic {
			ret = append(ret, v)
		}
	}

	return ret
}

// if dic code is block, then return true, else false.
func (dic DIC) CheckBlock(protocol P) (isBlock bool, ret []DIC) {
	if protocol == PV2007 {
		if (dic.Val() & 0xFF) == 0xFF {
			return true, getDICs(dic, 8)
		}

		if (dic.Val() >> 8 & 0xFF) == 0xFF {
			return true, getDICs(dic, 16)
		}

		if (dic.Val() >> 16 & 0xFF) == 0xFF {
			return true, getDICs(dic, 24)
		}

		return false, append([]DIC{}, dic)
	} else {
		return false, append([]DIC{}, dic)
	}
}
