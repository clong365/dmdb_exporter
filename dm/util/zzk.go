/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package util

import (
	"io"
	"math"
)

type Dm_build_912 struct {
	dm_build_913 bool
	dm_build_914 int
	dm_build_915 *Dm_build_925
	dm_build_916 *Dm_build_925
	dm_build_917 *Dm_build_925
}

func Dm_build_918(dm_build_919 int, dm_build_920 bool) *Dm_build_912 {
	dm_build_921 := new(Dm_build_925).dm_build_931(dm_build_919)
	return &Dm_build_912{dm_build_920, 0, dm_build_921, dm_build_921, dm_build_921}
}

func Dm_build_922(dm_build_923 []byte) *Dm_build_912 {
	dm_build_924 := new(Dm_build_925)
	dm_build_924.dm_build_926 = dm_build_923
	dm_build_924.dm_build_927 = 0
	dm_build_924.dm_build_928 = len(dm_build_923)
	dm_build_924.dm_build_929 = nil
	dm_build_924.dm_build_930 = nil
	return &Dm_build_912{true, 0, dm_build_924, dm_build_924, dm_build_924}
}

type Dm_build_925 struct {
	dm_build_926 []byte
	dm_build_927 int
	dm_build_928 int
	dm_build_929 *Dm_build_925
	dm_build_930 *Dm_build_925
}

func (dm_build_932 *Dm_build_925) dm_build_931(dm_build_933 int) *Dm_build_925 {
	dm_build_932.dm_build_926 = make([]byte, dm_build_933, dm_build_933)
	dm_build_932.dm_build_927 = 0
	dm_build_932.dm_build_928 = 0
	dm_build_932.dm_build_929 = nil
	dm_build_932.dm_build_930 = nil
	return dm_build_932
}

func (dm_build_935 *Dm_build_925) dm_build_934(dm_build_936 int) {
	dm_build_935.dm_build_928 = dm_build_936
	if dm_build_935.dm_build_927 > dm_build_936 {
		dm_build_935.dm_build_927 = dm_build_936
	}
}

func (dm_build_938 *Dm_build_925) dm_build_937(dm_build_939 int) {
	dm_build_938.dm_build_927 = dm_build_939
}

func (dm_build_941 *Dm_build_925) dm_build_940() int {
	return cap(dm_build_941.dm_build_926)
}

func (dm_build_943 *Dm_build_925) dm_build_942(dm_build_944 bool) int {
	if dm_build_944 {
		return dm_build_943.dm_build_928
	}

	return dm_build_943.dm_build_927
}

func (dm_build_946 *Dm_build_925) dm_build_945(dm_build_947 bool) int {
	if dm_build_947 {
		return cap(dm_build_946.dm_build_926) - dm_build_946.dm_build_928
	}

	return dm_build_946.dm_build_928 - dm_build_946.dm_build_927
}

func (dm_build_949 *Dm_build_925) dm_build_948(dm_build_950 int, dm_build_951 bool, dm_build_952 bool) int {
	var dm_build_953 int
	if dm_build_952 {
		dm_build_953 = dm_build_949.dm_build_945(dm_build_951)
	} else {
		dm_build_953 = dm_build_949.dm_build_942(dm_build_951)
	}

	if dm_build_950 > dm_build_953 {
		dm_build_950 = dm_build_953
	}

	if dm_build_952 {
		if dm_build_951 {

			for i := 0; i < dm_build_950; i++ {
				dm_build_949.dm_build_926[dm_build_949.dm_build_928] = 0
				dm_build_949.dm_build_928++
			}
		} else {
			dm_build_949.dm_build_927 += dm_build_950
		}
	} else {
		if dm_build_951 {
			dm_build_949.dm_build_928 -= dm_build_950
		} else {
			dm_build_949.dm_build_927 -= dm_build_950
		}
	}

	return dm_build_950
}

