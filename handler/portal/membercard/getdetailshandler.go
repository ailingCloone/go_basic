package membercard

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nrs_customer_module_backend/internal/logic/portal/membercard"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"
)

func GetDetailsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCardDetailsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := membercard.NewGetDetailsLogic(r.Context(), svcCtx)
		resp, err := l.GetDetails(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
