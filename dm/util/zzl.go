/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package util

import (
	"container/list"
	"io"
)

type Dm_build_834 struct {
	dm_build_835 *list.List
	dm_build_836 *dm_build_888
	dm_build_837 int
}

func Dm_build_838() *Dm_build_834 {
	return &Dm_build_834{
		dm_build_835: list.New(),
		dm_build_837: 0,
	}
}

func (dm_build_840 *Dm_build_834) Dm_build_839() int {
	return dm_build_840.dm_build_837
}

func (dm_build_842 *Dm_build_834) Dm_build_841(dm_build_843 *Dm_build_912, dm_build_844 int) int {
	var dm_build_845 = 0
	var dm_build_846 = 0
	for dm_build_845 < dm_build_844 && dm_build_842.dm_build_836 != nil {
		dm_build_846 = dm_build_842.dm_build_836.dm_build_896(dm_build_843, dm_build_844-dm_build_845)
		if dm_build_842.dm_build_836.dm_build_891 == 0 {
			dm_build_842.dm_build_878()
		}
		dm_build_845 += dm_build_846
		dm_build_842.dm_build_837 -= dm_build_846
	}
	return dm_build_845
}

func (dm_build_848 *Dm_build_834) Dm_build_847(dm_build_849 []byte, dm_build_850 int, dm_build_851 int) int {
	var dm_build_852 = 0
	var dm_build_853 = 0
	for dm_build_852 < dm_build_851 && dm_build_848.dm_build_836 != nil {
		dm_build_853 = dm_build_848.dm_build_836.dm_build_900(dm_build_849, dm_build_850, dm_build_851-dm_build_852)
		if dm_build_848.dm_build_836.dm_build_891 == 0 {
			dm_build_848.dm_build_878()
		}
		dm_build_852 += dm_build_853
		dm_build_848.dm_build_837 -= dm_build_853
		dm_build_850 += dm_build_853
	}
	return dm_build_852
}

func (dm_build_855 *Dm_build_834) Dm_build_854(dm_build_856 io.Writer, dm_build_857 int) int {
	var dm_build_858 = 0
	var dm_build_859 = 0
	for dm_build_858 < dm_build_857 && dm_build_855.dm_build_836 != nil {
		dm_build_859 = dm_build_855.dm_build_836.dm_build_905(dm_build_856, dm_build_857-dm_build_858)
		if dm_build_855.dm_build_836.dm_build_891 == 0 {
			dm_build_855.dm_build_878()
		}
		dm_build_858 += dm_build_859
		dm_build_855.dm_build_837 -= dm_build_859
	}
	return dm_build_858
}

func (dm_build_861 *Dm_build_834) Dm_build_860(dm_build_862 []byte, dm_build_863 int, dm_build_864 int) {
	if dm_build_864 == 0 {
		return
	}
	var dm_build_865 = dm_build_892(dm_build_862, dm_build_863, dm_build_864)
	if dm_build_861.dm_build_836 == nil {
		dm_build_861.dm_build_836 = dm_build_865
	} else {
		dm_build_861.dm_build_835.PushBack(dm_build_865)
	}
	dm_build_861.dm_build_837 += dm_build_864
}

func (dm_build_867 *Dm_build_834) dm_build_866(dm_build_868 int) byte {
	var dm_build_869 = dm_build_868
	var dm_build_870 = dm_build_867.dm_build_836
	for dm_build_869 > 0 && dm_build_870 != nil {
		if dm_build_870.dm_build_891 == 0 {
			continue
		}
		if dm_build_869 > dm_build_870.dm_build_891-1 {
			dm_build_869 -= dm_build_870.dm_build_891
			dm_build_870 = dm_build_867.dm_build_835.Front().Value.(*dm_build_888)
		} else {
			break
		}
	}
	return dm_build_870.dm_build_909(dm_build_869)
}
func (dm_build_872 *Dm_build_834) Dm_build_871(dm_build_873 *Dm_build_834) {
	if dm_build_873.dm_build_837 == 0 {
		return
	}
	var dm_build_874 = dm_build_873.dm_build_836
	for dm_build_874 != nil {
		dm_build_872.dm_build_875(dm_build_874)
		dm_build_873.dm_build_878()
		dm_build_874 = dm_build_873.dm_build_836
	}
	dm_build_873.dm_build_837 = 0
}
func (dm_build_876 *Dm_build_834) dm_build_875(dm_build_877 *dm_build_888) {
	if dm_build_877.dm_build_891 == 0 {
		return
	}
	if dm_build_876.dm_build_836 == nil {
		dm_build_876.dm_build_836 = dm_build_877
	} else {
		dm_build_876.dm_build_835.PushBack(dm_build_877)
	}
	dm_build_876.dm_build_837 += dm_build_877.dm_build_891
}

