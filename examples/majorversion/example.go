package majorversion

import (
	"github.com/gofrs/uuid"
	uuidv5 "github.com/gofrs/uuid/v5"
)

func example() { //nolint: deadcode,unused
	// v3+incompatible import — should be blocked
	_ = uuid.Must(uuid.NewV4())

	// v5 import — should NOT be blocked because it is a different module
	_ = uuidv5.Must(uuidv5.NewV4())
}
