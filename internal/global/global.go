package global

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateGuid() (Result string) {
	newUUID := uuid.New()
	return newUUID.String()
}

const DefaultTimeFormat = "2006-01-02 15:04:05"

func TimeInSingapore() (*time.Time, error) {
	format := DefaultTimeFormat
	// Load Singapore time zone
	singapore, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return nil, err
	}

	// Get current time in Singapore time zone
	now := time.Now().In(singapore)

	// Format the time using the provided format string
	formattedTimeStr := now.Format(format)

	// Parse the formatted string back into a time.Time object
	formattedTime, err := time.Parse(format, formattedTimeStr)
	if err != nil {
		return nil, err
	}

	return &formattedTime, nil
}

// extractFieldName parses the error message to extract the field name
func ExtractFieldName(errMsg string) string {
	// Assuming the error message format is "field 'fieldname' is not set"
	// Split the error message by single quotes to extract the field name
	parts := strings.Split(errMsg, "\"")

	if len(parts) == 3 {
		return parts[1]
	}
	return ""
}

// get current date time
func GetCurrentTime() string {
	currentTime := time.Now()
	currentTime2 := currentTime.Format("2006-01-02 15:04:05")
	return currentTime2
}

// convert string to int64
func ConvertStringToInt64(value string) (i int64, err error) {

	i, err = strconv.ParseInt(value, 10, 64)

	return i, err

}

// If no layout is provided, it uses the DefaultTimeFormat.
func FormatTime(t time.Time, layout ...string) string {
	var l string
	if len(layout) > 0 {
		l = layout[0] //call function with format e.g. global.FormatTime(currentTime, "Monday, Jan _2 2006")
	} else {
		l = DefaultTimeFormat
	}
	return t.Format(l)
}

func GenerateSha256(pass string) (sha256PassStr string) {
	hash := sha256.New()
	hash.Write([]byte(pass))
	sha256Pass := hash.Sum(nil)
	sha256PassStr = hex.EncodeToString(sha256Pass)

	return sha256PassStr
}

func GenerateAccessToken() string {

	accessToken := GenerateRefreshToken()

	accessToken = strings.ReplaceAll(accessToken, "-", "")

	return accessToken

}

// GenerateRefreshToken generates a universally unique refresh token
func GenerateRefreshToken() string {
	// Generate a UUID (version 4) as a refresh token
	refreshToken := uuid.New().String()
	return refreshToken
}

func GetFirstAndLastDayOfMonth() (firstDay, lastDay string) {
	t, err := TimeInSingapore()
	if err != nil {
		return
	}
	firstDay = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
	lastDay = time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
	return firstDay, lastDay
}

func PasswordValidation(password string, confirmPassword string) bool {
	if password == confirmPassword {
		return true
	} else {
		return false
	}
}

// MaskPhoneNumber masks most of the digits of a phone number and reveals the last few digits.
func MaskPhoneNumber(phoneNumber string) string {
	// Check if the phone number is valid (at least 4 digits)
	if len(phoneNumber) < 4 {
		return "Invalid phone number"
	}

	// Extract the last few digits to reveal
	lastDigits := phoneNumber[len(phoneNumber)-2:] // Show last two digits (adjust as needed)

	// Determine the number of 'x' characters to use for masking
	maskLength := len(phoneNumber) - len(lastDigits)

	// Generate the masked phone number
	maskedPhoneNumber := ""
	for i := 0; i < maskLength; i++ {
		maskedPhoneNumber += "x"
	}
	maskedPhoneNumber += lastDigits

	return maskedPhoneNumber
}

// func UserAccessPermission(conn sqlx.SqlConn, ctx context.Context, StaffId int64, referTable string) (*[]user_access_setting.UserAccessSettings, error) {
// 	userAccessSettingModel := user_access_setting.NewUserAccessSettingModel(conn)

// 	userAccessSetting, err := userAccessSettingModel.FindAll(ctx, StaffId, referTable)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return userAccessSetting, nil
// }

func FormatMDYTime(sendTime time.Time) string {
	sendTimeStr := FormatTime(sendTime, "01-02-2006 15:04:05")
	return sendTimeStr
}

type DataFilter struct {
	Title string `json:"title"`
	Code  int    `json:"code"`
}

func FilterDay(page string) (filterDay []DataFilter) {
	latest := DataFilter{
		Title: "Latest",
		Code:  0,
	}
	days30 := DataFilter{
		Title: "30 Days",
		Code:  30,
	}
	days60 := DataFilter{
		Title: "60 Days",
		Code:  60,
	}

	switch page {
	case "app_registration_list_get_list":
		filterDay = append(filterDay, latest)
		filterDay = append(filterDay, days30)
		filterDay = append(filterDay, days60)
	}

	return filterDay

}

func CheckIsEmpty(value string) bool {
	if value == "" || value == "null" || value == "nil" {
		return true
	}

	return false
}
