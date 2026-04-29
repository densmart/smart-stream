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
	0: "OK",

	// Auth errors (100-109)
	100: "Web user not found",
	101: "Invalid credentials or client inactive",
	102: "Failed to generate token",

	// Client errors (150-159)
	151: "Cannot create client",
	152: "Cannot retrieve client",
	153: "Cannot update client",
	154: "Cannot delete client",
	155: "Client not found",
	156: "Cannot search clients",
	157: "Cannot reset password",

	// Media errors (200-209)
	201: "Cannot create media",
	202: "Media not found",
	203: "Cannot update media",
	204: "Cannot delete media",
	205: "Media not found for update",
	206: "Cannot search media",

	// Playlist errors (210-219)
	211: "Cannot create playlist",
	212: "Playlist not found",
	213: "Cannot update playlist",
	214: "Cannot delete playlist",
	215: "Playlist not found for update",
	216: "Cannot search playlists",

	// Upload errors (220-229)
	221: "File size exceeds maximum allowed",
	222: "Invalid file type",
	223: "Invalid file extension",
	224: "Failed to create storage directory",
	225: "Failed to open uploaded file",
	226: "Failed to save file",
	227: "Failed to copy file content",

	// Streaming errors (240-249)
	240: "Media not found for streaming",
	241: "Media file not found on disk",
	242: "Error accessing media file",

	9999: "Unknown error",
}
