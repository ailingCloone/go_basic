package login

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/tokenglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"
	"nrs_customer_module_backend/internal/model/staff"
	"nrs_customer_module_backend/internal/model/staff_other_details"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *LoginLogic {
	return &LoginLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *LoginLogic) Login(req *types.PostLoginPortalReq) (resp *types.SuccessResponse, err error) {
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	exists, err := user.CheckStaff(l.ctx, conn, "EMAIL", req.Contact)
	if err != nil {
		return nil, err
	}

	hashPass := global.GenerateSha256(req.Code)

	if exists.Password != hashPass {
		err = fmt.Errorf(errorglobal.InvalidPassword)
		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	//check user have active record in oauth table or not, if yes set active = 0 to the active record
	err = user.CheckUserOauthActive(l.ctx, conn, exists.Id, currentTime, "staff")
	if err != nil {
		return nil, err
	}

	//email correct, password correct -> generate token & update into oauth
	//get data from middleware
	oauthData, expirationTime, err := user.GetOauthDataFromBeforeLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}

	tokenValiditySeconds := int64(c.TokenValiditySecondsLogin)
	tokenResp, err := tokenglobal.GenerateToken(tokenValiditySeconds)
	if err != nil {
		return nil, err
	}

	oauthData = SetUpOauthData(req, exists.Id, oauthData, exists.Role, tokenResp, currentTime)

	oauthModel := oauth.NewOauthModel(conn)
	if err := oauthModel.UpdateById(l.ctx, oauthData); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to login", map[string]interface{}{}), nil
	}

	response := StaffResponse(l.ctx, conn, oauthData, exists, expirationTime, c.UserImageUrl, c.UserImageUrl)

	return responseglobal.GenerateResponseBody(true, "Successfully login.", response), nil
}

func SetUpOauthData(req *types.PostLoginPortalReq, userId int64, oauthData *oauth.Oauth, role int64, tokenResp *types.TokenResponse, currentTime *time.Time) *oauth.Oauth {
	userIdNull := sql.NullInt64{Int64: userId, Valid: true}

	oauthData.StaffId = userIdNull
	oauthData.Scope = role // refer role table : -staff: staff.Role / customer: role.id = 1 - user

	oauthData.LoginBy = 1 //login by: -> 1- Email, 2- Contact, 3- Icno

	oauthData.AccessToken = tokenResp.AccessToken
	oauthData.RefreshToken = tokenResp.RefreshToken
	oauthData.ExpiresIn = tokenResp.ExpiresIn
	oauthData.Updated = *currentTime

	return oauthData
}

func StaffResponse(ctx context.Context, conn sqlx.SqlConn, oauthData *oauth.Oauth, exists *staff.Staff, expirationTime *time.Time, imageUrl string, webImageUrl string) (response map[string]interface{}) {
	staffOtherDetailsModel := staff_other_details.NewStaffOtherDetailsModel(conn)
	staffOtherInfo, err := staffOtherDetailsModel.FindOneByStaffId(ctx, oauthData.StaffId.Int64)
	if err == nil {
		if staffOtherInfo.ImageUrl != "" {
			imageUrl = staffOtherInfo.ImageUrl
		}
		if staffOtherInfo.WebImageUrl != "" {
			imageUrl = staffOtherInfo.WebImageUrl
		}
	}

	staffInfo := &staff.StaffInfo{
		Guid:        exists.Guid,
		ImageUrl:    imageUrl,
		WebImageUrl: webImageUrl,
		Name:        exists.Name,
		Contact:     exists.Contact,
		Email:       exists.Email,
		Icno:        exists.Icno,
		StaffCode:   exists.Code,
	}
	tokenInfo := &oauth.TokenInfo{
		AccessToken:  oauthData.AccessToken,
		RefreshToken: oauthData.RefreshToken,
		ExpiredAt:    fmt.Sprint((*expirationTime).Format(global.DefaultTimeFormat)),
		ExpiresIn:    oauthData.ExpiresIn,
	}

	response = map[string]interface{}{
		"token_info": tokenInfo,
		"staff_info": staffInfo,
	}
	return response
}
