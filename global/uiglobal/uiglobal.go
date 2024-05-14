package uiglobal

import (
	"context"
	"fmt"
	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/model/card_flow"
	"nrs_customer_module_backend/internal/model/ui"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func CheckContentTypeName(typeVal int64) (typeName string) {
	switch typeVal {
	case 1:
		typeName = "Text"
	case 2:
		typeName = "Int"
	case 3:
		typeName = "Textarea"
	case 4:
		typeName = "OTPCode"
	}
	return typeName
}

func CheckButtonTypeName(typeVal int64) (typeName string) {
	switch typeVal {
	case 1:
		typeName = "Btn_Normal"
	case 2:
		typeName = "Btn_Radio"
	case 3:
		typeName = "Btn_Text"
	case 4:
		typeName = "Btn_Switch"
	case 5:
		typeName = "Btn_Dropdown"
	case 6:
		typeName = "Btn_CheckBox"
	}
	return typeName
}

func GetUiHeader(uiModel ui.UiModel, uiHeaderId int64) (ui.Header, error) {
	uiHeader, err := uiModel.FindOneUiHeader(context.Background(), uiHeaderId)
	if err != nil {
		return ui.Header{}, err
	}

	header := ui.Header{
		ImageUrl:    uiHeader.ImageUrl,
		WebImageUrl: uiHeader.WebImageUrl,
		Title:       uiHeader.Title,
		Description: uiHeader.Description,
	}

	return header, nil
}

func BuildContentItem(uiContent ui.UIContent, c config.Config) ui.Content {
	typeName := CheckContentTypeName(uiContent.Type)
	description := strings.Replace(uiContent.Description, "{{minute}}", fmt.Sprint(c.SMSOtpExpire)+" seconds", -1)
	return ui.Content{
		Type:                  typeName,
		Dwname:                uiContent.Dwname,
		LeftIcon:              uiContent.LeftIcon,
		RightIcon:             uiContent.RightIcon,
		ImageUrl:              uiContent.ImageUrl,
		WebImageUrl:           uiContent.WebImageUrl,
		Title:                 uiContent.Title,
		Description:           description,
		Placeholder:           uiContent.Placeholder,
		Required:              uiContent.Required,
		ValidationDescription: uiContent.ValidationDescription,
	}
}

func BuildButtonItem(uiButton ui.UIButton) ui.Button {
	typeName := CheckButtonTypeName(uiButton.Type)

	return ui.Button{
		Type:                  typeName,
		Dwname:                uiButton.Dwname,
		Text:                  uiButton.Text,
		Action:                uiButton.Action,
		Required:              uiButton.Required,
		ValidationDescription: uiButton.ValidationDescription,
	}
}

func BuildButtonFromContent(uiContent ui.UIContent) ui.Button {
	typeName := CheckContentTypeName(uiContent.Type)

	return ui.Button{
		Type:                  typeName,
		Dwname:                uiContent.Dwname,
		Text:                  uiContent.Title, // Use content title as button text
		Required:              uiContent.Required,
		ValidationDescription: uiContent.ValidationDescription,
	}
}

func BuildContentFromButton(uiButton ui.UIButton) ui.Content {
	typeName := CheckButtonTypeName(uiButton.Type)

	return ui.Content{
		Type:                  typeName,
		Dwname:                uiButton.Dwname,
		Title:                 uiButton.Text,   // Use button text as content title
		Action:                uiButton.Action, // Use button text as content title
		Required:              uiButton.Required,
		ValidationDescription: uiButton.ValidationDescription,
	}
}

func SortedTabPriorities(tabMap map[int64]*ui.TabContent) []int64 {
	var keys []int64
	for key := range tabMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func GetTitle(summary ui.UISummary, uiModel ui.UiModel, ctx context.Context, id int64) (titles string) {
	if summary.UiTitleId.Int64 > 0 {
		uiTitle, err := uiModel.FindUiTitle(ctx, id)
		if err != nil {
			return ""
		}

		for _, title := range uiTitle {
			if title.Id == summary.UiTitleId.Int64 {
				titles = title.Title
			}
		}
	}
	return titles
}

func CardFlowInfo(ctx context.Context, conn sqlx.SqlConn, categoryId int64, index int) (int64, int64, card_flow.CardFlowUi, error) {
	cardFlowModel := card_flow.NewCardFlowModel(conn)
	cardFlow, err := cardFlowModel.FindOneCategoryId(ctx, categoryId)
	if err != nil {
		return 0, 0, card_flow.CardFlowUi{}, err
	}

	if len(*cardFlow) > index {
		cFlow := (*cardFlow)[index]
		return cFlow.UiHeaderId, cFlow.UiId, cFlow, nil
	}

	return 0, 0, card_flow.CardFlowUi{}, fmt.Errorf("card flow not found or index out of range")
}
