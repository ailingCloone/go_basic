package password

import (
	"context"
	"strings"

	"nrs_customer_module_backend/internal/config"
	emails "nrs_customer_module_backend/internal/global/email"
	otpFunc "nrs_customer_module_backend/internal/global/otp"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/email_template"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type ForgotPasswordLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewForgotPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *ForgotPasswordLogic {
	return &ForgotPasswordLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *ForgotPasswordLogic) ForgotPassword(req *types.PostForgotPasswordReq) (resp *types.SuccessResponse, err error) {
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	email := req.Email
	var username string
	exists, err := user.CheckStaff(l.ctx, conn, "EMAIL", email)
	if err != nil {
		return nil, err
	}
	username = exists.Name

	var sendStatus string
	emailModel := email_template.NewEmailTemplateModel(conn)
	templates, err := emailModel.FindOneReferId(l.ctx, "category", 27)
	if err != nil {
		return nil, err
	}

	subject := templates.Title
	body := templates.Description
	body = strings.ReplaceAll(body, "{{ .username }}", username)
	var auth int64 = 2 //1- SMS, 2- EMAIL [forgot password is using email for now]
	var from int64 = 1 //1- Forget Password, 2- Login Contact ,3- Login IC,  4- Register, 5- Profile Update
	otpRecord, sendStatus, err := otpFunc.GenerateEmailOtpRecord(l.ctx, conn, email, auth, from, c.EmailOtpExpire, body, subject)
	if err != nil {
		return nil, err
	}

	response := emails.EmailSentResponse(otpRecord, sendStatus)
	return responseglobal.GenerateResponseBody(true, "Successfully Sent Email.", response), nil
}
