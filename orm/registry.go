package orm

import (
	"pheet-fiber-backend/helper"
	"reflect"
	"time"

	"git.innovasive.co.th/backend/models"
	helperModel "git.innovasive.co.th/backend/models"
	"github.com/fatih/structs"
	"github.com/gofrs/uuid"
	"github.com/guregu/null/zero"
	"github.com/spf13/cast"
)

type Registry interface {
	TypeName() string
	RegisterPkId(val interface{}) string
	Bind(field *structs.Field, val interface{}) error
	Equal(x interface{}, y interface{}) bool
}

var GlobalRegistry = map[string]Registry{
	(uid{}).TypeName():                              uid{},
	(str("")).TypeName():                            str(""),
	(zerouid{}).TypeName():                          zerouid{},
	(integer(0)).TypeName():                         (integer(0)),
	(integer64(0)).TypeName():                       (integer64(0)),
	(floater32(float32(0))).TypeName():              (floater32(0)),
	(floater64(float64(0))).TypeName():              (floater64(0)),
	(timestamp(helperModel.Timestamp{})).TypeName(): timestamp(helperModel.Timestamp{}),
	(date(helperModel.Date{})).TypeName():           date(helperModel.Date{}),
	(zeroString(zeroString{})).TypeName():           zeroString(zero.String{}),
	(zeroInt(zero.Int{})).TypeName():                zeroInt(zero.Int{}),
	(zeroFloat(zero.Float{})).TypeName():            zeroFloat(zero.Float{}),
	(zeroBool(zero.Bool{})).TypeName():              zeroBool(zero.Bool{}),
	(boolean(true)).TypeName():                      (boolean(true)),
}

/*
----------------------------------------
|
|	UUID
|
----------------------------------------
*/
type uid uuid.UUID

func (elem uid) TypeName() string {
	return "uuid"
}

func (elem uid) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	if _, ok := val.(uuid.UUID); ok {
		return val.(uuid.UUID).String()
	}

	return val.(*uuid.UUID).String()
}

func (elem uid) Bind(field *structs.Field, val interface{}) error {
	parseVal, err := uuid.FromString(cast.ToString(val))
	if err == nil {
		field.Set(&parseVal)
	}
	return nil
}

func (elem uid) Equal(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return false
	}
	return x.(*uuid.UUID).String() == y.(*uuid.UUID).String()
}

/*
----------------------------------------
|
|	Zero UUID
|
----------------------------------------
*/
type zerouid models.ZeroUUID

func (elem zerouid) TypeName() string {
	return "zerouuid"
}

func (elem zerouid) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	return val.(models.ZeroUUID).String()
}

func (elem zerouid) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		parseVal, err := models.NewZeroUUIDFromstring(cast.ToString(val))
		if err == nil {
			return field.Set(parseVal)
		}
	}
	return nil
}

func (elem zerouid) Equal(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return false
	}

	if reflect.TypeOf(x).String() == "models.ZeroUUID" && reflect.TypeOf(y).String() == "models.ZeroUUID" {
		if x.(models.ZeroUUID) == (models.ZeroUUID{}) || y.(models.ZeroUUID) == (models.ZeroUUID{}) {
			return false
		}
		return x.(models.ZeroUUID).String() == y.(models.ZeroUUID).String()
	} else if reflect.TypeOf(x).String() == "models.ZeroUUID" && reflect.TypeOf(y).String() == "*uuid.UUID" {
		if x.(models.ZeroUUID) == (models.ZeroUUID{}) || y == nil {
			return false
		}
		return x.(models.ZeroUUID).String() == y.(*uuid.UUID).String()
	} else if reflect.TypeOf(x).String() == "*uuid.UUID" && reflect.TypeOf(y).String() == "models.ZeroUUID" {
		if x == nil || y.(models.ZeroUUID) == (models.ZeroUUID{}) {
			return false
		}
		return y.(models.ZeroUUID).String() == x.(*uuid.UUID).String()
	}
	return false
}

/*
----------------------------------------
|
|	String
|
----------------------------------------
*/
type str string

func (elem str) TypeName() string {
	return "string"
}

func (elem str) RegisterPkId(val interface{}) string {
	return cast.ToString(val)
}

func (elem str) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		if cast.ToString(val) != "" {
			field.Set(cast.ToString(val))
		}
	}
	return nil
}
func (elem str) Equal(x interface{}, y interface{}) bool {
	if x.(string) == "" || y.(string) == "" {
		return false
	}
	return x.(string) == y.(string)
}

