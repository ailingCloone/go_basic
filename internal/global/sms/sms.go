package sms

import (
	"nrs_customer_module_backend/internal/global/createfile"
)

func SendSMS(appLogger *createfile.LogDir, filename string, contact, smsContent string) error {
	// fmt.Println("Send email: ", " INFO:", info, "RABBITMQ: ", rmRes)

	// params := "user=" + url.QueryEscape("senhengotp") + "&" +
	// 	"pass=" + url.QueryEscape("Senheng@188") + "&" +
	// 	"type=" + url.QueryEscape("0") + "&" +
	// 	"to=" + url.QueryEscape(contact) + "&" +
	// 	"from=" + url.QueryEscape("Senheng") + "&" +
	// 	"text=" + url.QueryEscape(smsContent) + "&" +
	// 	"servid=" + url.QueryEscape("MES01") + "&" +
	// 	"title=" + url.QueryEscape("RABBITMQ") + "&" +
	// 	"detail=" + url.QueryEscape("1")

	// path := fmt.Sprintf(config.EtrackerUrl+"%s", params)

	// log.Println("Info Start: ", path)

	// res, err := http.Get(path)
	// if err != nil {
	// 	appLogger.Error(filename).Println(err, "Failed to Get SMS http")
	// }

	// _, err = ioutil.ReadAll(res.Body)

	// defer res.Body.Close()

	// if err != nil {
	// 	appLogger.Error(filename).Println("[x] Info: ", path, "[x] SMS ioutil.ReadAll ", err)
	// 	return
	// }
	return nil
}
