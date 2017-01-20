package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadPassFile(t *testing.T) {
	passfile, err := ReadPassFile("../../data/pgpass_invalid")
	assert.Error(t, err)
	assert.Equal(t, "stat ../../data/pgpass_invalid: no such file or directory", err.Error())
	assert.Nil(t, passfile)

	passfile, err = ReadPassFile("../../data/pgpass_invalid_perms")
	assert.Error(t, err)
	assert.Equal(t, "pgpass file has invalid permissions", err.Error())
	assert.Nil(t, passfile)

	passfile, err = ReadPassFile("../../data/pgpass")
	assert.NoError(t, err)
	assert.Equal(t, "../../data/pgpass", passfile.Path)
	assert.Equal(t, 1, len(passfile.Entries))
}
