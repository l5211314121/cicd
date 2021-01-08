package lib

func RespErr(desc string, err string) map[string]interface{}{
	return map[string]interface{}{
		"error": desc + err,
	}
}