func (dm_build_955 *Dm_build_925) dm_build_954(dm_build_956 io.Reader, dm_build_957 int) int {
	dm_build_958 := dm_build_955.dm_build_945(true)
	if dm_build_957 > dm_build_958 {
		dm_build_957 = dm_build_958
	}

	dm_build_959 := 0

	for dm_build_959 < dm_build_957 {
		ret, err := dm_build_956.Read(dm_build_955.dm_build_926[dm_build_955.dm_build_928 : dm_build_955.dm_build_928+dm_build_957-dm_build_959])
		if ret == 0 && err == io.EOF || err != nil && err != io.EOF {
			panic("NETWORK EOF ERROR")
		}

		dm_build_959 += ret
		dm_build_955.dm_build_928 += ret
	}

	return dm_build_959
}

func (dm_build_961 *Dm_build_925) dm_build_960(dm_build_962 io.Writer, dm_build_963 int, dm_build_964 int) int {
	dm_build_965, dm_build_966 := dm_build_962.Write(dm_build_961.dm_build_926[dm_build_963:dm_build_964])
	if dm_build_966 != nil {
		panic("flush fail")
	}

	return dm_build_965
}

func (dm_build_968 *Dm_build_912) Dm_build_967() int {

	dm_build_969 := 0

	dm_build_970 := dm_build_968.dm_build_916
	for dm_build_970 != nil {
		dm_build_969 += dm_build_970.dm_build_942(true)
		dm_build_970 = dm_build_970.dm_build_930
	}

	return dm_build_969
}

func (dm_build_972 *Dm_build_912) Dm_build_971(dm_build_973 int) *Dm_build_912 {
	dm_build_974, dm_build_975 := dm_build_972.dm_build_1010(dm_build_973)
	dm_build_974.dm_build_934(dm_build_975)

	dm_build_972.dm_build_915 = dm_build_974

	dm_build_974 = dm_build_974.dm_build_930
	for dm_build_974 != nil {
		dm_build_974.dm_build_934(0)
		dm_build_974 = dm_build_974.dm_build_930
	}

	return dm_build_972
}

func (dm_build_977 *Dm_build_912) Dm_build_976(dm_build_978 int) *Dm_build_912 {
	dm_build_979, dm_build_980 := dm_build_977.dm_build_1010(dm_build_978)
	dm_build_979.dm_build_937(dm_build_980)

	dm_build_977.dm_build_915 = dm_build_979

	dm_build_979 = dm_build_977.dm_build_916
	for dm_build_979 != dm_build_977.dm_build_915 {
		dm_build_979.dm_build_937(dm_build_979.dm_build_928)
		dm_build_979 = dm_build_979.dm_build_930
	}

	dm_build_979 = dm_build_977.dm_build_915.dm_build_930
	for dm_build_979 != nil {
		dm_build_979.dm_build_937(0)
		dm_build_979 = dm_build_979.dm_build_930
	}

	return dm_build_977
}

func (dm_build_982 *Dm_build_912) Dm_build_981(dm_build_983 bool) int {
	var dm_build_984 int

	if dm_build_983 {
		dm_build_984 = dm_build_982.dm_build_915.dm_build_928
	} else {
		dm_build_984 = dm_build_982.dm_build_915.dm_build_927
	}

	dm_build_985 := dm_build_982.dm_build_916
	for dm_build_985 != dm_build_982.dm_build_915 {
		dm_build_984 += dm_build_985.dm_build_940()
		dm_build_985 = dm_build_985.dm_build_930
	}

	return dm_build_984
}

func (dm_build_987 *Dm_build_912) Dm_build_986(dm_build_988 bool) int {
	if dm_build_988 {
		if dm_build_987.dm_build_913 {
			return math.MaxInt32
		} else {
			return dm_build_987.dm_build_915.dm_build_945(true)
		}
	}

	dm_build_989 := 0
	dm_build_990 := dm_build_987.dm_build_915
	for dm_build_990 != nil {
		dm_build_989 += dm_build_990.dm_build_945(false)
		dm_build_990 = dm_build_990.dm_build_930
	}

	return dm_build_989
}

