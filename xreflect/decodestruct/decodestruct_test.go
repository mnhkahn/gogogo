// Package decodestruct
package decodestruct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Config struct {
	Link     string `tag:"link"`
	NotExist int    `tag:"not_exist"`
}

func TestDecode(t *testing.T) {
	c := new(Config)
	_, err := Decode(c, "tag", FuncMap{
		"link": func() (interface{}, error) {
			return "http://cyeam.com", nil
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, c.Link, "http://cyeam.com")
	assert.Equal(t, c.NotExist, 0)
}
