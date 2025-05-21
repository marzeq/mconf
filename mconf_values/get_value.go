package mconf_values

import "math/big"

func UnwrapString(v MconfValue) (result string, ok bool) {
	ms, ok := v.(*MconfString)
	result = ms.Value
	return
}

func UnwrapBigInt(v MconfValue) (result *big.Int, ok bool) {
	mi, ok := v.(*MconfInt)
	result = mi.Value
	return
}

func UnwrapInt(v MconfValue) (result int64, ok bool) {
	mi, ok := v.(*MconfInt)
	if !ok {
		return 0, false
	}
	result = mi.Value.Int64()
	return
}

func UnwrapBigFloat(v MconfValue) (result *big.Float, ok bool) {
	mf, ok := v.(*MconfFloat)
	result = mf.Value
	return
}

func UnwrapFloat(v MconfValue) (result float64, ok bool) {
	mf, ok := v.(*MconfFloat)
	if !ok {
		return 0, false
	}
	result, _ = mf.Value.Float64()
	return
}

func UnwrapBool(v MconfValue) (result bool, ok bool) {
	mb, ok := v.(*MconfBool)
	result = mb.Value
	return
}

func UnwrapList(v MconfValue) (result []MconfValue, ok bool) {
	ml, ok := v.(*MconfList)
	result = ml.Value
	return
}

func UnwrapObject(v MconfValue) (result map[string]MconfValue, ok bool) {
	mo, ok := v.(*MconfObject)
	result = mo.Value
	return
}

func UnwrapNull(v MconfValue) (ok bool) {
	_, ok = v.(*MconfNull)
	return
}

func ObjGetString(v map[string]MconfValue, key string) (result string, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return "", false, false
	}
	ms, ok := got.(*MconfString)
	if !ok {
		return "", true, false
	}
	return ms.Value, true, true
}

func ObjGetBigInt(v map[string]MconfValue, key string) (result *big.Int, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return nil, false, false
	}
	mi, ok := got.(*MconfInt)
	if !ok {
		return nil, true, false
	}
	return mi.Value, true, true
}

func ObjGetInt(v map[string]MconfValue, key string) (result int64, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return 0, false, false
	}
	mi, ok := got.(*MconfInt)
	if !ok {
		return 0, true, false
	}
	result = mi.Value.Int64()
	return result, true, true
}

func ObjGetBigFloat(v map[string]MconfValue, key string) (result *big.Float, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return nil, false, false
	}
	mf, ok := got.(*MconfFloat)
	if !ok {
		return nil, true, false
	}
	return mf.Value, true, true
}

func ObjGetFloat(v map[string]MconfValue, key string) (result float64, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return 0, false, false
	}
	mf, ok := got.(*MconfFloat)
	if !ok {
		return 0, true, false
	}
	result, _ = mf.Value.Float64()
	return result, true, true
}

func ObjGetBool(v map[string]MconfValue, key string) (result bool, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return false, false, false
	}
	mb, ok := got.(*MconfBool)
	if !ok {
		return false, true, false
	}
	return mb.Value, true, true
}

func ObjGetList(v map[string]MconfValue, key string) (result []MconfValue, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return nil, false, false
	}
	ml, ok := got.(*MconfList)
	if !ok {
		return nil, true, false
	}
	return ml.Value, true, true
}

func ObjGetObject(v map[string]MconfValue, key string) (result map[string]MconfValue, exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return nil, false, false
	}
	mo, ok := got.(*MconfObject)
	if !ok {
		return nil, true, false
	}
	return mo.Value, true, true
}

func ObjGetNull(v map[string]MconfValue, key string) (exists bool, typeOk bool) {
	got, ok := v[key]
	if !ok {
		return false, false
	}
	_, ok = got.(*MconfNull)
	return true, ok
}