func (dm_build_992 *Dm_build_912) Dm_build_991(dm_build_993 int, dm_build_994 bool, dm_build_995 bool) *Dm_build_912 {
	for dm_build_993 > 0 {
		dm_build_993 -= dm_build_992.dm_build_915.dm_build_948(dm_build_993, dm_build_994, dm_build_995)
		if dm_build_993 == 0 {
			break
		}

		if dm_build_995 {
			dm_build_992.dm_build_915 = dm_build_992.dm_build_915.dm_build_930
		} else {
			dm_build_992.dm_build_915 = dm_build_992.dm_build_915.dm_build_929
		}

		if dm_build_992.dm_build_915 == nil {
			panic("index out of range")
		}
	}
	return dm_build_992
}

func (dm_build_997 *Dm_build_912) Dm_build_996(dm_build_998 io.Reader, dm_build_999 int) int {
	dm_build_1000 := 0

	for dm_build_999 > 0 {
		llen := dm_build_997.dm_build_915.dm_build_954(dm_build_998, dm_build_999)
		dm_build_1000 += llen
		dm_build_999 -= llen
		if dm_build_999 == 0 {
			break
		}

		dm_build_997.dm_build_915 = dm_build_997.dm_build_915.dm_build_930
		if dm_build_997.dm_build_915 == nil {
			if dm_build_997.dm_build_913 {
				dm_build_997.dm_build_1006(dm_build_999)
			} else {
				panic("index out of range")
			}
		}
	}
	return dm_build_1000
}

func (dm_build_1002 *Dm_build_912) Dm_build_1001(dm_build_1003 io.Writer, dm_build_1004 bool) *Dm_build_912 {
	dm_build_1005 := dm_build_1002.dm_build_916

	for dm_build_1005 != nil {
		if dm_build_1004 {
			dm_build_1005.dm_build_960(dm_build_1003, 0, dm_build_1005.dm_build_940())
		} else {
			dm_build_1005.dm_build_960(dm_build_1003, dm_build_1005.dm_build_927, dm_build_1005.dm_build_928)
		}
		dm_build_1005 = dm_build_1005.dm_build_930
	}

	return dm_build_1002
}

func (dm_build_1007 *Dm_build_912) dm_build_1006(dm_build_1008 int) *Dm_build_912 {
	dm_build_1009 := 2 * dm_build_1007.dm_build_917.dm_build_940()
	if dm_build_1008 < dm_build_1009 {
		dm_build_1008 = dm_build_1009
	}

	dm_build_1007.dm_build_915 = new(Dm_build_925).dm_build_931(dm_build_1008)
	dm_build_1007.dm_build_915.dm_build_929 = dm_build_1007.dm_build_917
	dm_build_1007.dm_build_917.dm_build_930 = dm_build_1007.dm_build_915
	dm_build_1007.dm_build_917 = dm_build_1007.dm_build_915
	dm_build_1007.dm_build_914++

	return dm_build_1007
}

func (dm_build_1011 *Dm_build_912) dm_build_1010(dm_build_1012 int) (*Dm_build_925, int) {
	dm_build_1013 := dm_build_1011.dm_build_916
	for dm_build_1013 != nil {
		skip := dm_build_1013.dm_build_940()
		if dm_build_1012 < skip {
			break
		}

		dm_build_1012 -= skip
		dm_build_1013 = dm_build_1013.dm_build_930
	}

	if dm_build_1013 == nil {
		panic("index out of range")
	}

	return dm_build_1013, dm_build_1012
}

func (dm_build_1015 *Dm_build_912) Dm_build_1014(dm_build_1016 bool) int {
	if dm_build_1016 {
		return dm_build_1015.Dm_build_1044([]byte{1})
	} else {
		return dm_build_1015.Dm_build_1044([]byte{0})
	}
}

func (dm_build_1018 *Dm_build_912) Dm_build_1017(dm_build_1019 byte) int {
	return dm_build_1018.Dm_build_1044(Dm_build_586.Dm_build_759(dm_build_1019))
}

func (dm_build_1021 *Dm_build_912) Dm_build_1020(dm_build_1022 int16) int {
	return dm_build_1021.Dm_build_1044(Dm_build_586.Dm_build_762(dm_build_1022))
}

func (dm_build_1024 *Dm_build_912) Dm_build_1023(dm_build_1025 int32) int {
	return dm_build_1024.Dm_build_1044(Dm_build_586.Dm_build_765(dm_build_1025))
}

