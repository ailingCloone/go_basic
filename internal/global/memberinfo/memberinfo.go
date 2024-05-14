package memberinfo

import (
	"fmt"
	"nrs_customer_module_backend/internal/model/request_status"
	"strings"
)

type TabList struct {
	Title   string        `json:"title"`
	Status  int           `json:"status"`
	Content []DataAppList `json:"content"`
}

type DataAppList struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
	Icno string `json:"icno"`
	// Email string `json:"email"`
	// Card   string           `json:"card"`
	// Status string           `json:"status"`
	// Action []DataListAction `json:"action"`
}

type DataList struct {
	Guid   string           `json:"guid"`
	Name   string           `json:"name"`
	Email  string           `json:"email"`
	Card   string           `json:"card"`
	Status string           `json:"status"`
	Action []DataListAction `json:"action"`
}

type DataListAction struct {
	Title    string `json:"title"`
	Edit     int    `json:"edit"`
	Show     int    `json:"show"`
	Redirect string `json:"redirect"`
}

func PortalStatusAction(cusCardCode string, status int64, record DataList) (Result DataList) {
	var approveBtn DataListAction
	approveBtn.Title = "Approve"
	approveBtn.Edit = 1
	approveBtn.Show = 1
	approveBtn.Redirect = ""
	var rejectBtn DataListAction
	rejectBtn.Title = "Reject"
	rejectBtn.Edit = 1
	rejectBtn.Show = 1
	rejectBtn.Redirect = ""
	// var viewBtn DataListAction
	// viewBtn.Title = "View"
	// viewBtn.Edit = 0
	// viewBtn.Show = 1
	// viewBtn.Redirect = ""

	isEbsc := strings.ToLower(cusCardCode) == "ebsc" || strings.ToLower(cusCardCode) == "bsc"
	// if ebsc no need show status and button
	if !isEbsc {
		/*
			"1": true, // Pending
			"2": true, // Approve
			"3": true, // Reject
			"4": true, // All
		*/
		switch status {
		case 1:
			record.Status = "Pending Approval"
		case 2:
			record.Status = "Approve"
		case 3:
			record.Status = "Reject"
		}
		isNoNeedEdit := status == 2 || status == 3
		if isNoNeedEdit {
			// show but no editable
			approveBtn.Edit = 0
			rejectBtn.Edit = 0
		}
		record.Action = append(record.Action, approveBtn)
		record.Action = append(record.Action, rejectBtn)
		// record.Action = append(record.Action, viewBtn)
	}

	return record
}

func AppContent(list *[]request_status.MemberRequestList, page string) (record []TabList) {
	var pendingTabList, approveTabList TabList
	// group by status
	grouped := make(map[int64][]DataAppList)
	for _, v := range *list {
		status := v.Status // Assuming Status is the field containing the status value
		data := DataAppList{
			Guid: v.Guid,
			Name: v.CusFullname,
			Icno: v.CusIcno,
		}
		grouped[status] = append(grouped[status], data)
	}

	for status, items := range grouped {
		for _, item := range items {
			fmt.Println(item)
			switch status {
			case 1: // "Pending Approval"
				pendingTabList.Title = "Pending Approval"
				pendingTabList.Status = 1
				pendingTabList.Content = append(pendingTabList.Content, item)
			case 2: // "Approve"
				approveTabList.Title = "Approved"
				approveTabList.Status = 2
				approveTabList.Content = append(approveTabList.Content, item)
			case 3: // "Reject"
			default: // "All"
			}
		}
	}

	if pendingTabList.Status == 0 {
		pendingTabList.Title = "Pending Approval"
		pendingTabList.Status = 1
		pendingTabList.Content = []DataAppList{}
	}

	if approveTabList.Status == 0 {
		approveTabList.Title = "Approved"
		approveTabList.Status = 2
		approveTabList.Content = []DataAppList{}
	}

	switch page {
	case "app_registration_list_get_list":
		record = append(record, pendingTabList)
		record = append(record, approveTabList)
	}

	return record
}
