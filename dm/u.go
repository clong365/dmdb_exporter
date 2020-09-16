/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

type DmStruct struct {
	TypeData
	m_strctDesc *StructDescriptor // 结构体的描述信息

	m_attribs []TypeData // 各属性值

	m_objCount int // 一个数组项中存在对象类型的个数（class、动态数组)

	m_strCount int // 一个数组项中存在字符串类型的个数

	typeName string

	elements []interface{}
}

func newDmStruct(typeName string, elements []interface{}) *DmStruct {
	ds := new(DmStruct)
	ds.typeName = typeName
	ds.elements = elements
	return ds
}

func (ds *DmStruct) create(dc *DmConnection) (*DmStruct, error) {
	desc, err := newStructDescriptor(ds.typeName, dc)
	if err != nil {
		return nil, err
	}
	return ds.createByStructDescriptor(desc, dc)
}

func newDmStructByTypeData(atData []TypeData, desc *TypeDescriptor) *DmStruct {
	ds := new(DmStruct)
	ds.initTypeData()
	ds.m_strctDesc = newStructDescriptorByTypeDescriptor(desc)
	ds.m_attribs = atData
	return ds
}

func (dest *DmStruct) Scan(src interface{}) error {
	switch src := src.(type) {
	case *DmStruct:
		*dest = *src
		return nil
	default:
		return UNSUPPORTED_SCAN
	}
}

func (ds *DmStruct) getAttribsTypeData() []TypeData {
	return ds.m_attribs
}

func (ds *DmStruct) createByStructDescriptor(desc *StructDescriptor, conn *DmConnection) (*DmStruct, error) {
	ds.initTypeData()

	if nil == desc {
		return nil, ECGO_INVALID_PARAMETER_VALUE.throw()
	}

	ds.m_strctDesc = desc
	if nil == ds.elements {
		ds.m_attribs = make([]TypeData, desc.getSize())
	} else {
		if desc.getSize() != len(ds.elements) && desc.getObjId() != 4 {
			return nil, ECGO_STRUCT_MEM_NOT_MATCH.throw()
		}
		var err error
		ds.m_attribs, err = TypeDataSV.toStruct(ds.elements, ds.m_strctDesc.m_typeDesc)
		if err != nil {
			return nil, err
		}
	}

	return ds, nil
}

func (ds *DmStruct) getSQLTypeName() (string, error) {
	return ds.m_strctDesc.m_typeDesc.getFulName()
}

func (ds *DmStruct) getAttributes() ([]interface{}, error) {
	return TypeDataSV.toJavaArrayByDmStruct(ds)
}

func (ds *DmStruct) checkCol(col int) error {
	if col < 1 || col > len(ds.m_attribs) {
		return ECGO_INVALID_SEQUENCE_NUMBER.throw()
	}
	return nil
}

// 获取指定索引的成员变量值，以TypeData的形式给出，col 1 based
func (ds *DmStruct) getAttrValue(col int) (*TypeData, error) {
	err := ds.checkCol(col)
	if err != nil {
		return nil, err
	}
	return &ds.m_attribs[col-1], nil
}
