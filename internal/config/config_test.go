package config

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestNewConfig(t *testing.T) {

	badfile, err := ioutil.TempFile("", "conf.")
	require.NoError(t, err, err)
	defer os.Remove(badfile.Name())
	badfile.WriteString(`aefSD
sadfg
RFABND FYGUMG
V`)
	badfile.Sync()

	goodfile, err := ioutil.TempFile("", "conf.")
	require.NoError(t, err, err)
	defer os.Remove(goodfile.Name())
	goodfile.WriteString(`[Server]
Address = "localhost"
Port = "80"`)
	goodfile.Sync()

	t.Run("No such file", func(t *testing.T) {
		c, e := NewConfig("adfergdth")
		require.Equal(t, Config{}, c)
		require.Error(t, e)
	})

	t.Run("Bad file", func(t *testing.T) {
		c, e := NewConfig(badfile.Name())
		require.Equal(t, Config{}, c)
		require.Error(t, e)
	})

	t.Run("TOML reading", func(t *testing.T) {
		c, e := NewConfig(goodfile.Name())
		require.Equal(t, "localhost", c.Server.Address)
		require.Equal(t, "80", c.Server.Port)
		require.NoError(t, e)
	})

}
