package fbx

import (
	"reflect"

	"github.com/pkg/errors"
)

type DataType rune

const TYPE_BOOL DataType = 'C'
const TYPE_INT16 DataType = 'Y'
const TYPE_INT32 DataType = 'I'
const TYPE_INT64 DataType = 'L'
const TYPE_FLOAT32 DataType = 'F'
const TYPE_FLOAT64 DataType = 'D'

const TYPE_ARRAY_FLOAT32 DataType = 'f'
const TYPE_ARRAY_FLOAT64 DataType = 'd'
const TYPE_ARRAY_INT64 DataType = 'l'
const TYPE_ARRAY_INT32 DataType = 'i'
const TYPE_ARRAY_BOOL DataType = 'b'

const TYPE_STRING DataType = 'S'
const TYPE_RAW DataType = 'R'

func (dt DataType) IsArray() bool {
	switch dt {
	case TYPE_ARRAY_FLOAT64, TYPE_ARRAY_FLOAT32,
		TYPE_ARRAY_INT64, TYPE_ARRAY_INT32,
		TYPE_ARRAY_BOOL:
		return true
	}
	return false
}

func (dt DataType) IsSpecial() bool {
	switch dt {
	case TYPE_STRING, TYPE_RAW:
		return true
	}
	return false
}

func dataTypeFromInterface(v interface{}) (dt DataType, err error) {
	switch v.(type) {
	case bool:
		dt = TYPE_BOOL
	case int16:
		dt = TYPE_INT16
	case int32:
		dt = TYPE_INT32
	case int64:
		dt = TYPE_INT64
	case float32:
		dt = TYPE_FLOAT32
	case float64:
		dt = TYPE_FLOAT64
	case []bool:
		dt = TYPE_ARRAY_BOOL
	case []int32:
		dt = TYPE_ARRAY_INT32
	case []int64:
		dt = TYPE_ARRAY_INT64
	case []float32:
		dt = TYPE_ARRAY_FLOAT32
	case []float64:
		dt = TYPE_ARRAY_FLOAT64
	case []byte:
		dt = TYPE_RAW
	case string:
		dt = TYPE_STRING
	default:
		err = errors.Errorf("Can't find data type for type %v", reflect.TypeOf(v))
	}
	return dt, err
}
