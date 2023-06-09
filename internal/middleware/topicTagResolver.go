package middleware

import (
	"os"

	"github.com/Bofry/structproto"
)

func TopicTagResolve(fieldname, token string) (*structproto.Tag, error) {
	var tag *structproto.Tag
	if token != "" && token != "-" {
		tag = &structproto.Tag{
			Name: os.ExpandEnv(token),
		}
	}
	return tag, nil
}
