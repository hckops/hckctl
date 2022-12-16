package box

import (
	"fmt"
)

// starts and/or attach to local container
func OpenLocalBox() string {
	// TODO docker and podman
	return fmt.Sprint("open-local-box")
}

// starts and/or attach to remote container
func OpenBox() string {
	// TODO invoke api
	return fmt.Sprint("open-box")
}