func (dm_build_1027 *Dm_build_912) Dm_build_1026(dm_build_1028 uint8) int {
	return dm_build_1027.Dm_build_1044(Dm_build_586.Dm_build_777(dm_build_1028))
}

func (dm_build_1030 *Dm_build_912) Dm_build_1029(dm_build_1031 uint16) int {
	return dm_build_1030.Dm_build_1044(Dm_build_586.Dm_build_780(dm_build_1031))
}

func (dm_build_1033 *Dm_build_912) Dm_build_1032(dm_build_1034 uint32) int {
	return dm_build_1033.Dm_build_1044(Dm_build_586.Dm_build_783(dm_build_1034))
}

func (dm_build_1036 *Dm_build_912) Dm_build_1035(dm_build_1037 uint64) int {
	return dm_build_1036.Dm_build_1044(Dm_build_586.Dm_build_786(dm_build_1037))
}

func (dm_build_1039 *Dm_build_912) Dm_build_1038(dm_build_1040 float32) int {
	return dm_build_1039.Dm_build_1044(Dm_build_586.Dm_build_783(math.Float32bits(dm_build_1040)))
}

func (dm_build_1042 *Dm_build_912) Dm_build_1041(dm_build_1043 float64) int {
	return dm_build_1042.Dm_build_1044(Dm_build_586.Dm_build_786(math.Float64bits(dm_build_1043)))
}

func (dm_build_1045 *Dm_build_912) Dm_build_1044(dm_build_1046 []byte) int {
	dm_build_1047 := len(dm_build_1046)

	for i := 0; i < dm_build_1047; i++ {
		if dm_build_1045.dm_build_915.dm_build_945(true) > 0 {
			dm_build_1045.dm_build_915.dm_build_926[dm_build_1045.dm_build_915.dm_build_928] = dm_build_1046[i]
			dm_build_1045.dm_build_915.dm_build_928++
			continue
		}

		dm_build_1045.dm_build_915 = dm_build_1045.dm_build_915.dm_build_930
		if dm_build_1045.dm_build_915 != nil {
			i--
			continue
		}

		if !dm_build_1045.dm_build_913 {
			panic("index out of range")
		}

		dm_build_1045.dm_build_1006(dm_build_1047 - i)
		i--
	}

	return dm_build_1047
}

func (dm_build_1049 *Dm_build_912) Dm_build_1048(dm_build_1050 []byte) int {
	return dm_build_1049.Dm_build_1023(int32(len(dm_build_1050))) + dm_build_1049.Dm_build_1044(dm_build_1050)
}

func (dm_build_1052 *Dm_build_912) Dm_build_1051(dm_build_1053 []byte) int {
	return dm_build_1052.Dm_build_1026(uint8(len(dm_build_1053))) + dm_build_1052.Dm_build_1044(dm_build_1053)
}

func (dm_build_1055 *Dm_build_912) Dm_build_1054(dm_build_1056 []byte) int {
	return dm_build_1055.Dm_build_1029(uint16(len(dm_build_1056))) + dm_build_1055.Dm_build_1044(dm_build_1056)
}

func (dm_build_1058 *Dm_build_912) Dm_build_1057(dm_build_1059 []byte) int {
	return dm_build_1058.Dm_build_1044(dm_build_1059) + dm_build_1058.Dm_build_1017(0)
}

func (dm_build_1061 *Dm_build_912) Dm_build_1060(dm_build_1062 string, dm_build_1063 string) int {
	dm_build_1064 := Dm_build_586.Dm_build_793(dm_build_1062, dm_build_1063)
	return dm_build_1061.Dm_build_1048(dm_build_1064)
}

func (dm_build_1066 *Dm_build_912) Dm_build_1065(dm_build_1067 string, dm_build_1068 string) int {
	dm_build_1069 := Dm_build_586.Dm_build_793(dm_build_1067, dm_build_1068)
	return dm_build_1066.Dm_build_1051(dm_build_1069)
}

