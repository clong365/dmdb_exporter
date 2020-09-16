/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package util

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math"
)

type dm_build_585 struct{}

var Dm_build_586 = &dm_build_585{}

func (Dm_build_588 *dm_build_585) Dm_build_587(dm_build_589 []byte, dm_build_590 int, dm_build_591 byte) int {
	dm_build_589[dm_build_590] = dm_build_591
	return 1
}

func (Dm_build_593 *dm_build_585) Dm_build_592(dm_build_594 []byte, dm_build_595 int, dm_build_596 int8) int {
	dm_build_594[dm_build_595] = byte(dm_build_596)
	return 1
}

func (Dm_build_598 *dm_build_585) Dm_build_597(dm_build_599 []byte, dm_build_600 int, dm_build_601 int16) int {
	dm_build_599[dm_build_600] = byte(dm_build_601)
	dm_build_600++
	dm_build_599[dm_build_600] = byte(dm_build_601 >> 8)
	return 2
}

func (Dm_build_603 *dm_build_585) Dm_build_602(dm_build_604 []byte, dm_build_605 int, dm_build_606 int32) int {
	dm_build_604[dm_build_605] = byte(dm_build_606)
	dm_build_605++
	dm_build_604[dm_build_605] = byte(dm_build_606 >> 8)
	dm_build_605++
	dm_build_604[dm_build_605] = byte(dm_build_606 >> 16)
	dm_build_605++
	dm_build_604[dm_build_605] = byte(dm_build_606 >> 24)
	dm_build_605++
	return 4
}

func (Dm_build_608 *dm_build_585) Dm_build_607(dm_build_609 []byte, dm_build_610 int, dm_build_611 int64) int {
	dm_build_609[dm_build_610] = byte(dm_build_611)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 8)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 16)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 24)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 32)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 40)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 48)
	dm_build_610++
	dm_build_609[dm_build_610] = byte(dm_build_611 >> 56)
	return 8
}

func (Dm_build_613 *dm_build_585) Dm_build_612(dm_build_614 []byte, dm_build_615 int, dm_build_616 float32) int {
	return Dm_build_613.Dm_build_632(dm_build_614, dm_build_615, math.Float32bits(dm_build_616))
}

func (Dm_build_618 *dm_build_585) Dm_build_617(dm_build_619 []byte, dm_build_620 int, dm_build_621 float64) int {
	return Dm_build_618.Dm_build_637(dm_build_619, dm_build_620, math.Float64bits(dm_build_621))
}

func (Dm_build_623 *dm_build_585) Dm_build_622(dm_build_624 []byte, dm_build_625 int, dm_build_626 uint8) int {
	dm_build_624[dm_build_625] = byte(dm_build_626)
	return 1
}

func (Dm_build_628 *dm_build_585) Dm_build_627(dm_build_629 []byte, dm_build_630 int, dm_build_631 uint16) int {
	dm_build_629[dm_build_630] = byte(dm_build_631)
	dm_build_630++
	dm_build_629[dm_build_630] = byte(dm_build_631 >> 8)
	return 2
}

func (Dm_build_633 *dm_build_585) Dm_build_632(dm_build_634 []byte, dm_build_635 int, dm_build_636 uint32) int {
	dm_build_634[dm_build_635] = byte(dm_build_636)
	dm_build_635++
	dm_build_634[dm_build_635] = byte(dm_build_636 >> 8)
	dm_build_635++
	dm_build_634[dm_build_635] = byte(dm_build_636 >> 16)
	dm_build_635++
	dm_build_634[dm_build_635] = byte(dm_build_636 >> 24)
	return 3
}

func (Dm_build_638 *dm_build_585) Dm_build_637(dm_build_639 []byte, dm_build_640 int, dm_build_641 uint64) int {
	dm_build_639[dm_build_640] = byte(dm_build_641)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 8)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 16)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 24)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 32)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 40)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 48)
	dm_build_640++
	dm_build_639[dm_build_640] = byte(dm_build_641 >> 56)
	return 3
}

func (Dm_build_643 *dm_build_585) Dm_build_642(dm_build_644 []byte, dm_build_645 int, dm_build_646 []byte, dm_build_647 int, dm_build_648 int) int {
	copy(dm_build_644[dm_build_645:dm_build_645+dm_build_648], dm_build_646[dm_build_647:dm_build_647+dm_build_648])
	return dm_build_648
}

