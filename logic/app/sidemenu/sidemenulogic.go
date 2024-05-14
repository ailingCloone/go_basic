package sidemenu

import (
	"context"
	"fmt"
	"sort"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"
	"nrs_customer_module_backend/internal/model/sidemenu"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	// Initialize menuMap before using it
	menuMap := make(map[string][]sidemenu.PortalSidemenuResponseInterface)

	if oauthData.StaffId.Int64 != 0 {
		menuMap, err = StaffSidemenu(conn, l.ctx, oauthData)
		if err != nil {
			return nil, err
		}
	}

	if oauthData.CustomerId.Int64 != 0 {
		menuMap, err = CustomerSidemenu(conn, l.ctx, oauthData)
		if err != nil {
			return nil, err
		}
	}

	var menus []sidemenu.PortalSidemenu
	for categoryTitle, menuList := range menuMap {
		menus = append(menus, sidemenu.PortalSidemenu{
			CategoryTitle:          categoryTitle,
			PortalSidemenuResponse: menuList,
		})
	}

	sort.Slice(menus, func(i, j int) bool {
		return menus[i].CategoryTitle < menus[j].CategoryTitle
	})

	responseData := map[string]interface{}{
		"menu": menus,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", responseData), nil
}

func StaffSidemenu(conn sqlx.SqlConn, ctx context.Context, oauthData *oauth.Oauth) (map[string][]sidemenu.PortalSidemenuResponseInterface, error) {
	userPermission, err := user.GetUserPermission(ctx)
	if err != nil {
		return nil, err
	}

	categoryMenus := make(map[string][]sidemenu.PortalSidemenuResponseInterface)

	var resp *[]sidemenu.FindAll

	resp, err = sidemenu.NewSidemenuModel(conn).FindAll(ctx, fmt.Sprint(oauthData.Scope), oauthData.TenantId.Int64)
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
						Guid:       item.Guid.String,
						ImageUrl:   item.ImageUrl.String,
						Title:      item.Title.String,
						Action:     item.Action.String,
						Children:   children,
						Permission: []sidemenu.UserPermission{permission},
					}
				} else {
					responseItem = sidemenu.PortalSidemenuResponseWithoutSubmenu{
						Guid:       item.Guid.String,
						ImageUrl:   item.ImageUrl.String,
						Title:      item.Title.String,
						Action:     item.Action.String,
						Permission: []sidemenu.UserPermission{permission},
					}
				}
				categoryMenus[item.CategoryTitle] = append(categoryMenus[item.CategoryTitle], responseItem)
			}
		}
	}

	return categoryMenus, nil
}

func CustomerSidemenu(conn sqlx.SqlConn, ctx context.Context, oauthData *oauth.Oauth) (map[string][]sidemenu.PortalSidemenuResponseInterface, error) {
	menuMap := make(map[string][]sidemenu.PortalSidemenuResponseInterface)
	sidemenuModel := sidemenu.NewSidemenuModel(conn)

	resp, err := sidemenuModel.FindAll(ctx, fmt.Sprint(oauthData.Scope), oauthData.TenantId.Int64)
	if err != nil {
		return nil, err
	}

	// Organize submenus by their parent ID
	submenusMap := make(map[int64][]sidemenu.Submenu)
	for _, item := range *resp {
		submenu := sidemenu.Submenu{
			Title: item.Title.String,
			To:    item.Action.String,
		}
		submenusMap[item.ParentId] = append(submenusMap[item.ParentId], submenu)

	}

	// Build categoryMenus with main menu items and their submenus
	for _, item := range *resp {
		var responseItem sidemenu.PortalSidemenuResponseInterface
		if item.ParentId == 0 {
			children := submenusMap[item.Id]
			if len(children) > 0 {
				responseItem = sidemenu.PortalSidemenuResponseWithSubmenu{
					Guid:     item.Guid.String,
					ImageUrl: item.ImageUrl.String,
					Title:    item.Title.String,
					Action:   item.Action.String,
					Children: children,
				}
			} else {
				responseItem = sidemenu.PortalSidemenuResponseWithoutSubmenu{
					Guid:     item.Guid.String,
					ImageUrl: item.ImageUrl.String,
					Title:    item.Title.String,
					Action:   item.Action.String,
				}
			}
			menuMap[item.CategoryTitle] = append(menuMap[item.CategoryTitle], responseItem)
		}

	}

	return menuMap, nil
}