func (dm_build_1071 *Dm_build_912) Dm_build_1070(dm_build_1072 string, dm_build_1073 string) int {
	dm_build_1074 := Dm_build_586.Dm_build_793(dm_build_1072, dm_build_1073)
	return dm_build_1071.Dm_build_1054(dm_build_1074)
}

func (dm_build_1076 *Dm_build_912) Dm_build_1075(dm_build_1077 string, dm_build_1078 string) int {
	dm_build_1079 := Dm_build_586.Dm_build_793(dm_build_1077, dm_build_1078)
	return dm_build_1076.Dm_build_1057(dm_build_1079)
}

func (dm_build_1081 *Dm_build_912) Dm_build_1080() byte {
	return Dm_build_586.Dm_build_797(dm_build_1081.Dm_build_1098(make([]byte, 1, 1)))
}

func (dm_build_1083 *Dm_build_912) Dm_build_1082() int16 {
	return Dm_build_586.Dm_build_800(dm_build_1083.Dm_build_1098(make([]byte, 2, 2)))
}

func (dm_build_1085 *Dm_build_912) Dm_build_1084() int32 {
	return Dm_build_586.Dm_build_803(dm_build_1085.Dm_build_1098(make([]byte, 4, 4)))
}

func (dm_build_1087 *Dm_build_912) Dm_build_1086() int64 {
	return Dm_build_586.Dm_build_806(dm_build_1087.Dm_build_1098(make([]byte, 8, 8)))
}

func (dm_build_1089 *Dm_build_912) Dm_build_1088() float32 {
	return Dm_build_586.Dm_build_809(dm_build_1089.Dm_build_1098(make([]byte, 4, 4)))
}

func (dm_build_1091 *Dm_build_912) Dm_build_1090() float64 {
	return Dm_build_586.Dm_build_812(dm_build_1091.Dm_build_1098(make([]byte, 8, 8)))
}

func (dm_build_1093 *Dm_build_912) Dm_build_1092() uint8 {
	return Dm_build_586.Dm_build_815(dm_build_1093.Dm_build_1098(make([]byte, 1, 1)))
}

func (dm_build_1095 *Dm_build_912) Dm_build_1094() uint16 {
	return Dm_build_586.Dm_build_818(dm_build_1095.Dm_build_1098(make([]byte, 2, 2)))
}

func (dm_build_1097 *Dm_build_912) Dm_build_1096() uint32 {
	return Dm_build_586.Dm_build_821(dm_build_1097.Dm_build_1098(make([]byte, 4, 4)))
}

func (dm_build_1099 *Dm_build_912) Dm_build_1098(dm_build_1100 []byte) []byte {
	for i := 0; i < len(dm_build_1100); i++ {
		if dm_build_1099.dm_build_915 == nil {

			return dm_build_1100
		}
		if dm_build_1099.dm_build_915.dm_build_945(false) > 0 {
			dm_build_1100[i] = dm_build_1099.dm_build_915.dm_build_926[dm_build_1099.dm_build_915.dm_build_927]
			dm_build_1099.dm_build_915.dm_build_927++
			continue
		}

		dm_build_1099.dm_build_915 = dm_build_1099.dm_build_915.dm_build_930

		i--
	}
	return dm_build_1100
}

func (dm_build_1102 *Dm_build_912) Dm_build_1101() []byte {
	dm_build_1103 := dm_build_1102.Dm_build_1084()
	return dm_build_1102.Dm_build_1098(make([]byte, dm_build_1103, dm_build_1103))
}

func (dm_build_1105 *Dm_build_912) Dm_build_1104() []byte {
	dm_build_1106 := dm_build_1105.Dm_build_1080()
	return dm_build_1105.Dm_build_1098(make([]byte, dm_build_1106, dm_build_1106))
}

func (dm_build_1108 *Dm_build_912) Dm_build_1107() []byte {
	dm_build_1109 := dm_build_1108.Dm_build_1082()
	return dm_build_1108.Dm_build_1098(make([]byte, dm_build_1109, dm_build_1109))
}

func (dm_build_1111 *Dm_build_912) Dm_build_1110(dm_build_1112 int) []byte {
	return dm_build_1111.Dm_build_1098(make([]byte, dm_build_1112, dm_build_1112))
}