func (Dm_build_650 *dm_build_585) Dm_build_649(dm_build_651 []byte, dm_build_652 int, dm_build_653 []byte, dm_build_654 int, dm_build_655 int) int {
	dm_build_652 += Dm_build_650.Dm_build_632(dm_build_651, dm_build_652, uint32(dm_build_655))
	return 4 + Dm_build_650.Dm_build_642(dm_build_651, dm_build_652, dm_build_653, dm_build_654, dm_build_655)
}

func (Dm_build_657 *dm_build_585) Dm_build_656(dm_build_658 []byte, dm_build_659 int, dm_build_660 []byte, dm_build_661 int, dm_build_662 int) int {
	dm_build_659 += Dm_build_657.Dm_build_627(dm_build_658, dm_build_659, uint16(dm_build_662))
	return 2 + Dm_build_657.Dm_build_642(dm_build_658, dm_build_659, dm_build_660, dm_build_661, dm_build_662)
}

func (Dm_build_664 *dm_build_585) Dm_build_663(dm_build_665 []byte, dm_build_666 int, dm_build_667 string, dm_build_668 string) int {
	dm_build_669 := Dm_build_664.Dm_build_793(dm_build_667, dm_build_668)
	dm_build_666 += Dm_build_664.Dm_build_632(dm_build_665, dm_build_666, uint32(len(dm_build_669)))
	return 4 + Dm_build_664.Dm_build_642(dm_build_665, dm_build_666, dm_build_669, 0, len(dm_build_669))
}

func (Dm_build_671 *dm_build_585) Dm_build_670(dm_build_672 []byte, dm_build_673 int, dm_build_674 string, dm_build_675 string) int {
	dm_build_676 := Dm_build_671.Dm_build_793(dm_build_674, dm_build_675)

	dm_build_673 += Dm_build_671.Dm_build_627(dm_build_672, dm_build_673, uint16(len(dm_build_676)))
	return 2 + Dm_build_671.Dm_build_642(dm_build_672, dm_build_673, dm_build_676, 0, len(dm_build_676))
}

func (Dm_build_678 *dm_build_585) Dm_build_677(dm_build_679 []byte, dm_build_680 int) byte {
	return dm_build_679[dm_build_680]
}

func (Dm_build_682 *dm_build_585) Dm_build_681(dm_build_683 []byte, dm_build_684 int) int16 {
	var dm_build_685 int16
	dm_build_685 = int16(dm_build_683[dm_build_684] & 0xff)
	dm_build_684++
	dm_build_685 |= int16(dm_build_683[dm_build_684]&0xff) << 8
	return dm_build_685
}

func (Dm_build_687 *dm_build_585) Dm_build_686(dm_build_688 []byte, dm_build_689 int) int32 {
	var dm_build_690 int32
	dm_build_690 = int32(dm_build_688[dm_build_689] & 0xff)
	dm_build_689++
	dm_build_690 |= int32(dm_build_688[dm_build_689]&0xff) << 8
	dm_build_689++
	dm_build_690 |= int32(dm_build_688[dm_build_689]&0xff) << 16
	dm_build_689++
	dm_build_690 |= int32(dm_build_688[dm_build_689]&0xff) << 24
	return dm_build_690
}

func (Dm_build_692 *dm_build_585) Dm_build_691(dm_build_693 []byte, dm_build_694 int) int64 {
	var dm_build_695 int64
	dm_build_695 = int64(dm_build_693[dm_build_694] & 0xff)
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 8
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 16
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 24
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 32
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 40
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 48
	dm_build_694++
	dm_build_695 |= int64(dm_build_693[dm_build_694]&0xff) << 56
	return dm_build_695
}

func (Dm_build_697 *dm_build_585) Dm_build_696(dm_build_698 []byte, dm_build_699 int) float32 {
	return math.Float32frombits(Dm_build_697.Dm_build_713(dm_build_698, dm_build_699))
}

func (Dm_build_701 *dm_build_585) Dm_build_700(dm_build_702 []byte, dm_build_703 int) float64 {
	return math.Float64frombits(Dm_build_701.Dm_build_718(dm_build_702, dm_build_703))
}

func (Dm_build_705 *dm_build_585) Dm_build_704(dm_build_706 []byte, dm_build_707 int) uint8 {
	return uint8(dm_build_706[dm_build_707] & 0xff)
}

func (Dm_build_709 *dm_build_585) Dm_build_708(dm_build_710 []byte, dm_build_711 int) uint16 {
	var dm_build_712 uint16
	dm_build_712 = uint16(dm_build_710[dm_build_711] & 0xff)
	dm_build_711++
	dm_build_712 |= uint16(dm_build_710[dm_build_711]&0xff) << 8
	return dm_build_712
}

