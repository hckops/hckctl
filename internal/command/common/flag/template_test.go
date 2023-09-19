package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTemplateSourceFlag(t *testing.T) {
	provider := CloudProviderFlag

	errRevision := ValidateTemplateSourceFlag(&provider, &TemplateSourceFlag{Revision: "invalid"})
	assert.EqualError(t, errRevision, "flag not supported: provider=cloud revision=invalid")

	errLocal := ValidateTemplateSourceFlag(&provider, &TemplateSourceFlag{Revision: "main", Local: true})
	assert.EqualError(t, errLocal, "flag not supported: provider=cloud local=true")
}
