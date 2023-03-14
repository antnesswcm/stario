package stario

import "strconv"

type InputMsg struct {
	msg string
	err error
}

func (im InputMsg) String() (string, error) {
	// 如果InputMsg包含错误，直接返回空字符串和错误信息
	if im.err != nil {
		return "", im.err
	}
	// 否则返回InputMsg中的msg信息和空的错误信息
	return im.msg, nil
}

func (im InputMsg) MustString() string {
	res, _ := im.String()
	return res
}

func (im InputMsg) Int() (int, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.Atoi(im.msg)
}

func (im InputMsg) MustInt() int {
	res, _ := im.Int()
	return res
}

func (im InputMsg) Int64() (int64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseInt(im.msg, 10, 64)
}

func (im InputMsg) MustInt64() int64 {
	res, _ := im.Int64()
	return res
}

func (im InputMsg) Uint64() (uint64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseUint(im.msg, 10, 64)
}

func (im InputMsg) MustUint64() uint64 {
	res, _ := im.Uint64()
	return res
}

func (im InputMsg) Bool() (bool, error) {
	if im.err != nil {
		return false, im.err
	}
	return strconv.ParseBool(im.msg)
}

func (im InputMsg) MustBool() bool {
	res, _ := im.Bool()
	return res
}

func (im InputMsg) Float64() (float64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseFloat(im.msg, 64)
}

func (im InputMsg) MustFloat64() float64 {
	res, _ := im.Float64()
	return res
}

func (im InputMsg) Float32() (float32, error) {
	if im.err != nil {
		return 0, im.err
	}
	res, err := strconv.ParseFloat(im.msg, 32)
	return float32(res), err
}

func (im InputMsg) MustFloat32() float32 {
	res, _ := im.Float32()
	return res
}