func (Dm_build_714 *dm_build_585) Dm_build_713(dm_build_715 []byte, dm_build_716 int) uint32 {
	var dm_build_717 uint32
	dm_build_717 = uint32(dm_build_715[dm_build_716] & 0xff)
	dm_build_716++
	dm_build_717 |= uint32(dm_build_715[dm_build_716]&0xff) << 8
	dm_build_716++
	dm_build_717 |= uint32(dm_build_715[dm_build_716]&0xff) << 16
	dm_build_716++
	dm_build_717 |= uint32(dm_build_715[dm_build_716]&0xff) << 24
	return dm_build_717
}

func (Dm_build_719 *dm_build_585) Dm_build_718(dm_build_720 []byte, dm_build_721 int) uint64 {
	var dm_build_722 uint64
	dm_build_722 = uint64(dm_build_720[dm_build_721] & 0xff)
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 8
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 16
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 24
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 32
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 40
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 48
	dm_build_721++
	dm_build_722 |= uint64(dm_build_720[dm_build_721]&0xff) << 56
	return dm_build_722
}

func (Dm_build_724 *dm_build_585) Dm_build_723(dm_build_725 []byte, dm_build_726 int) []byte {
	dm_build_727 := Dm_build_724.Dm_build_713(dm_build_725, dm_build_726)
	dm_build_728 := make([]byte, dm_build_727)
	copy(dm_build_728[:int(dm_build_727)], dm_build_725[dm_build_726+4:dm_build_726+4+int(dm_build_727)])
	return dm_build_728
}

func (Dm_build_730 *dm_build_585) Dm_build_729(dm_build_731 []byte, dm_build_732 int) []byte {
	dm_build_733 := Dm_build_730.Dm_build_708(dm_build_731, dm_build_732)
	dm_build_734 := make([]byte, dm_build_733)
	copy(dm_build_734[:int(dm_build_733)], dm_build_731[dm_build_732+2:dm_build_732+2+int(dm_build_733)])
	return dm_build_734
}

func (Dm_build_736 *dm_build_585) Dm_build_735(dm_build_737 []byte, dm_build_738 int, dm_build_739 int) []byte {
	dm_build_740 := make([]byte, dm_build_739)
	copy(dm_build_740[:dm_build_739], dm_build_737[dm_build_738:dm_build_738+dm_build_739])
	return dm_build_740
}

func (Dm_build_742 *dm_build_585) Dm_build_741(dm_build_743 []byte, dm_build_744 int, dm_build_745 int, dm_build_746 string) string {
	return Dm_build_742.Dm_build_828(dm_build_743[dm_build_744:dm_build_744+dm_build_745], dm_build_746)
}

func (Dm_build_748 *dm_build_585) Dm_build_747(dm_build_749 []byte, dm_build_750 int, dm_build_751 string) string {
	dm_build_752 := Dm_build_748.Dm_build_713(dm_build_749, dm_build_750)
	dm_build_750 += 4
	return Dm_build_748.Dm_build_741(dm_build_749, dm_build_750, int(dm_build_752), dm_build_751)
}

func (Dm_build_754 *dm_build_585) Dm_build_753(dm_build_755 []byte, dm_build_756 int, dm_build_757 string) string {
	dm_build_758 := Dm_build_754.Dm_build_708(dm_build_755, dm_build_756)
	dm_build_756 += 2
	return Dm_build_754.Dm_build_741(dm_build_755, dm_build_756, int(dm_build_758), dm_build_757)
}

func (Dm_build_760 *dm_build_585) Dm_build_759(dm_build_761 byte) []byte {
	return []byte{dm_build_761}
}

func (Dm_build_763 *dm_build_585) Dm_build_762(dm_build_764 int16) []byte {
	return []byte{byte(dm_build_764), byte(dm_build_764 >> 8)}
}

func (Dm_build_766 *dm_build_585) Dm_build_765(dm_build_767 int32) []byte {
	return []byte{byte(dm_build_767), byte(dm_build_767 >> 8), byte(dm_build_767 >> 16), byte(dm_build_767 >> 24)}
}

func (Dm_build_769 *dm_build_585) Dm_build_768(dm_build_770 int64) []byte {
	return []byte{byte(dm_build_770), byte(dm_build_770 >> 8), byte(dm_build_770 >> 16), byte(dm_build_770 >> 24), byte(dm_build_770 >> 32),
		byte(dm_build_770 >> 40), byte(dm_build_770 >> 48), byte(dm_build_770 >> 56)}
}

func (Dm_build_772 *dm_build_585) Dm_build_771(dm_build_773 float32) []byte {
	return Dm_build_772.Dm_build_783(math.Float32bits(dm_build_773))
}

