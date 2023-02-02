package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_field(t *testing.T) {
	field, err := newField("val")
	assert.NoError(t, err)
	assert.Equal(t, "val", field.value)
	assert.Equal(t, true, field.matches("val"))
	assert.Equal(t, false, field.matches("value"))

	field, err = newField("*")
	assert.NoError(t, err)
	assert.Equal(t, "*", field.value)
	assert.NotNil(t, field.re)
	assert.Equal(t, true, field.matches("val"))
	assert.Equal(t, true, field.matches("value"))

	field, err = newField("(.+")
	assert.EqualError(t, err, "error parsing regexp: missing closing ): `^(.+$`")
	assert.NotNil(t, field)

	field, err = newField("foo_*")
	assert.NoError(t, err)
	assert.Equal(t, "foo_*", field.value)
	assert.NotNil(t, field.re)
	assert.Equal(t, false, field.matches("foo"))
	assert.Equal(t, true, field.matches("foo_bar"))
	assert.Equal(t, true, field.matches("foo_bar_widget"))

}

func Test_fieldString(t *testing.T) {
	field, err := newField("val")
	assert.NoError(t, err)
	assert.Equal(t, "val", field.String())
}
