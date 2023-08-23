package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSourceFlag(t *testing.T) {
	provider := CloudProviderFlag

	errRevision := ValidateSourceFlag(&provider, &SourceFlag{Revision: "invalid"})
	assert.EqualError(t, errRevision, "flag not supported: provider=cloud revision=invalid")

	errLocal := ValidateSourceFlag(&provider, &SourceFlag{Revision: "main", Local: true})
	assert.EqualError(t, errLocal, "flag not supported: provider=cloud local=true")
}
