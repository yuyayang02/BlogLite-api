package httpresponse

//type response struct {
//	Code    string `json:"code"`
//	Message string `json:"message,omitempty"`
//	Data    any    `json:"data,omitempty"`
//}
//
//var ok = errors.ErrorCode{Code: 0}
//
//func Response(c *gin.Context, err error, data any) {
//	var (
//		code    = ok
//		message string
//	)
//	if err != nil {
//		if errors.As(err, &code) {
//			if errors.Is(code, errors.InternalServiceErr) {
//				message = ""
//			} else {
//				message = err.Error()
//			}
//		} else {
//			code = errors.Unimplemented
//			message = err.Error()
//		}
//	}
//
//	c.JSON(http.StatusOK, response{
//		Code:    code.Code,
//		Message: message,
//		Data:    data,
//	})
//}
//
//func Data(c *gin.Context, data any) {
//	c.JSON(http.StatusOK, response{
//		Code:    "OK",
//		Message: "",
//		Data:    data,
//	})
//}
//
//func Error(c *gin.Context, err error) {
//	Response(c, err, nil)
//}
//
//func ErrorOrOK(c *gin.Context, err error) {
//	if err != nil {
//		Response(c, err, nil)
//		return
//	}
//	OK(c)
//}
//
//func ErrorOrData(c *gin.Context, err error, data any) {
//	if err != nil {
//		Response(c, err, nil)
//		return
//	}
//	Data(c, data)
//}
//
//func OK(c *gin.Context) {
//	c.JSON(http.StatusOK, response{
//		Code:    "OK",
//		Message: "Successful",
//		Data:    nil,
//	})
//}
