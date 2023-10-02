package marshaller

import (
	"github.com/liip/sheriff"
	"github.com/wernerdweight/api-auth-go/auth/contract"
)

func marshal(v interface{}, groups []string) (interface{}, *contract.AuthError) {
	o := &sheriff.Options{
		Groups:          groups,
		IncludeEmptyTag: true,
	}
	data, err := sheriff.Marshal(o, v)
	if nil != err {
		return nil, contract.NewAuthError(contract.MarshallingError, map[string]string{"details": err.Error()})
	}
	return data, nil
}

func MarshalPublic(v interface{}) (interface{}, *contract.AuthError) {
	return marshal(v, []string{"public"})
}

func MarshalInternal(v interface{}) (interface{}, *contract.AuthError) {
	return marshal(v, []string{"internal"})
}