func (dm_build_879 *Dm_build_834) dm_build_878() {
	var dm_build_880 = dm_build_879.dm_build_835.Front()
	if dm_build_880 == nil {
		dm_build_879.dm_build_836 = nil
	} else {
		dm_build_879.dm_build_836 = dm_build_880.Value.(*dm_build_888)
		dm_build_879.dm_build_835.Remove(dm_build_880)
	}
}

func (dm_build_882 *Dm_build_834) Dm_build_881() []byte {
	var dm_build_883 = make([]byte, dm_build_882.dm_build_837)
	var dm_build_884 = dm_build_882.dm_build_836
	var dm_build_885 = 0
	var dm_build_886 = len(dm_build_883)
	var dm_build_887 = 0
	for dm_build_884 != nil {
		if dm_build_884.dm_build_891 > 0 {
			if dm_build_886 > dm_build_884.dm_build_891 {
				dm_build_887 = dm_build_884.dm_build_891
			} else {
				dm_build_887 = dm_build_886
			}
			copy(dm_build_883[dm_build_885:dm_build_885+dm_build_887], dm_build_884.dm_build_889[dm_build_884.dm_build_890:dm_build_884.dm_build_890+dm_build_887])
			dm_build_885 += dm_build_887
			dm_build_886 -= dm_build_887
		}
		dm_build_884 = dm_build_882.dm_build_835.Front().Value.(*dm_build_888)
	}
	return dm_build_883
}

type dm_build_888 struct {
	dm_build_889 []byte
	dm_build_890 int
	dm_build_891 int
}

func dm_build_892(dm_build_893 []byte, dm_build_894 int, dm_build_895 int) *dm_build_888 {
	return &dm_build_888{
		dm_build_893,
		dm_build_894,
		dm_build_895,
	}
}

func (dm_build_897 *dm_build_888) dm_build_896(dm_build_898 *Dm_build_912, dm_build_899 int) int {
	if dm_build_897.dm_build_891 <= dm_build_899 {
		dm_build_899 = dm_build_897.dm_build_891
	}
	dm_build_898.Dm_build_1044(dm_build_897.dm_build_889[dm_build_897.dm_build_890 : dm_build_897.dm_build_890+dm_build_899])
	dm_build_897.dm_build_890 += dm_build_899
	dm_build_897.dm_build_891 -= dm_build_899
	return dm_build_899
}

func (dm_build_901 *dm_build_888) dm_build_900(dm_build_902 []byte, dm_build_903 int, dm_build_904 int) int {
	if dm_build_901.dm_build_891 <= dm_build_904 {
		dm_build_904 = dm_build_901.dm_build_891
	}
	copy(dm_build_902[dm_build_903:dm_build_903+dm_build_904], dm_build_901.dm_build_889[dm_build_901.dm_build_890:dm_build_901.dm_build_890+dm_build_904])
	dm_build_901.dm_build_890 += dm_build_904
	dm_build_901.dm_build_891 -= dm_build_904
	return dm_build_904
}

func (dm_build_906 *dm_build_888) dm_build_905(dm_build_907 io.Writer, dm_build_908 int) int {
	if dm_build_906.dm_build_891 <= dm_build_908 {
		dm_build_908 = dm_build_906.dm_build_891
	}
	dm_build_907.Write(dm_build_906.dm_build_889[dm_build_906.dm_build_890 : dm_build_906.dm_build_890+dm_build_908])
	dm_build_906.dm_build_890 += dm_build_908
	dm_build_906.dm_build_891 -= dm_build_908
	return dm_build_908
}
func (dm_build_910 *dm_build_888) dm_build_909(dm_build_911 int) byte {
	return dm_build_910.dm_build_889[dm_build_910.dm_build_890+dm_build_911]
}
