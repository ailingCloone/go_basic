package activitylog

import (
	"context"
	"nrs_customer_module_backend/internal/model/activity_log"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func AddActivityLog(info *activity_log.ActivityLog, conn sqlx.SqlConn) (err error) {

	data := &activity_log.ActivityLog{
		StaffId:    info.StaffId,
		CustomerId: info.CustomerId,
		ReferTable: info.ReferTable,
		Action:     strings.ToLower(info.Action),
		Changes:    info.Changes,
		Created:    info.Created,
	}
	activityLogModel := activity_log.NewActivityLogModel(conn)
	// Insert data into the database
	_, err = activityLogModel.Insert(context.Background(), data)
	// insert record
	if err != nil {
		return err
	}

	return nil

}
