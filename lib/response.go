package lib

func RespErr(desc string, err error) map[string]interface{}{
	return map[string]interface{}{
		"error": desc + err.Error(),
	}
}
