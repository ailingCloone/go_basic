package setting

import (
	"context"
	"fmt"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/email_template"
	"nrs_customer_module_backend/internal/model/ui"
	"nrs_customer_module_backend/internal/model/user_access_setting"
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

func (l *GetLogic) Get(req *types.PostCardSettingGetReq) (resp *types.SuccessResponse, err error) {
	userPermission, err := user.GetUserPermission(l.ctx)
	if err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	cardModel := card.NewCardModel(conn)
	uiModel := ui.NewUiModel(conn)

	cardDetails, err := cardModel.FindOneGuid(l.ctx, req.Guid)
	if err != nil {
		return nil, err
	}

	ui, err := uiModel.FindOne(l.ctx, 30)
	if err != nil {
		return nil, err
	}

	uiButton, err := uiModel.FindUiButton(l.ctx, ui.Id)
	if err != nil {
		return nil, err
	}

	cardPrice, err := uiModel.FindUiOneContent(l.ctx, ui.Id)
	if err != nil {
		return nil, err
	}

	emailTemplateModel := email_template.NewEmailTemplateModel(conn)
	emailTemplateInfo, err := cardModel.FindOneEmailTemplateInfo(l.ctx)
	if err != nil {
		return nil, err
	}
	emailTemplate, err := emailTemplateModel.FindOne(l.ctx, emailTemplateInfo.ReferID)
	if err != nil {
		return nil, err
	}

	// var uiTemplate []card.CardSettingResponse
	var uiTemplate []interface{}
	var selectedStatus int64

	if cardPrice != nil {
		uiTemplate = append(uiTemplate, card.CardSettingResponseWithoutSelected{
			Guid:        cardDetails.Guid,
			Title:       cardPrice.Title,
			Type:        "int",
			Dwname:      cardPrice.Dwname,
			Description: cardPrice.Description,
			Value:       cardDetails.Price,
			// Value:       cardPrice.Placeholder,
		})
	}

	for _, u := range uiButton {
		switch u.Id {
		case 12: //ekyc button
			selectedStatus = cardDetails.Ekyc
		case 13: //Register button
			selectedStatus = cardDetails.Register
		case 14: //Renew button
			selectedStatus = cardDetails.Renew
		case 15: //Upgrade button
			selectedStatus = cardDetails.Upgrade
		case 16: //Send Invitation
			selectedStatus = cardDetails.Invitation
		}
		subUi, err := uiModel.FindSubUiButton(l.ctx, u.Id)
		if err != nil {
			return nil, err
		}

		switch u.Type {
		case 2: // Type 2: Radio
			if len(*subUi) > 0 {
				var submenus []card.TypeValue
				for _, subUi := range *subUi {
					selectedStatus = 0
					if subUi.Text == fmt.Sprint(cardDetails.Validity) {
						selectedStatus = 1
					}

					submenus = append(submenus, card.TypeValue{
						Guid:     subUi.Guid,
						Title:    subUi.Text,
						Selected: selectedStatus,
						Dwname:   subUi.Dwname,
					})
				}
				cardSettingResponse := card.CardSettingResponseWithoutSelected{
					Guid:  cardDetails.Guid,
					Title: u.Text,
					Type:  "radio_button",
					Value: submenus,
				}

				uiTemplate = append(uiTemplate, cardSettingResponse)
			}
		case 4: // Type 4: switch button
			cardSettingResponse := card.CardSettingResponse{
				Guid:     cardDetails.Guid,
				Title:    u.Text,
				Type:     "switch_button",
				Selected: selectedStatus,
				Dwname:   u.Dwname,
			}

			uiTemplate = append(uiTemplate, cardSettingResponse)
		default:
			cardSettingResponse := card.CardSettingResponse{
				Guid:     cardDetails.Guid,
				Title:    u.Text,
				Type:     fmt.Sprint(u.Type),
				Selected: selectedStatus,
				Dwname:   u.Dwname,
			}
			uiTemplate = append(uiTemplate, cardSettingResponse)

		}
	}

	if emailTemplate != nil {
		uiTemplate = append(uiTemplate, card.CardSettingResponseWithoutSelected{
			Guid:        emailTemplate.Guid,
			Title:       emailTemplate.Title,
			Type:        "html",
			Description: emailTemplate.Description,
			Dwname:      emailTemplateInfo.ReferTable,
			// Value:       cardPrice.Placeholder,
		})
	}

	var permissionSummaries []user_access_setting.Permission
	for _, permission := range *userPermission {
		summary := user_access_setting.Permission{
			Title: permission.Title,
			Allow: user_access_setting.Allow{
				Edit: permission.AllowEdit,
				View: permission.AllowView,
			},
		}
		permissionSummaries = append(permissionSummaries, summary)
	}

	content := append([]interface{}{}, uiTemplate...)
	response := map[string]interface{}{
		"content":    content,
		"permission": permissionSummaries,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", response), nil
}
