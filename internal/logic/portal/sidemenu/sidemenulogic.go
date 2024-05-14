package sidemenu

import (
	"context"
	"fmt"
	"sort"

	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model/sidemenu"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"nrs_customer_module_backend/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// redisCon  = newConfig.RedisConnect()
	filename  = "portal_sidemenulogic"
	appLogger = createfile.New(filename)
)

type SidemenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSidemenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SidemenuLogic {
	return &SidemenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SidemenuLogic) Sidemenu() (*types.SuccessResponse, error) {
	oauthData, err := user.GetOauthDataFromAfterLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}
	userPermission, err := user.GetUserPermission(l.ctx)
	if err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	sidemenuModel := sidemenu.NewSidemenuModel(conn)
	resp, err := sidemenuModel.FindAll(l.ctx, fmt.Sprint(oauthData.Scope), oauthData.TenantId.Int64)
	if err != nil {
		return nil, err
	}

	permissionMap := make(map[int64]sidemenu.UserPermission)
	for _, perm := range *userPermission {
		permissionMap[perm.SubModuleId.Int64] = sidemenu.UserPermission{
			AllowAdd:    perm.AllowAdd,
			AllowEdit:   perm.AllowEdit,
			AllowDelete: perm.AllowDelete,
			AllowView:   perm.AllowView,
		}
	}

	categoryMenus := make(map[string][]sidemenu.PortalSidemenuResponseInterface)

	// Organize submenus by their parent ID
	submenusMap := make(map[int64][]sidemenu.Submenu)
	for _, item := range *resp {
		if permission, ok := permissionMap[item.Id]; ok && permission.AllowView != 0 {
			submenu := sidemenu.Submenu{
				Title:      item.Title.String,
				To:         item.Action.String,
				Permission: []sidemenu.UserPermission{permission},
			}
			submenusMap[item.ParentId] = append(submenusMap[item.ParentId], submenu)
		}
	}

	// Build categoryMenus with main menu items and their submenus
	for _, item := range *resp {
		if permission, ok := permissionMap[item.Id]; ok && permission.AllowView != 0 {
			var responseItem sidemenu.PortalSidemenuResponseInterface
			if item.ParentId == 0 {
				children := submenusMap[item.Id]
				if len(children) > 0 {
					responseItem = sidemenu.PortalSidemenuResponseWithSubmenu{
						Guid:        item.Guid.String,
						ImageUrl:    item.ImageUrl.String,
						WebImageUrl: item.ImageUrl.String,
						Title:       item.Title.String,
						Action:      item.Action.String,
						Children:    children,
						Permission:  []sidemenu.UserPermission{permission},
					}
				} else {
					responseItem = sidemenu.PortalSidemenuResponseWithoutSubmenu{
						Guid:        item.Guid.String,
						ImageUrl:    item.ImageUrl.String,
						WebImageUrl: item.WebImageUrl.String,
						Title:       item.Title.String,
						Action:      item.Action.String,
						Permission:  []sidemenu.UserPermission{permission},
					}
				}
				categoryMenus[item.CategoryTitle] = append(categoryMenus[item.CategoryTitle], responseItem)
			}
		}
	}

	var menus []sidemenu.PortalSidemenu
	for category, menuItems := range categoryMenus {
		menus = append(menus, sidemenu.PortalSidemenu{
			CategoryTitle:          category,
			PortalSidemenuResponse: menuItems,
		})
	}

	// Sort menus by category title
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].CategoryTitle < menus[j].CategoryTitle
	})

	responseData := map[string]interface{}{
		"menu": menus,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", responseData), nil
}
