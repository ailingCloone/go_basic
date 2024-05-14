package obsglobal

import (
	"fmt"
	"nrs_customer_module_backend/internal/config"
	"strings"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/zeromicro/go-zero/core/conf"
)

type UploadFileStruct struct {
	FolderPath       string
	UploadedFilename string
	Extension        string
	SourceFile       string
	ContentType      string
}

func ConnectObs() (obsClient *obs.ObsClient, bucket string, err error) {
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, "", err
	}

	var obsKey string = c.OBSCred.ObsKey
	var secret string = c.OBSCred.ObsSecret
	var endPoint string = c.OBSCred.ObsEndPoint
	bucket = c.OBSCred.ObsBucketStaging
	if c.Environment == "prod" {
		bucket = c.OBSCred.ObsBucketProd
	}

	ak := obsKey
	sk := secret

	obsClient, err = obs.New(ak, sk, endPoint /*obs.WithSecurityToken(securityToken)*/)
	if err != nil {
		fmt.Printf("Create obsClient error, errMsg: %s", err.Error())
	}

	return obsClient, bucket, err
}

func UploadFile(q UploadFileStruct) (string, error) {

	obsClient, bucket, err := ConnectObs()

	if err != nil {
		return "", err
	}

	folderPath := q.FolderPath
	uploadedFilename := q.UploadedFilename
	extension := q.Extension
	contentType := q.ContentType

	extension = strings.ReplaceAll(extension, ".", "")

	bucketFolder := "image"
	if strings.Contains(extension, "zip") {
		bucketFolder = "log"
	}

	var key string = fmt.Sprintf("%s/%s/%s", bucketFolder, folderPath, uploadedFilename)
	key = strings.ReplaceAll(key, "//", "/")
	var domain = fmt.Sprintf("https://%s.obs.ap-southeast-3.myhuaweicloud.com/%s", bucket, key)

	input := &obs.PutFileInput{}
	// Specify a bucket name.
	input.Bucket = bucket
	// Specify the object (example/objectname as an example) to upload.

	input.Key = key

	input.ACL = obs.AclPublicRead
	if strings.Contains(extension, "zip") {
		input.ACL = obs.AclPrivate
	}

	input.ContentType = contentType
	input.StorageClass = obs.StorageClassWarm

	// Specify a local file (localfile as an example).
	input.SourceFile = q.SourceFile
	// Perform the file-based upload.
	output, err := obsClient.PutFile(input)
	if err == nil {
		fmt.Printf("Put file(%s) under the bucket(%s) successful!\n", input.Key, input.Bucket)
		fmt.Printf("StorageClass:%s, ETag:%s\n",
			output.StorageClass, output.ETag)
		return domain, nil
	}

	fmt.Printf("Put file(%s) under the bucket(%s) fail!\n", input.Key, input.Bucket)
	if obsError, ok := err.(obs.ObsError); ok {
		fmt.Println("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
		fmt.Println(obsError.Error())
		err = fmt.Errorf(obsError.Error())
	} else {
		fmt.Println("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.")
		fmt.Println(err)
	}

	return "", err

}
