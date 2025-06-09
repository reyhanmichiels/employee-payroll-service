package usecase

import (
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/hash"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/domain"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/usecase/user"
)

type Usecases struct {
	User user.Interface
}

type InitParam struct {
	Dom  *domain.Domains
	Json parser.JSONInterface
	Log  log.Interface
	Hash hash.Interface
	Auth auth.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User: user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
	}
}
