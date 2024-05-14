package login

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/errorglobal"
	otpFunc "nrs_customer_module_backend/internal/global/otp"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/tokenglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/customer"
	"nrs_customer_module_backend/internal/model/customer_card"
	"nrs_customer_module_backend/internal/model/customer_other_details"
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

func (l *LoginLogic) Login(req *types.PostLoginAppReq) (resp *types.SuccessResponse, err error) {
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	var tokenInfo *oauth.TokenInfo

	var userId int64
	var password string

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	if req.From == "customer" {
		userId, password, resp, err = l.CustomerResponse(req, conn, tokenInfo, c.UserImageUrl, c.UserImageUrl)
		if err != nil {
			return nil, err
		}
	}

	var role int64

	if req.From == "staff" {
		userId, role, password, resp, err = l.StaffResponse(req, conn, tokenInfo, c.UserImageUrl, c.UserImageUrl)
		if err != nil {
			return nil, err
		}
	}

	err = l.UserVerification(req, conn, password, currentTime)
	if err != nil {
		return nil, err
	}

	//check user have active record in oauth table or not, if yes set active = 0 to the active record
	err = user.CheckUserOauthActive(l.ctx, conn, userId, currentTime, req.From)
	if err != nil {
		return nil, err
	}

	//email correct, password correct || otp correct -> generate token & update into oauth
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
	oauthData = SetUpOauthData(req, userId, oauthData, role, tokenResp, currentTime)

	oauthModel := oauth.NewOauthModel(conn)
	if err := oauthModel.UpdateById(l.ctx, oauthData); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to login.", map[string]interface{}{}), nil
	}

	tokenInfo = &oauth.TokenInfo{
		AccessToken:  oauthData.AccessToken,
		RefreshToken: oauthData.RefreshToken,
		ExpiredAt:    fmt.Sprint((*expirationTime).Format(global.DefaultTimeFormat)),
		ExpiresIn:    oauthData.ExpiresIn,
	}

	tokenInfoResp, ok := resp.Data.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("unable to convert interface{} to map[string]interface{}")
		return nil, err
	}
	tokenInfoResp["token_info"] = tokenInfo
	return resp, nil
}

func (l *LoginLogic) CustomerResponse(req *types.PostLoginAppReq, conn sqlx.SqlConn, tokenInfo *oauth.TokenInfo, imageUrl string, webImageUrl string) (userId int64, password string, resp *types.SuccessResponse, err error) {
	exists, err := user.CheckCustomer(l.ctx, conn, req.Type, req.Contact)
	if err != nil {
		return 0, "", nil, err
	}
	userId = exists.Id
	password = exists.Password

	customerCardModel := customer_card.NewCustomerCardModel(conn)
	cardInfor, err := customerCardModel.FindOneCustomerId(l.ctx, exists.Id)
	if err != nil {
		return 0, "", nil, err
	}

	var cardInfo []customer_card.CustomerCardInfo
	for _, info := range *cardInfor {
		cardInfo = append(cardInfo, customer_card.CustomerCardInfo{
			Guid:        info.Guid,
			Type:        info.Type,
			Description: info.Description,
			Expiry:      info.Expiry,
			Number:      info.Number,
		})
	}

	customerOtherDetailsModel := customer_other_details.NewCustomerOtherDetailsModel(conn)
	customerOtherInfo, err := customerOtherDetailsModel.FindOneByCustomerId(l.ctx, exists.Id)
	if err == nil {
		if customerOtherInfo.ImageUrl != "" {
			imageUrl = customerOtherInfo.ImageUrl
		}
		if customerOtherInfo.WebImageUrl != "" {
			webImageUrl = customerOtherInfo.WebImageUrl
		}
	}

	customerInfo := &customer.CustomerInfo{
		Guid:        exists.Guid,
		ImageUrl:    imageUrl,
		WebImageUrl: webImageUrl,
		Name:        exists.Username,
		Contact:     exists.Contact,
		Email:       exists.Email,
		Icno:        exists.Icno,
	}

	dataResp := map[string]interface{}{
		"token_info":    tokenInfo,
		"customer_info": customerInfo,
		"card_info":     cardInfo,
	}
	return userId, password, responseglobal.GenerateResponseBody(true, "Successfully login.", dataResp), err
}

func (l *LoginLogic) StaffResponse(req *types.PostLoginAppReq, conn sqlx.SqlConn, tokenInfo *oauth.TokenInfo, imageUrl, webImageUrl string) (userId, role int64, password string, resp *types.SuccessResponse, err error) {
	exists, err := user.CheckStaff(l.ctx, conn, req.Type, req.Contact)
	if err != nil {
		return 0, 0, "", nil, err
	}
	userId = exists.Id
	password = exists.Password
	role = exists.Role

	staffOtherDetailsModel := staff_other_details.NewStaffOtherDetailsModel(conn)
	staffOtherInfo, err := staffOtherDetailsModel.FindOneByStaffId(l.ctx, exists.Id)
	if err == nil {
		if staffOtherInfo.ImageUrl != "" {
			imageUrl = staffOtherInfo.ImageUrl
		}
		if staffOtherInfo.WebImageUrl != "" {
			webImageUrl = staffOtherInfo.WebImageUrl
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

	dataResp := map[string]interface{}{
		"staff_info": staffInfo,
		"token_info": tokenInfo,
	}
	return userId, role, password, responseglobal.GenerateResponseBody(true, "Successfully login.", dataResp), err
}

func (l *LoginLogic) UserVerification(req *types.PostLoginAppReq, conn sqlx.SqlConn, password string, currentTime *time.Time) (err error) {
	if req.Type == "EMAIL" {
		hashPass := global.GenerateSha256(req.Code)

		if password != hashPass {
			err = fmt.Errorf(errorglobal.InvalidPassword)
			return err
		}
	}

	if req.Type == "CONTACT" || req.Type == "IC" {
		otpList, err := otpFunc.OtpExpirationChecking(l.ctx, conn, req.TxId, currentTime)
		if err != nil {
			return err
		}

		if req.Code != fmt.Sprint(otpList.Code) {
			err = fmt.Errorf(errorglobal.InvalidOtp)
			return err
		}

		//update otp active status if verified otp
		otpList.Active = 2
		err = otpFunc.UpdateOTP(l.ctx, conn, otpList)
		if err != nil {
			return err
		}
	}
	return err

}

func SetUpOauthData(req *types.PostLoginAppReq, userId int64, oauthData *oauth.Oauth, role int64, tokenResp *types.TokenResponse, currentTime *time.Time) *oauth.Oauth {
	userIdNull := sql.NullInt64{Int64: userId, Valid: true}

	if req.From == "staff" {
		oauthData.StaffId = userIdNull
		oauthData.Scope = role // refer role table : -staff: staff.Role / customer: role.id = 1 - user
	}

	if req.From == "customer" {
		oauthData.CustomerId = userIdNull
		oauthData.Scope = 1 // refer role table : -staff: staff.Role / customer: role.id = 1 - user

	}
	if req.Type == "EMAIL" {
		//login by: -> 1- Email, 2- Contact, 3- Icno
		oauthData.LoginBy = 1
	}
	if req.Type == "CONTACT" {
		oauthData.LoginBy = 2
	}
	if req.Type == "IC" {
		oauthData.LoginBy = 3
	}

	oauthData.AccessToken = tokenResp.AccessToken
	oauthData.RefreshToken = tokenResp.RefreshToken
	oauthData.ExpiresIn = tokenResp.ExpiresIn
	oauthData.Updated = *currentTime

	return oauthData
}
