package storage

import (
	"github.com/traPtitech/Jomon/internal/service"
)

var ErrFileNotFound = service.NewNotFoundError("file not found")
