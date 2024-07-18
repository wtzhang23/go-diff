package diffmatchpatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffLinesPretty(t *testing.T) {
	s1 := `
foo
line 1
line 2
bar
line 3
line 4
line 5
line 6
baz
line 7
line 8
line 9
line 10
line 11
`
	s2 := `
line 1
foo
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
bar
line 10
line 11
baz
`
	expected := `  
- foo
  line 1
+ foo
  line 2
- bar
  line 3
  line 4
  line 5
  line 6
- baz
  line 7
  line 8
  line 9
+ bar
  line 10
  line 11
+ baz
+ `
	dmp := New()
	out := dmp.DiffLinesPrettyText(LinesPrettyConfig{
		Spacing: " ",
		Context: 4,
	}, s1, s2)
	assert.Equal(t, expected, out)
}