/*
----------------------------------------
|
|	int32
|
----------------------------------------
*/
type integer int

func (elem integer) TypeName() string {
	return "int32"
}

func (elem integer) RegisterPkId(val interface{}) string {
	return cast.ToString(val)
}

func (elem integer) Bind(field *structs.Field, val interface{}) error {
	if val == nil {
		return nil
	}
	field.Set(cast.ToInt(val))
	return nil
}
func (elem integer) Equal(x interface{}, y interface{}) bool {
	if cast.ToInt(x) == 0 || cast.ToInt(y) == 0 {
		return false
	}
	return cast.ToInt(x) == cast.ToInt(y)
}

/*
----------------------------------------
|
|	int64
|
----------------------------------------
*/
type integer64 int64

func (elem integer64) TypeName() string {
	return "int64"
}

func (elem integer64) RegisterPkId(val interface{}) string {
	return cast.ToString(val)
}

func (elem integer64) Bind(field *structs.Field, val interface{}) error {
	if val == nil {
		return nil
	}
	field.Set(cast.ToInt64(val))
	return nil
}
func (elem integer64) Equal(x interface{}, y interface{}) bool {
	if cast.ToInt64(x) == 0 || cast.ToInt64(y) == 0 {
		return false
	}
	return cast.ToInt64(x) == cast.ToInt64(y)
}

/*
----------------------------------------
|
|	float32
|
----------------------------------------
*/
type floater32 float64

func (elem floater32) TypeName() string {
	return "float32"
}

func (elem floater32) RegisterPkId(val interface{}) string {
	return cast.ToString(val)
}

func (elem floater32) Bind(field *structs.Field, val interface{}) error {
	if val == nil {
		return nil
	}
	field.Set(cast.ToFloat32(val))
	return nil
}
func (elem floater32) Equal(x interface{}, y interface{}) bool {
	if cast.ToFloat32(x) == 0 || cast.ToFloat32(y) == 0 {
		return false
	}
	return cast.ToFloat32(x) == cast.ToFloat32(y)
}

/*
----------------------------------------
|
|	float64
|
----------------------------------------
*/
type floater64 float64

func (elem floater64) TypeName() string {
	return "float64"
}

func (elem floater64) RegisterPkId(val interface{}) string {
	return cast.ToString(val)
}

func (elem floater64) Bind(field *structs.Field, val interface{}) error {
	if val == nil {
		return nil
	}
	field.Set(cast.ToFloat64(val))
	return nil
}
func (elem floater64) Equal(x interface{}, y interface{}) bool {
	if cast.ToFloat64(x) == 0 || cast.ToFloat64(y) == 0 {
		return false
	}
	return cast.ToFloat64(x) == cast.ToFloat64(y)
}

/*
----------------------------------------
|
|	timestamp
|
----------------------------------------
*/
type timestamp helper.Timestamp

func (elem timestamp) TypeName() string {
	return "timestamp"
}

func (elem timestamp) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	if v, ok := val.(*helper.Timestamp); ok {
		return v.String()
	}
	if v, ok := val.(helper.Timestamp); ok {
		return v.String()
	}
	return cast.ToString(val)
}

func (elem timestamp) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "time.Time":
			timestamp := helper.Timestamp(val.(time.Time))
			field.Set(&timestamp)
		case "string":
			timestamp := helper.NewTimestampFromString(cast.ToString(val))
			field.Set(&timestamp)
		}
	}
	return nil
}

func (elem timestamp) Equal(x interface{}, y interface{}) bool {
	p1, p1OK := x.(*helperModel.Timestamp)
	p2, p2OK := y.(*helperModel.Timestamp)
	if p1OK && p2OK {
		return p1.ToUnix() == p2.ToUnix()
	}
	return x.(helperModel.Timestamp).ToUnix() == y.(helperModel.Timestamp).ToUnix()
}

/*
----------------------------------------
|
|	date
|
----------------------------------------
*/
type date helperModel.Date

func (elem date) TypeName() string {
	return "date"
}

func (elem date) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	if v, ok := val.(*helperModel.Date); ok {
		return v.String()
	}
	if v, ok := val.(helperModel.Date); ok {
		return v.String()
	}
	return cast.ToString(val)
}

func (elem date) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "time.Time":
			dt := helperModel.Date(val.(time.Time))
			field.Set(&dt)
		case "string":
			dt := models.NewDateFromString(cast.ToString(val))
			field.Set(&dt)
		}
	}
	return nil
}

