package setting

import (
	"context"
	"reflect"
	"strings"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/email_template"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *EditLogic {
	return &EditLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *EditLogic) Edit(req *types.PostCardSettingEditReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	if req.Dwname == "email_template" {
		// Prepare data for insertion
		data := &email_template.EmailTemplate{
			Guid:    req.Guid,
			Updated: *currentTime,
		}
		emailTemplateModel := email_template.NewEmailTemplateModel(conn)
		if err := emailTemplateModel.UpdateEmailTemplate(l.ctx, data, req.Value); err != nil {
			return responseglobal.GenerateResponseBody(false, "Failed to update template", map[string]interface{}{}), nil
		}
		return responseglobal.GenerateResponseBody(true, "Successfully update email template.", map[string]interface{}{}), nil
	}
	cardModel := card.NewCardModel(conn)
	cardDetails, err := cardModel.FindOneGuid(l.ctx, req.Guid)
	if err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to get card details.", map[string]interface{}{"content": cardDetails}), err
	}

	cardAttributes := card.Card{}

	matchingAttributes := getMatchingFieldsAsString(cardAttributes, req.Dwname)

	if matchingAttributes != "" {
		data := &card.Card{
			Guid:    req.Guid,
			Updated: *currentTime,
		}
		if err := cardModel.UpdateCard(l.ctx, data, matchingAttributes, req.Value); err != nil {
			return responseglobal.GenerateResponseBody(false, "Failed to update card.", map[string]interface{}{}), nil
		}

	} else {
		return responseglobal.GenerateResponseBody(false, "Failed to get dwname.", map[string]interface{}{}), nil
	}

	return responseglobal.GenerateResponseBody(true, "Successfully update record.", map[string]interface{}{}), nil
}

func getMatchingFieldsAsString(obj interface{}, substr string) string {
	var matchingFields []string

	t := reflect.TypeOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Id" {
			continue // excluded Id field
		}

		if strings.Contains(substr, strings.ToLower(field.Name)) {
			matchingFields = append(matchingFields, field.Name)
		}
	}

	return strings.Join(matchingFields, ",")
}
