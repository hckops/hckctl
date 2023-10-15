package kubernetes

import (
	"io"

	_ "k8s.io/kubectl/pkg/cmd/cp"
	_ "unsafe"
)

// linkname exposes private packages

//go:linkname cpMakeTar k8s.io/kubectl/pkg/cmd/cp.makeTar
func cpMakeTar(srcPath, destPath string, writer io.Writer) error

//go:linkname cpStripPathShortcuts k8s.io/kubectl/pkg/cmd/cp.stripPathShortcuts
func cpStripPathShortcuts(p string) string
