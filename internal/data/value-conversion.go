package data

import (
	"fmt"
	"strconv"
)

func ToBool(value any) (bool, error) {
	var boolValue bool
	switch value := value.(type) {
	case bool:
		boolValue = value
	case string:
		var err error
		boolValue, err = strconv.ParseBool(value)
		if err != nil {
			return boolValue, err
		}
	case int, float64:
		if value != 0 && value != 1 {
			return boolValue, fmt.Errorf("cannot convert value %v to bool", value)
		}
		boolValue = value != 0
	default:
		return boolValue, fmt.Errorf("cannot convert value %v with type %T to bool", value, value)
	}
	return boolValue, nil
}

func ToFloat(value any) (float64, error) {
	var floatValue float64
	switch value := value.(type) {
	case bool:
		if value {
			floatValue = 1
		} else {
			floatValue = 0
		}
	case string:
		var err error
		floatValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return floatValue, err
		}
	case int, float64:
		floatValue = value.(float64)
		return floatValue, nil
	default:
		return floatValue, fmt.Errorf("cannot convert value %v with type %T to float64", value, value)
	}
	return floatValue, nil
}

func ToInt(value any) (int, error) {
	var intValue int
	switch value := value.(type) {
	case bool:
		if value {
			intValue = 1
		} else {
			intValue = 0
		}
	case string:
		var err error
		intValue, err = strconv.Atoi(value)
		if err != nil {
			return intValue, err
		}
	case int:
		intValue = value
	case float64:
		intValue = int(value)
	default:
		return intValue, fmt.Errorf("cannot convert value %v with type %T to int", value, value)
	}
	return intValue, nil
}
