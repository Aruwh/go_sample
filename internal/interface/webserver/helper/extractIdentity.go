package helper

import (
	"fewoserv/internal/infrastructure/common"
	helper "fewoserv/internal/interface/webserver/helper/authentication"
	"net/http"
)

func ExtractIdentity(r *http.Request) helper.Identity {
	value := r.Context().Value(common.IdentityIdentifier)
	return value.(helper.Identity)
}