func (dm_build_1114 *Dm_build_912) Dm_build_1113() []byte {
	dm_build_1115 := 0
	for dm_build_1114.Dm_build_1080() != 0 {
		dm_build_1115++
	}
	dm_build_1114.Dm_build_991(dm_build_1115, false, false)
	return dm_build_1114.Dm_build_1098(make([]byte, dm_build_1115-1, dm_build_1115-1))
}

func (dm_build_1117 *Dm_build_912) Dm_build_1116(dm_build_1118 int, dm_build_1119 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1117.Dm_build_1098(make([]byte, dm_build_1118, dm_build_1118)), dm_build_1119)
}

func (dm_build_1121 *Dm_build_912) Dm_build_1120(dm_build_1122 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1121.Dm_build_1101(), dm_build_1122)
}

func (dm_build_1124 *Dm_build_912) Dm_build_1123(dm_build_1125 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1124.Dm_build_1104(), dm_build_1125)
}

func (dm_build_1127 *Dm_build_912) Dm_build_1126(dm_build_1128 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1127.Dm_build_1107(), dm_build_1128)
}

func (dm_build_1130 *Dm_build_912) Dm_build_1129(dm_build_1131 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1130.Dm_build_1113(), dm_build_1131)
}

func (dm_build_1133 *Dm_build_912) Dm_build_1132(dm_build_1134 int, dm_build_1135 byte) int {
	return dm_build_1133.Dm_build_1168(dm_build_1134, Dm_build_586.Dm_build_759(dm_build_1135))
}

func (dm_build_1137 *Dm_build_912) Dm_build_1136(dm_build_1138 int, dm_build_1139 int16) int {
	return dm_build_1137.Dm_build_1168(dm_build_1138, Dm_build_586.Dm_build_762(dm_build_1139))
}

func (dm_build_1141 *Dm_build_912) Dm_build_1140(dm_build_1142 int, dm_build_1143 int32) int {
	return dm_build_1141.Dm_build_1168(dm_build_1142, Dm_build_586.Dm_build_765(dm_build_1143))
}

func (dm_build_1145 *Dm_build_912) Dm_build_1144(dm_build_1146 int, dm_build_1147 int64) int {
	return dm_build_1145.Dm_build_1168(dm_build_1146, Dm_build_586.Dm_build_768(dm_build_1147))
}

func (dm_build_1149 *Dm_build_912) Dm_build_1148(dm_build_1150 int, dm_build_1151 float32) int {
	return dm_build_1149.Dm_build_1168(dm_build_1150, Dm_build_586.Dm_build_771(dm_build_1151))
}

func (dm_build_1153 *Dm_build_912) Dm_build_1152(dm_build_1154 int, dm_build_1155 float64) int {
	return dm_build_1153.Dm_build_1168(dm_build_1154, Dm_build_586.Dm_build_774(dm_build_1155))
}

func (dm_build_1157 *Dm_build_912) Dm_build_1156(dm_build_1158 int, dm_build_1159 uint8) int {
	return dm_build_1157.Dm_build_1168(dm_build_1158, Dm_build_586.Dm_build_777(dm_build_1159))
}

func (dm_build_1161 *Dm_build_912) Dm_build_1160(dm_build_1162 int, dm_build_1163 uint16) int {
	return dm_build_1161.Dm_build_1168(dm_build_1162, Dm_build_586.Dm_build_780(dm_build_1163))
}

func (dm_build_1165 *Dm_build_912) Dm_build_1164(dm_build_1166 int, dm_build_1167 uint32) int {
	return dm_build_1165.Dm_build_1168(dm_build_1166, Dm_build_586.Dm_build_783(dm_build_1167))
}

func (dm_build_1169 *Dm_build_912) Dm_build_1168(dm_build_1170 int, dm_build_1171 []byte) int {
	dm_build_1172, dm_build_1173 := dm_build_1169.dm_build_1010(dm_build_1170)
	for i := 0; i < len(dm_build_1171); i++ {
		if dm_build_1173 < dm_build_1172.dm_build_940() {
			dm_build_1172.dm_build_926[dm_build_1173] = dm_build_1171[i]
			dm_build_1173++
			continue
		}

		dm_build_1172 = dm_build_1172.dm_build_930
		if dm_build_1172 == nil {
			panic("index of range")
		}

		i--
		dm_build_1173 = 0
	}

	return len(dm_build_1171)
}

