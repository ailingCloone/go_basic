package sms_template

import (
	"context"
	"database/sql"
	"strconv"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/sms_template"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *GetLogic {
	return &GetLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *GetLogic) Get(req *types.PostSTGetReq) (*types.SuccessResponse, error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	cardModel := card.NewCardModel(conn)
	cardResult, err := cardModel.FindOneGuid(context.Background(), req.Guid)
	if err != nil {
		return nil, err
	}

	smsTemplateModel := sms_template.NewSmsTemplateModel(conn)
	smsTemplates, err := smsTemplateModel.FindAllByCard(l.ctx, cardResult.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseglobal.GenerateResponseBody(false, "SMS Template is not found.", map[string]interface{}{}), nil
		}
		return nil, err
	}

	if smsTemplates == nil || len(*smsTemplates) == 0 {
		return responseglobal.GenerateResponseBody(false, "Failed to retrieve record.", map[string]interface{}{}), nil
	}

	var smsTemplateResponses []sms_template.SmsTemplateResponse
	for _, smsTemplate := range *smsTemplates {
		smsTemplateResponses = append(smsTemplateResponses, sms_template.SmsTemplateResponse{
			Guid:        smsTemplate.Guid,
			Title:       smsTemplate.Title,
			Description: smsTemplate.Description,
			Status:      strconv.Itoa(int(smsTemplate.Status)),
		})
	}

	data := map[string]interface{}{
		"templates": smsTemplateResponses,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", data), nil
}
