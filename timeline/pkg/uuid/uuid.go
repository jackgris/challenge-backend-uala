package uuid

import "github.com/rs/xid"

func New() string {
	return xid.New().String()
}

func IsValid(id string) bool {
	_, err := xid.FromString(id)
	return err == nil
}
