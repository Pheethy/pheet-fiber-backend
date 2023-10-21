package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/BlackMocca/sqlx"
	"github.com/fatih/structs"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
)

type Mapper struct {
	modelStructs  []modelStruct
	rowCount      int
	paginateTotal int
}

func newMapper(mainModel interface{}, options MapperOption) (Mapper, error) {
	mapper := Mapper{
		modelStructs: make([]modelStruct, 0),
	}
	ms, err := newModelStruct(mainModel, options)
	if err != nil {
		return mapper, err
	}
	if len(ms) > 0 && options.autobinding {
		for index := range ms {
			if ms[index].IsMainModel() {
				continue
			}
			subModels, err := newModelStruct(ms[index].model, options)
			if err != nil {
				return mapper, err
			}
			if len(subModels) > 1 {
				for subIndex := range subModels {
					subModels[subIndex].isReferenceModel = false
				}
				ms[index].subRefModel = subModels
			}
		}
	}

	mapper.modelStructs = ms
	return mapper, nil
}

func (m Mapper) GetData() interface{} {
	return modelStructs(m.modelStructs).GetMainModel().modelSlice.Interface()
}

func (m Mapper) GetRowCount() int {
	return m.rowCount
}

func (m Mapper) GetPaginateTotal() int {
	return m.paginateTotal
}

func validateModel(model interface{}) error {
	if isNil := isNil(model); isNil {
		return ErrMustNotNil
	}
	model = reflect.ValueOf(model).Elem().Interface()
	if err := mustbeStruct(model); err != nil {
		return err
	}

	return nil
}

func recovery() error {
	if r := recover(); r != nil {
		if v, ok := r.(error); ok {
			return v
		} else if reflect.TypeOf(v).Kind() == reflect.String {
			return errors.New(cast.ToString(v))
		}
	}
	return nil
}