func (dm_build_1175 *Dm_build_912) Dm_build_1174(dm_build_1176 int, dm_build_1177 []byte) int {
	return dm_build_1175.Dm_build_1140(dm_build_1176, int32(len(dm_build_1177))) + dm_build_1175.Dm_build_1168(dm_build_1176+4, dm_build_1177)
}

func (dm_build_1179 *Dm_build_912) Dm_build_1178(dm_build_1180 int, dm_build_1181 []byte) int {
	return dm_build_1179.Dm_build_1132(dm_build_1180, byte(len(dm_build_1181))) + dm_build_1179.Dm_build_1168(dm_build_1180+1, dm_build_1181)
}

func (dm_build_1183 *Dm_build_912) Dm_build_1182(dm_build_1184 int, dm_build_1185 []byte) int {
	return dm_build_1183.Dm_build_1136(dm_build_1184, int16(len(dm_build_1185))) + dm_build_1183.Dm_build_1168(dm_build_1184+2, dm_build_1185)
}

func (dm_build_1187 *Dm_build_912) Dm_build_1186(dm_build_1188 int, dm_build_1189 []byte) int {
	return dm_build_1187.Dm_build_1168(dm_build_1188, dm_build_1189) + dm_build_1187.Dm_build_1132(dm_build_1188+len(dm_build_1189), 0)
}

func (dm_build_1191 *Dm_build_912) Dm_build_1190(dm_build_1192 int, dm_build_1193 string, dm_build_1194 string) int {
	return dm_build_1191.Dm_build_1174(dm_build_1192, Dm_build_586.Dm_build_793(dm_build_1193, dm_build_1194))
}

func (dm_build_1196 *Dm_build_912) Dm_build_1195(dm_build_1197 int, dm_build_1198 string, dm_build_1199 string) int {
	return dm_build_1196.Dm_build_1178(dm_build_1197, Dm_build_586.Dm_build_793(dm_build_1198, dm_build_1199))
}

func (dm_build_1201 *Dm_build_912) Dm_build_1200(dm_build_1202 int, dm_build_1203 string, dm_build_1204 string) int {
	return dm_build_1201.Dm_build_1182(dm_build_1202, Dm_build_586.Dm_build_793(dm_build_1203, dm_build_1204))
}

func (dm_build_1206 *Dm_build_912) Dm_build_1205(dm_build_1207 int, dm_build_1208 string, dm_build_1209 string) int {
	return dm_build_1206.Dm_build_1186(dm_build_1207, Dm_build_586.Dm_build_793(dm_build_1208, dm_build_1209))
}

func (dm_build_1211 *Dm_build_912) Dm_build_1210(dm_build_1212 int) byte {
	return Dm_build_586.Dm_build_797(dm_build_1211.Dm_build_1237(dm_build_1212, make([]byte, 1, 1)))
}

func (dm_build_1214 *Dm_build_912) Dm_build_1213(dm_build_1215 int) int16 {
	return Dm_build_586.Dm_build_800(dm_build_1214.Dm_build_1237(dm_build_1215, make([]byte, 2, 2)))
}

func (dm_build_1217 *Dm_build_912) Dm_build_1216(dm_build_1218 int) int32 {
	return Dm_build_586.Dm_build_803(dm_build_1217.Dm_build_1237(dm_build_1218, make([]byte, 4, 4)))
}

func (dm_build_1220 *Dm_build_912) Dm_build_1219(dm_build_1221 int) int64 {
	return Dm_build_586.Dm_build_806(dm_build_1220.Dm_build_1237(dm_build_1221, make([]byte, 8, 8)))
}

func (dm_build_1223 *Dm_build_912) Dm_build_1222(dm_build_1224 int) float32 {
	return Dm_build_586.Dm_build_809(dm_build_1223.Dm_build_1237(dm_build_1224, make([]byte, 4, 4)))
}

