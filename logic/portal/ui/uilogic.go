package ui

import (
	"context"
	"errors"
	"sort"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/uiglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/ui"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UiLogic {
	return &UiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UiLogic) Ui(req *types.PostLoginUIReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	dataResp, err := GetUI(conn, l.ctx)
	if err != nil {
		return nil, err
	}
	return responseglobal.GenerateResponseBody(true, "Successfully retrieved records.", dataResp), nil
}

func GetUI(conn sqlx.SqlConn, ctx context.Context) (dataResp *ui.UiResp, err error) {
	uiModel := ui.NewUiModel(conn)
	uiResult, err := uiModel.FindOneCategoryId(context.Background(), 22)
	if err != nil {
		return nil, err
	}

	header, err := getUiHeader(uiModel, uiResult.UiHeaderId.Int64)
	if err != nil {
		return nil, err
	}

	uiSummaryList, err := uiModel.FindUiSummary(context.Background(), uiResult.Id)
	if err != nil {
		return nil, err
	}

	// Map to store tabs based on TabPriority
	tabMap := make(map[int64]*ui.TabContent)
	contentMap := make(map[string]ui.Content)
	buttonMap := make(map[string]ui.Button)
	for _, summary := range uiSummaryList {

		var titles string

		if summary.ReferTable == "" {
			return nil, errors.New("refer table value is empty")
		}

		switch summary.ReferTable {
		case "ui_content":
			uiContentList, err := uiModel.FindUiContent(ctx, uiResult.Id)
			if err != nil {
				return nil, err
			}

			for _, content := range uiContentList {
				if content.Id == summary.ReferId {
					titles = GetTitle(summary, uiModel, ctx, uiResult.Id)

					if summary.Type == 1 {
						// Build content item from content and add to contentMap
						contentMap[content.Guid] = buildContentItem(content)

						// Add content item to its respective tab's Content based on TabPriority
						tabPriority := summary.TabPriority
						if tab, ok := tabMap[tabPriority]; ok {
							// Tab exists, append content to existing tab's Content
							tab.Content = append(tab.Content, buildContentItem(content))
						} else {
							// Tab does not exist, create new tab with the content item
							tabMap[tabPriority] = &ui.TabContent{
								Title:   titles,
								Content: []ui.Content{buildContentItem(content)},
								Button:  []ui.Button{}, // Initialize empty slice for buttons
							}
						}
					} else if summary.Type == 2 {
						// Build button from content and add to buttonMap
						buttonMap[content.Guid] = buildButtonFromContent(content)

						// Add button item to its respective tab's Button based on TabPriority
						tabPriority := summary.TabPriority
						if tab, ok := tabMap[tabPriority]; ok {
							// Tab exists, append button to existing tab's Button
							tab.Button = append(tab.Button, buildButtonFromContent(content))
						} else {
							// Tab does not exist, create new tab with the button item
							tabMap[tabPriority] = &ui.TabContent{
								Title:   titles,
								Content: []ui.Content{}, // Initialize empty slice for content
								Button:  []ui.Button{buildButtonFromContent(content)},
							}
						}
					}

					break // Move to the next summary
				}
			}

		case "ui_button":

			uiButtonList, err := uiModel.FindUiButton(ctx, uiResult.Id)
			if err != nil {
				return nil, err
			}

			for _, button := range uiButtonList {
				if button.Id == summary.ReferId {
					titles = GetTitle(summary, uiModel, ctx, uiResult.Id)

					if summary.Type == 1 {
						// Build content item from button and add to contentMap
						contentMap[button.Guid] = buildContentFromButton(button)

						// Add content item to its respective tab's Content based on TabPriority
						tabPriority := summary.TabPriority
						if tab, ok := tabMap[tabPriority]; ok {
							// Tab exists, append content to existing tab's Content
							tab.Content = append(tab.Content, buildContentFromButton(button))
						} else {
							// Tab does not exist, create new tab with the content item
							tabMap[tabPriority] = &ui.TabContent{
								Title:   titles,
								Content: []ui.Content{buildContentFromButton(button)},
								Button:  []ui.Button{}, // Initialize empty slice for buttons
							}
						}
					} else if summary.Type == 2 {
						// Build button item from button and add to buttonMap
						buttonMap[button.Guid] = buildButtonItem(button)

						// Add button item to its respective tab's Button based on TabPriority
						tabPriority := summary.TabPriority
						if tab, ok := tabMap[tabPriority]; ok {
							// Tab exists, append button to existing tab's Button
							tab.Button = append(tab.Button, buildButtonItem(button))
						} else {
							// Tab does not exist, create new tab with the button item
							tabMap[tabPriority] = &ui.TabContent{
								Title:   titles,
								Content: []ui.Content{}, // Initialize empty slice for content
								Button:  []ui.Button{buildButtonItem(button)},
							}
						}
					}

					break // Move to the next summary
				}
			}
		}
	}

	// Prepare tabs in sorted order by TabPriority
	var tabs []ui.TabContent
	for _, tabPriority := range sortedTabPriorities(tabMap) {
		tabs = append(tabs, *tabMap[tabPriority])
	}

	// Prepare UI response with header and tabs
	dataResp = &ui.UiResp{
		SideImage: header.WebImageUrl,
		Header:    header,
		Tab:       tabs,
	}
	return dataResp, nil
}

func getUiHeader(uiModel ui.UiModel, uiHeaderId int64) (ui.Header, error) {
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

func buildContentItem(uiContent ui.UIContent) ui.Content {
	typeName := uiglobal.CheckContentTypeName(uiContent.Type)
	return ui.Content{
		Type:                  typeName,
		Dwname:                uiContent.Dwname,
		LeftIcon:              uiContent.LeftIcon,
		RightIcon:             uiContent.RightIcon,
		Title:                 uiContent.Title,
		Placeholder:           uiContent.Placeholder,
		Required:              uiContent.Required,
		ValidationDescription: uiContent.ValidationDescription,
	}
}

func buildButtonItem(uiButton ui.UIButton) ui.Button {
	typeName := uiglobal.CheckButtonTypeName(uiButton.Type)
	return ui.Button{
		Type:                  typeName,
		Dwname:                uiButton.Dwname,
		Text:                  uiButton.Text,
		Action:                uiButton.Action,
		Required:              uiButton.Required,
		ValidationDescription: uiButton.ValidationDescription,
	}
}

func buildButtonFromContent(uiContent ui.UIContent) ui.Button {
	typeName := uiglobal.CheckContentTypeName(uiContent.Type)

	return ui.Button{
		Type:                  typeName,
		Dwname:                uiContent.Dwname,
		Text:                  uiContent.Title, // Use content title as button text
		Required:              uiContent.Required,
		ValidationDescription: uiContent.ValidationDescription,
	}
}

func buildContentFromButton(uiButton ui.UIButton) ui.Content {
	typeName := uiglobal.CheckButtonTypeName(uiButton.Type)
	return ui.Content{
		Type:                  typeName,
		Dwname:                uiButton.Dwname,
		Title:                 uiButton.Text, // Use button text as content title
		Required:              uiButton.Required,
		ValidationDescription: uiButton.ValidationDescription,
	}
}

func sortedTabPriorities(tabMap map[int64]*ui.TabContent) []int64 {
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