func orm(ctx context.Context, model interface{}, rows *sqlx.Rows, options MapperOption) (Mapper, error) {
	if err := validateModel(model); err != nil {
		return Mapper{}, err
	}
	var mapper, err = newMapper(model, options)
	if err != nil {
		return mapper, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return mapper, err
	}

	var paginateColumnIndex = -1
	if len(columns) > 0 {
		for index, col := range columns {
			if strings.EqualFold(col.Name(), PAGINATE_COLUMN_NAME) {
				paginateColumnIndex = index
			}
		}
	}

	var rowCount int
	var paginateTotal int
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return mapper, err
		}
		values, err := rows.SliceScan()
		if err != nil {
			return mapper, err
		}
		rowCount++
		if paginateColumnIndex != -1 {
			paginateTotal = cast.ToInt(values[paginateColumnIndex])
		}

		var group, _ = errgroup.WithContext(ctx)
		var fillData = func(ms *modelStruct) {
			group.Go(func() error {
				defer func() {
					if panicErr := recovery(); panicErr != nil {
						err = panicErr
					}
				}()

				slice, err := fillValueList(ms, columns, values, options)
				if err != nil {
					return err
				}
				ms.modelSlice = slice

				return nil
			})
		}
		for index := range mapper.modelStructs {
			if mapper.modelStructs[index].IsMainModel() {
				fillData(&mapper.modelStructs[index])
				continue
			}
			/* root reference model */
			if !options.autobinding {
				continue
			}
			fillData(&mapper.modelStructs[index])
			if len(mapper.modelStructs[index].subRefModel) > 0 {
				for subModelIndex := range mapper.modelStructs[index].subRefModel {
					if mapper.modelStructs[index].name != mapper.modelStructs[index].subRefModel[subModelIndex].name {
						fillData(&mapper.modelStructs[index].subRefModel[subModelIndex])
					}
				}
			}
		}

		if err := group.Wait(); err != nil {
			return mapper, err
		}
	}

	/* orm relation with pkMainModel Id */
	if len(mapper.modelStructs) > 1 && options.autobinding {
		mainModel := modelStructs(mapper.modelStructs).GetMainModel()
		var bind = func(ctx context.Context, group *errgroup.Group, elem reflect.Value, refFields []string, allmodels []modelStruct) {
			group.Go(func() error {
				defer func() {
					if panicErr := recovery(); panicErr != nil {
						err = panicErr
					}
				}()

				return bindReference(ctx, elem, refFields, allmodels)
			})
		}
		/* orm sub component */
		if refIndexes := modelStructs(mapper.modelStructs).GetListReferenceModelIndex(); len(refIndexes) > 0 && len(options.pkFields) > 0 {
			var refGroup, ctx = errgroup.WithContext(ctx)
			for _, refIndex := range refIndexes {
				refModel := mapper.modelStructs[refIndex]
				if refModel.modelSlice.Len() == 0 || len(refModel.subRefModel) == 0 {
					continue
				}
				func(ctx context.Context, refModel modelStruct, refIndex int) {
					refGroup.Go(func() error {
						var group, subGroupCtx = errgroup.WithContext(ctx)
						subMainModel := modelStructs(refModel.subRefModel).GetMainModel()
						subMainModel.modelSlice = refModel.modelSlice

						for index := 0; index < subMainModel.modelSlice.Len(); index++ {
							bind(subGroupCtx, group, subMainModel.modelSlice.Index(index), subMainModel.refFields, refModel.subRefModel)
						}
						if err := group.Wait(); err != nil {
							return err
						}

						mapper.modelStructs[refIndex].modelSlice = subMainModel.modelSlice

						return nil
					})
				}(ctx, refModel, refIndex)
			}
			if err := refGroup.Wait(); err != nil {
				return mapper, err
			}
		}

		if mainModel.modelSlice.Len() > 0 && len(mainModel.refFields) > 0 {
			var group, ctx = errgroup.WithContext(ctx)
			for index := 0; index < mainModel.modelSlice.Len(); index++ {
				bind(ctx, group, mainModel.modelSlice.Index(index), mainModel.refFields, mapper.modelStructs)
			}
			if err := group.Wait(); err != nil {
				return mapper, err
			}
		}
	}

	mapper.rowCount = rowCount
	mapper.paginateTotal = paginateTotal
	return mapper, nil
}

func OrmContext(ctx context.Context, model interface{}, rows *sqlx.Rows, options MapperOption) (Mapper, error) {
	return orm(ctx, model, rows, options)
}

func Orm(model interface{}, rows *sqlx.Rows, options MapperOption) (Mapper, error) {
	return orm(context.Background(), model, rows, options)
}

func GetSelector(models interface{}) string {
	faith := structs.New(models)
	fields := faith.Fields()
	tablename := getTableName(faith)
	var selectors = make([]string, 0)
	var patternSelector = func(tablename string, fieldDB string) string {
		return fmt.Sprintf(`%s.%s "%s.%s"`, tablename, fieldDB, tablename, fieldDB)
	}

	if len(fields) > 0 {
		for _, field := range fields {
			if field.Name() != TABLE_FIELD_NAME {
				columename := field.Tag(TAGNAME)
				if columename != "" && columename != "-" {
					selectors = append(selectors, patternSelector(tablename, columename))
				}
			}
		}
	}

	return strings.Join(selectors, ",")
}

func fillValueList(ms *modelStruct, columns []*sql.ColumnType, values []interface{}, options MapperOption) (reflect.Value, error) {
	slice := ms.modelSlice
	model := ms.model
	ptr := copy(reflect.ValueOf(model)).Interface()
	if err := fillValue(ptr, columns, values); err != nil {
		return slice, err
	}
	reflectValPtr := reflect.ValueOf(ptr)
	exists, err := isDuplicateByPK(ms, slice, reflectValPtr, options)
	if err != nil {
		return slice, err
	}
	if !exists {
		slice = reflect.Append(slice, reflectValPtr)
	}
	return slice, nil
}

