package marshaller

import (
	"github.com/liip/sheriff"
	"github.com/wernerdweight/api-auth-go/v2/auth/contract"
)

func Marshal(v interface{}, groups []string) (interface{}, *contract.AuthError) {
	o := &sheriff.Options{
		Groups:          groups,
		IncludeEmptyTag: true,
	}
	data, err := sheriff.Marshal(o, v)
	if nil != err {
		return nil, contract.NewInternalError(contract.MarshallingError, map[string]string{"details": err.Error()})
	}
	return data, nil
}

func MarshalPublic(v interface{}) (interface{}, *contract.AuthError) {
	return Marshal(v, []string{"public"})
}

func MarshalInternal(v interface{}) (interface{}, *contract.AuthError) {
	return Marshal(v, []string{"internal"})
}
