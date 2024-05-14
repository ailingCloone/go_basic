package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Environment string `json:"environment"`
	DBStaging   struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"dBStaging"`
	DevDBStaging struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"devDBStaging"`
	DBProd struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"dbProd"`

	OBSCred struct {
		ObsKey           string `json:"obsKey"`
		ObsSecret        string `json:"obsSecret"`
		ObsEndPoint      string `json:"obsEndPoint"`
		ObsBucketStaging string `json:"obsBucketStaging"`
		ObsBucketProd    string `json:"obsBucketProd"`
	} `json:"obsCred"`
	TokenValiditySecondsAuthorized int64  `json:"TokenValiditySecondsAuthorized"`
	TokenValiditySecondsLogin      int64  `json:"TokenValiditySecondsLogin"`
	EmailOtpExpire                 int    `json:"EmailOtpExpire"`
	SMSOtpExpire                   int    `json:"SMSOtpExpire"`
	UserImageUrl                   string `json:"UserImageUrl"`

	SMTPEmail struct {
		Username string
		Password string
		Port     int
		From     string
		Host     string
	}
}