func isDuplicateByPK(ms *modelStruct, slice reflect.Value, ptr reflect.Value, options MapperOption) (bool, error) {
	pkM := ms.pkM
	faith := structs.New(ptr.Interface())
	var storePK = func(elemFaith *structs.Struct, fields []string) (bool, error) {
		pkId, err := getIds(elemFaith, fields)
		if err != nil {
			return false, err
		}
		if pkId == "" || pkId == "0" || pkId == "false" {
			return true, nil
		}
		if _, ok := pkM.Load(pkId); ok {
			return true, nil
		}
		pkM.Store(pkId, slice.Len()-1)

		return false, nil
	}

	var allFields = []string{}
	var elemFaith *structs.Struct
	switch ms.IsMainModel() {
	case true:
		/* for main model */
		allFields, _ = getFieldMetaData(faith, options)
		elemFaith = structs.New(ptr.Interface())
	default:
		/* for reference model */
		allFields, _ = getFieldMetaData(faith, options)
		elemFaith = structs.New(ptr.Interface())
		if len(ms.refFields) > 0 {
			allFields = append(allFields, ms.refFields...)
		}
	}

	return storePK(elemFaith, allFields)
}

func bindReference(ctx context.Context, mainElem reflect.Value, mainRefFieldNames []string, allModels []modelStruct) error {
	faith := structs.New(mainElem.Interface())
	if len(mainRefFieldNames) > 0 {
		var group, _ = errgroup.WithContext(ctx)
		for _, refField := range mainRefFieldNames {
			func(mainElem reflect.Value, refField string, mainRefFieldNames []string, allModels []modelStruct) {
				group.Go(func() (err error) {
					defer func() {
						if panicErr := recovery(); panicErr != nil {
							err = panicErr
						}
					}()
					pkFieldRefDataField := mainElem.Elem().FieldByName(refField)
					if pkFieldRefDataField.Type().Kind() == reflect.Slice {
						pkFieldRefDataField.Set(reflect.MakeSlice(pkFieldRefDataField.Type(), 0, 0))
					}

					if tagVal := getTagValue(faith, refField, TAG_FK); tagVal != "" {
						fk := newForeignKeyFromTag(tagVal)
						if err := fk.Validate(); err != nil {
							return err
						}
						refModel := modelStructs(allModels).GetRefModelByFieldName(refField)
						if !refModel.IsZero() && refModel.modelSlice.Len() > 0 {
							for i := 0; i < refModel.modelSlice.Len(); i++ {
								refVal := copy(refModel.modelSlice.Index(i))
								if isJoin(faith, refVal.Interface(), fk.fkField1, fk.fkField2) {
									if pkFieldRefDataField.Type().Kind() == reflect.Ptr {
										/* object */
										pkFieldRefDataField = refVal
										break
									} else {
										/* slice */
										pkFieldRefDataField = reflect.Append(pkFieldRefDataField, refVal)
									}
								}
							}
						}
					}
					faith.Field(refField).Set(pkFieldRefDataField.Interface())
					return nil
				})
			}(mainElem, refField, mainRefFieldNames, allModels)
		}
		if err := group.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func isJoin(mainFaith *structs.Struct, refData interface{}, fkCol1Keys []string, fkCol2Keys []string) bool {
	checkEqual := func(fkCol1, fkCol2 string) bool {
		parentID := mainFaith.Field(fkCol1).Value()
		parentType := mainFaith.Field(fkCol1).Tag(TAG_TYPE)
		linkID := structs.New(refData).Field(fkCol2).Value()
		return equal(parentType, parentID, linkID)
	}
	var totalValid = len(fkCol1Keys)
	var isValid int
	for index, _ := range fkCol1Keys {
		if checkEqual(fkCol1Keys[index], fkCol2Keys[index]) {
			isValid++
		}
	}
	return totalValid == isValid
}