func (Dm_build_775 *dm_build_585) Dm_build_774(dm_build_776 float64) []byte {
	return Dm_build_775.Dm_build_786(math.Float64bits(dm_build_776))
}

func (Dm_build_778 *dm_build_585) Dm_build_777(dm_build_779 uint8) []byte {
	return []byte{byte(dm_build_779)}
}

func (Dm_build_781 *dm_build_585) Dm_build_780(dm_build_782 uint16) []byte {
	return []byte{byte(dm_build_782), byte(dm_build_782 >> 8)}
}

func (Dm_build_784 *dm_build_585) Dm_build_783(dm_build_785 uint32) []byte {
	return []byte{byte(dm_build_785), byte(dm_build_785 >> 8), byte(dm_build_785 >> 16), byte(dm_build_785 >> 24)}
}

func (Dm_build_787 *dm_build_585) Dm_build_786(dm_build_788 uint64) []byte {
	return []byte{byte(dm_build_788), byte(dm_build_788 >> 8), byte(dm_build_788 >> 16), byte(dm_build_788 >> 24), byte(dm_build_788 >> 32), byte(dm_build_788 >> 40), byte(dm_build_788 >> 48), byte(dm_build_788 >> 56)}
}

func (Dm_build_790 *dm_build_585) Dm_build_789(dm_build_791 []byte, dm_build_792 string) []byte {
	if dm_build_792 == "UTF-8" {
		return dm_build_791
	}

	if e := dm_build_832(dm_build_792); e != nil {
		tmp, err := ioutil.ReadAll(
			transform.NewReader(bytes.NewReader(dm_build_791), e.NewEncoder()),
		)
		if err != nil {
			panic("UTF8 To Charset error!")
		}

		return tmp
	}

	panic("Unsupported Charset!")
}

func (Dm_build_794 *dm_build_585) Dm_build_793(dm_build_795 string, dm_build_796 string) []byte {
	return Dm_build_794.Dm_build_789([]byte(dm_build_795), dm_build_796)
}

func (Dm_build_798 *dm_build_585) Dm_build_797(dm_build_799 []byte) byte {
	return Dm_build_798.Dm_build_677(dm_build_799, 0)
}

func (Dm_build_801 *dm_build_585) Dm_build_800(dm_build_802 []byte) int16 {
	return Dm_build_801.Dm_build_681(dm_build_802, 0)
}

func (Dm_build_804 *dm_build_585) Dm_build_803(dm_build_805 []byte) int32 {
	return Dm_build_804.Dm_build_686(dm_build_805, 0)
}

func (Dm_build_807 *dm_build_585) Dm_build_806(dm_build_808 []byte) int64 {
	return Dm_build_807.Dm_build_691(dm_build_808, 0)
}

func (Dm_build_810 *dm_build_585) Dm_build_809(dm_build_811 []byte) float32 {
	return Dm_build_810.Dm_build_696(dm_build_811, 0)
}

func (Dm_build_813 *dm_build_585) Dm_build_812(dm_build_814 []byte) float64 {
	return Dm_build_813.Dm_build_700(dm_build_814, 0)
}

func (Dm_build_816 *dm_build_585) Dm_build_815(dm_build_817 []byte) uint8 {
	return Dm_build_816.Dm_build_704(dm_build_817, 0)
}

func (Dm_build_819 *dm_build_585) Dm_build_818(dm_build_820 []byte) uint16 {
	return Dm_build_819.Dm_build_708(dm_build_820, 0)
}

func (Dm_build_822 *dm_build_585) Dm_build_821(dm_build_823 []byte) uint32 {
	return Dm_build_822.Dm_build_713(dm_build_823, 0)
}

func (Dm_build_825 *dm_build_585) Dm_build_824(dm_build_826 []byte, dm_build_827 string) []byte {
	if dm_build_827 == "UTF-8" {
		return dm_build_826
	}

	if e := dm_build_832(dm_build_827); e != nil {
		tmp, err := ioutil.ReadAll(
			transform.NewReader(bytes.NewReader(dm_build_826), e.NewDecoder()),
		)
		if err != nil {

			panic("Charset To UTF8 error!")
		}

		return tmp
	}

	panic("Unsupported Charset!")
}

func (Dm_build_829 *dm_build_585) Dm_build_828(dm_build_830 []byte, dm_build_831 string) string {
	return string(Dm_build_829.Dm_build_824(dm_build_830, dm_build_831))
}

func dm_build_832(dm_build_833 string) encoding.Encoding {
	if e, err := ianaindex.MIB.Encoding(dm_build_833); err == nil && e != nil {
		return e
	}
	return nil
}