func (dm_build_1226 *Dm_build_912) Dm_build_1225(dm_build_1227 int) float64 {
	return Dm_build_586.Dm_build_812(dm_build_1226.Dm_build_1237(dm_build_1227, make([]byte, 8, 8)))
}

func (dm_build_1229 *Dm_build_912) Dm_build_1228(dm_build_1230 int) uint8 {
	return Dm_build_586.Dm_build_815(dm_build_1229.Dm_build_1237(dm_build_1230, make([]byte, 1, 1)))
}

func (dm_build_1232 *Dm_build_912) Dm_build_1231(dm_build_1233 int) uint16 {
	return Dm_build_586.Dm_build_818(dm_build_1232.Dm_build_1237(dm_build_1233, make([]byte, 2, 2)))
}

func (dm_build_1235 *Dm_build_912) Dm_build_1234(dm_build_1236 int) uint32 {
	return Dm_build_586.Dm_build_821(dm_build_1235.Dm_build_1237(dm_build_1236, make([]byte, 4, 4)))
}

func (dm_build_1238 *Dm_build_912) Dm_build_1237(dm_build_1239 int, dm_build_1240 []byte) []byte {
	dm_build_1241, dm_build_1242 := dm_build_1238.dm_build_1010(dm_build_1239)

	for i := 0; i < len(dm_build_1240); i++ {

		if dm_build_1242 < dm_build_1241.dm_build_940() {
			dm_build_1240[i] = dm_build_1238.dm_build_915.dm_build_926[dm_build_1242]
			dm_build_1242++
			continue
		}

		dm_build_1241 = dm_build_1241.dm_build_930
		if dm_build_1241 == nil {
			panic("index out of range")
		}

		i--
		dm_build_1242 = 0
	}

	return dm_build_1240
}

func (dm_build_1244 *Dm_build_912) Dm_build_1243(dm_build_1245 int) []byte {
	dm_build_1246 := dm_build_1244.Dm_build_1216(dm_build_1245)
	return dm_build_1244.Dm_build_1237(dm_build_1245+4, make([]byte, dm_build_1246, dm_build_1246))
}

func (dm_build_1248 *Dm_build_912) Dm_build_1247(dm_build_1249 int) []byte {
	dm_build_1250 := dm_build_1248.Dm_build_1210(dm_build_1249)
	return dm_build_1248.Dm_build_1237(dm_build_1249+1, make([]byte, dm_build_1250, dm_build_1250))
}

func (dm_build_1252 *Dm_build_912) Dm_build_1251(dm_build_1253 int) []byte {
	dm_build_1254 := dm_build_1252.Dm_build_1213(dm_build_1253)
	return dm_build_1252.Dm_build_1237(dm_build_1253+2, make([]byte, dm_build_1254, dm_build_1254))
}

func (dm_build_1256 *Dm_build_912) Dm_build_1255(dm_build_1257 int) []byte {
	dm_build_1258 := 0
	for dm_build_1256.Dm_build_1210(dm_build_1257) != 0 {
		dm_build_1257++
		dm_build_1258++
	}

	return dm_build_1256.Dm_build_1237(dm_build_1257-dm_build_1258, make([]byte, dm_build_1258, dm_build_1258))
}

func (dm_build_1260 *Dm_build_912) Dm_build_1259(dm_build_1261 int, dm_build_1262 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1260.Dm_build_1243(dm_build_1261), dm_build_1262)
}

func (dm_build_1264 *Dm_build_912) Dm_build_1263(dm_build_1265 int, dm_build_1266 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1264.Dm_build_1247(dm_build_1265), dm_build_1266)
}

func (dm_build_1268 *Dm_build_912) Dm_build_1267(dm_build_1269 int, dm_build_1270 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1268.Dm_build_1251(dm_build_1269), dm_build_1270)
}

func (dm_build_1272 *Dm_build_912) Dm_build_1271(dm_build_1273 int, dm_build_1274 string) string {
	return Dm_build_586.Dm_build_828(dm_build_1272.Dm_build_1255(dm_build_1273), dm_build_1274)
}
