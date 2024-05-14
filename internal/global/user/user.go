package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nrs_customer_module_backend/internal/model/customer"
	"nrs_customer_module_backend/internal/model/oauth"
	"nrs_customer_module_backend/internal/model/staff"
	"nrs_customer_module_backend/internal/model/user_access_setting"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func CheckCustomer(ctx context.Context, conn sqlx.SqlConn, identifierType string, contact string) (exists *customer.Customer, err error) {
	customerModel := customer.NewCustomerModel(conn)
	if identifierType == "CONTACT" {
		exists, err = customerModel.FindOneContact(ctx, contact)
		if err != nil {
			return nil, err
		}
	}
	if identifierType == "EMAIL" {
		exists, err = customerModel.FindOneEmail(ctx, contact)
		if err != nil {
			return nil, err
		}
	}
	if identifierType == "IC" {
		exists, err = customerModel.FindOneIC(ctx, contact)
		if err != nil {
			return nil, err
		}
	}

	if exists.Id == 0 {
		err = fmt.Errorf("record not found.")
		return nil, err
	}

	return exists, nil

}

func CheckStaff(ctx context.Context, conn sqlx.SqlConn, identifierType string, contact string) (exists *staff.Staff, err error) {
	staffModel := staff.NewStaffModel(conn)
	if identifierType == "CONTACT" {
		exists, err = staffModel.FindOneContact(ctx, contact)
		fmt.Println("exists staff for CONTACT", exists)
		fmt.Println("exists staff for CONTACT err", err)
		if err != nil {
			return nil, err
		}
	}
	if identifierType == "EMAIL" {
		exists, err = staffModel.FindOneEmail(ctx, contact)
		fmt.Println("exists staff for EMAIL", exists)
		fmt.Println("exists staff for EMAIL err", err)
		if err != nil {
			return nil, err
		}
		fmt.Println("after if err!= nil")
	}
	if identifierType == "IC" {
		exists, err = staffModel.FindOneIC(ctx, contact)
		fmt.Println("exists staff for IC", exists)
		if err != nil {
			return nil, err
		}
	}

	if exists.Id == 0 {
		err = fmt.Errorf("record not found.")
		return nil, err
	}

	return exists, nil
}

func CheckUserOauthActive(ctx context.Context, conn sqlx.SqlConn, userId int64, currentTime *time.Time, from string) error {
	oauthModel := oauth.NewOauthModel(conn)

	// Retrieve OAuth record based on 'from'
	var accActive *oauth.Oauth
	var err error

	switch from {
	case "customer":
		accActive, err = oauthModel.CheckOauthCustomerActive(ctx, userId)
	case "staff":
		accActive, err = oauthModel.CheckOauthStaffActive(ctx, userId)
	default:
		return fmt.Errorf("unsupported user type: %s", from)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			// Record not found, proceed without error
			return nil
		}
		return err // Return other errors
	}

	if accActive != nil {
		// Update existing OAuth record
		accActive.Active = 0
		accActive.Updated = *currentTime
		if err := oauthModel.UpdateById(ctx, accActive); err != nil {
			return err // Return error if update fails
		}
	}

	return nil
}

func GetOauthDataFromBeforeLoginMiddleware(ctx context.Context) (*oauth.Oauth, *time.Time, error) {
	oauthData, ok := ctx.Value("oauthData").(*oauth.Oauth)
	if !ok {
		return nil, nil, errors.New("OAuth data not found in context")
	}
	expirationTime, ok := ctx.Value("expirationTime").(*time.Time)
	if !ok {
		return nil, nil, errors.New("expiration time not found in context")
	}

	return oauthData, expirationTime, nil
}

func GetOauthDataFromAfterLoginMiddleware(ctx context.Context) (*oauth.Oauth, error) {
	oauthData, ok := ctx.Value("oauthData").(*oauth.Oauth)
	if !ok {
		return nil, errors.New("OAuth data not found in context")
	}

	return oauthData, nil
}

func GetUserPermission(ctx context.Context) (*[]user_access_setting.UserAccessSettings, error) {
	userPermission, ok := ctx.Value("userPermission").(*[]user_access_setting.UserAccessSettings)
	if !ok {
		return nil, errors.New("user permission data not found in context")
	}

	return userPermission, nil
}
