package mconf_values

import "math/big"

func MconfUnwrapString(v MconfValue) (result string, ok bool) {
	ms, ok := v.(*MconfString)
	result = ms.Value
	return
}

func MconfUnwrapInt(v MconfValue) (result *big.Int, ok bool) {
	mi, ok := v.(*MconfInt)
	result = mi.Value
	return
}

func MconfUnwrapFloat(v MconfValue) (result *big.Float, ok bool) {
	mf, ok := v.(*MconfFloat)
	result = mf.Value
	return
}

func MconfUnwrapBool(v MconfValue) (result bool, ok bool) {
	mb, ok := v.(*MconfBool)
	result = mb.Value
	return
}

func MconfUnwrapList(v MconfValue) (result []MconfValue, ok bool) {
	ml, ok := v.(*MconfList)
	result = ml.Value
	return
}

func MconfUnwrapObject(v MconfValue) (result map[string]MconfValue, ok bool) {
	mo, ok := v.(*MconfObject)
	result = mo.Value
	return
}

func MconfUnwrapNull(v MconfValue) (ok bool) {
	_, ok = v.(*MconfNull)
	return
}
