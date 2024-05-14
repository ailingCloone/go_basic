package dashboard

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/request_status"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *DashboardLogic {
	return &DashboardLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *DashboardLogic) Dashboard(req *types.PostDashboardReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	requestStatusModel := request_status.NewRequestStatusModel(conn)
	cardModel := card.NewCardModel(conn)

	subtitles, err := requestStatusModel.FindAllTitle(l.ctx)
	if err != nil {
		return nil, err
	}
	cardInfo, err := cardModel.FindAllCard(l.ctx)
	if err != nil {
		return nil, err
	}

	var registrationSummary []request_status.Summary
	var memberSummary []request_status.Summary
	var cardSummary []card.CardInfo

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect shared data

	for _, subtitle := range *subtitles { // Dereference the pointer to slice
		wg.Add(1)
		go func(subtitle request_status.Title) {
			defer wg.Done()

			var statusCount int
			var statusCountMonthly int
			var statusCountPreviousMonth int
			var err error

			switch subtitle.Id {
			case 7: //id "7" -> pending approval
				statusCount, err = requestStatusModel.FindAllPending(l.ctx)
			case 8:
				statusCount, err = requestStatusModel.FindAllMonthly(l.ctx, 2) //status = 2 => Approved
			case 9:
				statusCount, err = requestStatusModel.FindAllMonthly(l.ctx, 3) //status = 3 => Rejected
			case 10:
				statusCount, err = cardModel.FindAllCustomerCard(l.ctx)
				if err == nil {
					statusCountMonthly, _ = cardModel.FindAllCustomerCardMonthly(l.ctx)
					statusCountPreviousMonth, _ = cardModel.FindAllCustomerCardPreviousMonth(l.ctx)
				}
			case 11:
				statusCount, err = cardModel.FindAllCustomerCardRenew(l.ctx)
				if err == nil {
					statusCountMonthly, _ = cardModel.FindAllCustomerCardRenewMonthly(l.ctx)
					statusCountPreviousMonth, _ = cardModel.FindAllCustomerCardRenewPreviousMonth(l.ctx)
				}
			default:
				statusCount = 0
			}

			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()

			switch subtitle.Id {
			case 7, 8, 9:
				registrationSummary = append(registrationSummary, request_status.Summary{
					Title:       subtitle.Title,
					Description: subtitle.Description,
					Value:       fmt.Sprint(statusCount),
				})
			case 10, 11:
				var increment int
				var difference float64
				if statusCount != 0 && statusCountPreviousMonth != 0 {
					difference = (float64(statusCountMonthly) / float64(statusCountPreviousMonth)) * 100
					increment = int(difference)
				} else if statusCountPreviousMonth == 0 {
					difference = (float64(statusCountMonthly) / float64(statusCountPreviousMonth)) * 100
					increment = int(difference)
				} else {
					increment = statusCount
				}
				percent := fmt.Sprintf("%+d%%", increment)

				comparedPercentage := strings.Replace(subtitle.Description, "{{ .percent }}", percent, -1)
				memberSummary = append(memberSummary, request_status.Summary{
					Title:       subtitle.Title,
					Description: comparedPercentage,
					Value:       fmt.Sprint(statusCount),
					Percentage:  percent,
				})
			case 12:
				for _, info := range *cardInfo {
					cardName := strings.Replace(subtitle.Title, "{{ .card }}", info.Title, -1)
					statusCount, _ = cardModel.FindAllMember(l.ctx, int64(info.Id))
					cardSummary = append(cardSummary, card.CardInfo{
						Title:       cardName,
						Value:       fmt.Sprint(statusCount),
						ImageUrl:    info.ImageUrl,
						WebImageUrl: info.WebImageUrl,
					})
				}
			}
		}(subtitle) // Pass the subtitle by value to prevent race conditions
	}

	wg.Wait()

	data := request_status.Data{
		RegistrationSummary: registrationSummary,
		MemberSummary:       memberSummary,
		CardSummary:         cardSummary,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", data), nil
}
