package splash

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/splash"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SplashLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSplashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SplashLogic {
	return &SplashLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SplashLogic) GetSplashes(info map[string]interface{}) (*types.SuccessResponse, error) {
	isFirst, err := strconv.ParseBool(fmt.Sprint(info["IsFirst"]))
	if err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	splashModel := splash.NewSplashModel(conn)

	resp, err := splashModel.FindAll(l.ctx, isFirst)
	if err != nil {
		return nil, err
	}
	fmt.Println("data from db", resp)

	pageMap := make(map[int]interface{})

	for _, find := range *resp {
		pageNumber := int(find.PageNumber.Int64)

		if _, value := pageMap[pageNumber]; !value {
			if find.ButtonText.Valid && find.ButtonAction.Valid && find.ButtonAction.String != "" {
				pageMap[pageNumber] = splash.SplashEntryWithButton{
					Guid:      find.Guid.String,
					ImageData: []splash.ImageData{},
					Button: splash.Button{
						ButtonText:   find.ButtonText.String,
						ButtonAction: find.ButtonAction.String,
					},
				}
			} else {
				pageMap[pageNumber] = splash.SplashEntryWithoutButton{
					Guid:      find.Guid.String,
					ImageData: []splash.ImageData{},
				}
			}
		}

		imageData := splash.ImageData{
			ImageUrl:     find.ImageUrl.String,
			WebImageUrl:  find.WebImageUrl.String,
			Description:  find.Description.String,
			Title:        find.Title.String,
			RedirectLink: find.RedirectLink.String,
		}

		switch entry := pageMap[pageNumber].(type) {
		case splash.SplashEntryWithButton:
			entry.ImageData = append(entry.ImageData, imageData)
			pageMap[pageNumber] = entry
		case splash.SplashEntryWithoutButton:
			entry.ImageData = append(entry.ImageData, imageData)
			pageMap[pageNumber] = entry
		}
	}

	// Extract keys from the map
	keys := make([]int, 0, len(pageMap))
	for key := range pageMap {
		keys = append(keys, key)
	}

	// Sort keys in ascending order
	sort.Ints(keys)

	var pageEntries []interface{}
	for _, key := range keys {
		pageEntries = append(pageEntries, pageMap[key])
	}

	responseData := map[string]interface{}{
		"list": pageEntries,
	}

	return responseglobal.GenerateResponseBody(true, "Success retrieve record.", responseData), nil
}
