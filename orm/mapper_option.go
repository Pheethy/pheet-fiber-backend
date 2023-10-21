package orm

import "github.com/fatih/structs"

type MapperOption struct {
	autobinding bool
	pkFields    []MapperOptionPkField
}

type MapperOptionPkField struct {
	model     interface{}
	fieldName []string
	faith     *structs.Struct
}

func NewMapperOption() MapperOption {
	return MapperOption{
		autobinding: true,
		pkFields:    make([]MapperOptionPkField, 0),
	}
}

func NewMapperOptionPKField(model interface{}, fieldname []string) MapperOptionPkField {
	return MapperOptionPkField{
		model:     model,
		fieldName: fieldname,
		faith:     structs.New(model),
	}
}

func (m MapperOption) SetDisableBinding() MapperOption {
	m.autobinding = false
	return m
}

func (m MapperOption) SetOverridePKField(fields ...MapperOptionPkField) MapperOption {
	m.pkFields = fields
	return m
}
