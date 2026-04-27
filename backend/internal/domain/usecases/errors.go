package usecases

import "github.com/densmart/smart-stream/pkg/logger"

type UCError struct {
	Code     uint
	HttpCode int
	Message  string
}

func (e *UCError) Error() string {
	logger.Debugf("[%d] %s", e.Code, e.Message)
	if e.Message == "" {
		if msg, exists := UCErrorCode[e.Code]; exists {
			return msg
		}
	}
	return e.Message
}

var UCErrorCode = map[uint]string{
	0:   "OK",
	100: "Web user not found",

	151:  "Cannot create client",
	152:  "Cannot retrieve client",
	153:  "Cannot update client",
	154:  "Cannot delete client",
	155:  "Client not found",
	156:  "Cannot search clients",
	157:  "Cannot reset password",
	9999: "Unknown error",
}
