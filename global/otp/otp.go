package otp

import (
	"context"
	"fmt"
	"math/rand"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	emails "nrs_customer_module_backend/internal/global/email"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/global/sms"
	"strings"

	otps "nrs_customer_module_backend/internal/model/otp"

	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func GenerateOTP() int64 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	otp := int64(rng.Intn(900000) + 100000) // Generate a 6-digit OTP as an int64
	return otp
}

func InsertOTP(ctx context.Context, conn sqlx.SqlConn, contact string, auth, from int64, expiresIn int) (*otps.Otp, error) {
	otp := GenerateOTP()

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	data := &otps.Otp{
		Guid:       global.GenerateGuid(),
		AuthType:   auth,    // 1- SMS, 2- EMAIL
		SendStatus: 0,       // 1- Success, 2- Failed
		From:       from,    // 1- Forget Password, 2- Login Contact ,3- Login IC,  4- Register, 5- Profile Update
		Value:      contact, // Can be contact or email
		Code:       otp,
		ExpiresIn:  expiresIn,
		Created:    *currentTime,
		Active:     1,
	}

	otpModel := otps.NewOtpModel(conn)
	_, err = otpModel.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateOTP updates an existing OTP record in the database by GUID.
func UpdateOTP(ctx context.Context, conn sqlx.SqlConn, otp *otps.Otp) error {
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return err
	}
	otp.SendTime = *currentTime
	otpModel := otps.NewOtpModel(conn)
	err = otpModel.UpdateByGuid(ctx, otp)
	if err != nil {
		return err
	}

	return nil
}

// check otp valid or not
func ExpiredChecking(createdTime string, expiresInSeconds int, currentTime time.Time) bool {
	expirationTime, err := GetTokenExpirationTime(createdTime, expiresInSeconds)
	if err != nil {
		return false
	}

	if currentTime.After(*expirationTime) {
		return true //expired
	}
	return false
}

// get expiration time func
func GetTokenExpirationTime(createdTime string, expiresInSeconds int) (*time.Time, error) {
	createdAt, err := time.Parse(global.DefaultTimeFormat, createdTime)
	if err != nil {
		return nil, err
	}

	expirationTime := createdAt.Add(time.Duration(expiresInSeconds) * time.Second)

	return &expirationTime, err
}

func GenerateEmailOtpRecord(ctx context.Context, conn sqlx.SqlConn, email string, auth, from int64, expiresIn int, body, subject string) (otpRecord *otps.Otp, sendStatus string, err error) {
	otpRecord, err = InsertOTP(ctx, conn, email, auth, from, expiresIn)
	if err != nil {
		return nil, "", err
	}
	body = strings.ReplaceAll(body, "{{ .otp }}", fmt.Sprint(otpRecord.Code))

	err = emails.SendOTPEmail(email, otpRecord.Code, body, subject, "")
	if err != nil {
		otpRecord.SendStatus = 2 // Failed
		sendStatus = "FAILED"
		err = UpdateOTP(ctx, conn, otpRecord)
		if err != nil {
			return nil, "", err
		}
		return nil, "", err
	}
	sendStatus = "SUCCESS"
	otpRecord.SendStatus = 1 // Success
	err = UpdateOTP(ctx, conn, otpRecord)
	if err != nil {
		return nil, "", err
	}

	return otpRecord, sendStatus, nil
}

func GenerateSMSOtpRecord(ctx context.Context, conn sqlx.SqlConn, contactNo string, auth, from int64, expiresIn int, appLogger *createfile.LogDir, logFile, smsContent, responseMessage string) (otpRecord *otps.Otp, sendStatus string, err error) {
	otpRecord, err = InsertOTP(ctx, conn, contactNo, auth, from, expiresIn)
	if err != nil {
		return nil, "", err
	}
	smsMessage := strings.Replace(smsContent, "{{OTP}}", fmt.Sprint(otpRecord.Code), -1)
	fmt.Println("sms content is", smsMessage)
	err = sms.SendSMS(appLogger, logFile, contactNo, smsMessage)
	if err != nil {
		otpRecord.SendStatus = 2 // Failed
		sendStatus = "FAILED"
		otpRecord.SendResponse = sendStatus
		err = UpdateOTP(ctx, conn, otpRecord)
		if err != nil {
			appLogger.Error(logFile).Printf("Failed to send SMS to %s: %v", contactNo, err)

			return nil, "", err
		}
		return nil, "", err
	}
	sendStatus = "SUCCESS"
	otpRecord.SendStatus = 1                 // Success
	otpRecord.SendResponse = responseMessage //TO DO:: update sms response
	err = UpdateOTP(ctx, conn, otpRecord)
	if err != nil {
		return nil, "", err
	}

	return otpRecord, sendStatus, nil
}

func OtpExpirationChecking(ctx context.Context, conn sqlx.SqlConn, txId string, currentTime *time.Time) (otpList *otps.Otp, err error) {
	otpModel := otps.NewOtpModel(conn)
	otpList, err = otpModel.FindOneGuid(ctx, txId)
	if err != nil {
		err = fmt.Errorf(errorglobal.InvalidTxid)
		return nil, err
	}

	//check is otp valid to use
	expired := ExpiredChecking(otpList.SendTime.Format(global.DefaultTimeFormat), otpList.ExpiresIn, *currentTime)

	if expired {
		err = fmt.Errorf(errorglobal.InvalidOtp)
		return nil, err
	}
	return otpList, nil
}