func (elem date) Equal(x interface{}, y interface{}) bool {
	p1, p1OK := x.(*helperModel.Date)
	p2, p2OK := y.(*helperModel.Date)
	if p1OK && p2OK {
		return p1.String() == p2.String()
	}
	return x.(helperModel.Date).String() == y.(helperModel.Date).String()
}

/*
----------------------------------------
|
|	zerostring
|
----------------------------------------
*/
type zeroString zero.String

func (elem zeroString) TypeName() string {
	return "zerostring"
}

func (elem zeroString) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	return val.(zero.String).ValueOrZero()
}

func (elem zeroString) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "zero.String":
			field.Set(val.(zero.String))
		case "string":
			dt := zero.StringFrom(cast.ToString(val))
			field.Set(dt)
		case "[]uint8":
			field.Set(zero.StringFrom(cast.ToString(val)))
		}
	}
	return nil
}

func (elem zeroString) Equal(x interface{}, y interface{}) bool {
	return x.(zero.String).ValueOrZero() == y.(zero.String).ValueOrZero()
}

/*
----------------------------------------
|
|	zeroint
|
----------------------------------------
*/
type zeroInt zero.Int

func (elem zeroInt) TypeName() string {
	return "zeroint"
}

func (elem zeroInt) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	return cast.ToString(val.(zero.Int).ValueOrZero())
}

func (elem zeroInt) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "zero.Int":
			field.Set(val.(zero.Int))
		case "string":
			dt := zero.IntFrom(int64(cast.ToInt(val)))
			field.Set(dt)
		case "int":
			dt := zero.IntFrom(int64(cast.ToInt(val)))
			field.Set(dt)
		}
	}
	return nil
}

func (elem zeroInt) Equal(x interface{}, y interface{}) bool {
	return x.(zero.Int).ValueOrZero() == y.(zero.Int).ValueOrZero()
}

/*
----------------------------------------
|
|	zerofloat
|
----------------------------------------
*/
type zeroFloat zero.Float

func (elem zeroFloat) TypeName() string {
	return "zerofloat"
}

func (elem zeroFloat) RegisterPkId(val interface{}) string {
	if val == nil || reflect.ValueOf(val).IsNil() || reflect.ValueOf(val).IsZero() {
		return ""
	}
	return cast.ToString(val.(zero.Float).ValueOrZero())
}

func (elem zeroFloat) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "zero.Float":
			field.Set(val.(zero.Float))
		case "string":
			dt := zero.FloatFrom(cast.ToFloat64(val))
			field.Set(dt)
		case "int":
			dt := zero.FloatFrom(cast.ToFloat64(val))
			field.Set(dt)
		case "float64":
			dt := zero.FloatFrom(cast.ToFloat64(val))
			field.Set(dt)
		}
	}
	return nil
}

func (elem zeroFloat) Equal(x interface{}, y interface{}) bool {
	return x.(zero.Float).ValueOrZero() == y.(zero.Float).ValueOrZero()
}

/*
----------------------------------------
|
|	zerobool
|
----------------------------------------
*/
type zeroBool zero.Bool

func (elem zeroBool) TypeName() string {
	return "zerobool"
}

func (elem zeroBool) RegisterPkId(val interface{}) string {
	return ""
}

func (elem zeroBool) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "zero.Bool":
			field.Set(val.(zero.Bool))
		case "string":
			dt := zero.BoolFrom(cast.ToBool(val))
			field.Set(dt)
		case "bool":
			dt := zero.BoolFrom(cast.ToBool(val))
			field.Set(dt)
		}
	}
	return nil
}

func (elem zeroBool) Equal(x interface{}, y interface{}) bool {
	return x.(zero.Bool).ValueOrZero() == y.(zero.Bool).ValueOrZero()
}

/*
----------------------------------------
|
|	bool
|
----------------------------------------
*/
type boolean bool

func (elem boolean) TypeName() string {
	return "bool"
}

func (elem boolean) RegisterPkId(val interface{}) string {
	return ""
}

func (elem boolean) Bind(field *structs.Field, val interface{}) error {
	if val != nil {
		switch reflect.TypeOf(val).String() {
		case "string":
			dt := cast.ToBool(val)
			field.Set(dt)
		case "bool":
			dt := cast.ToBool(val)
			field.Set(dt)
		}
	}
	return nil
}

func (elem boolean) Equal(x interface{}, y interface{}) bool {
	return cast.ToBool(x) == cast.ToBool(y)
}
