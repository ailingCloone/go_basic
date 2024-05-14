package uploadfile

import (
	"fmt"
	"io"
	"net/http"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/obsglobal"

	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CommonStruct struct {
	W            http.ResponseWriter
	R            *http.Request
	AppLogger    *createfile.LogDir
	LogFilename  string
	SelectedFile string
	ApiName      string
	FormFile     string
}

type CommonResStruct struct {
	Status     string
	Err        error
	StatusCode int
	Filename   string
	PathName   string
}

type UploadOBSAndRemoveLocalStruct struct {
	SelectedFile string
	RandomName   string
	Extension    string
	Name         string
	ContentType  string
	Err          error
	AppLogger    *createfile.LogDir
	LogFilename  string
	ApiName      string
}

func ImagePath() (path string, err error) {
	path = filepath.Join("image/")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func Common(q CommonStruct) (Result CommonResStruct) {
	r := q.R
	appLogger := q.AppLogger
	logFilename := q.LogFilename
	selectedFile := q.SelectedFile
	apiName := q.ApiName
	randomName := ""
	file_name := ""
	contentType := ""
	path, err := ImagePath()
	if selectedFile == "" || err != nil {
		Result = CommonResStruct{
			Status:     "error",
			Err:        fmt.Errorf("some thing wrong"),
			StatusCode: http.StatusInternalServerError,
		}
		return Result
	}

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile(q.FormFile)
	if err != nil {
		switch err {
		case http.ErrMissingFile:
			Result = CommonResStruct{
				Status:     "errorF",
				Err:        fmt.Errorf(err.Error()),
				StatusCode: http.StatusInternalServerError,
			}
		default:
			Result = CommonResStruct{
				Status:     "error",
				Err:        fmt.Errorf(err.Error()),
				StatusCode: http.StatusInternalServerError,
			}
		}
		return Result
	}
	fileName := handler.Filename
	header := handler.Header
	var extension = filepath.Ext(fileName)
	if header["Content-Type"] != nil {
		contentType = header["Content-Type"][0]
	}
	defer file.Close()
	t := strconv.FormatInt(time.Now().Unix(), 10)
	file_name = strings.TrimSuffix(fileName, extension)
	file_name = strings.ReplaceAll(file_name, " ", "")
	randomName = file_name + "_" + string(t) + extension
	name := filepath.Join(path, filepath.Base(randomName))
	// Create file
	dst, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0766)
	if err != nil {
		appLogger.Error(logFilename).Println("[Common] OpenFile: ", err)
		Result = CommonResStruct{
			Status:     "error",
			Err:        fmt.Errorf(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
		return Result
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		appLogger.Error(logFilename).Println("[Common] Copy: ", err)
		Result = CommonResStruct{
			Status:     "error",
			Err:        fmt.Errorf(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
		return Result
	}

	param := UploadOBSAndRemoveLocalStruct{
		SelectedFile: selectedFile,
		RandomName:   randomName,
		Extension:    extension,
		Name:         name,
		ContentType:  contentType,
		Err:          err,
		AppLogger:    appLogger,
		LogFilename:  logFilename,
		ApiName:      apiName,
	}
	Result = UploadOBSAndRemoveLocal(param)

	return Result
}

func UploadOBSAndRemoveLocal(q UploadOBSAndRemoveLocalStruct) (Result CommonResStruct) {
	uploadFileParam := obsglobal.UploadFileStruct{
		FolderPath:       q.SelectedFile,
		UploadedFilename: q.RandomName,
		Extension:        q.Extension,
		SourceFile:       q.Name,
		ContentType:      q.ContentType,
	}
	q.RandomName, q.Err = obsglobal.UploadFile(uploadFileParam)
	if q.Err != nil {
		q.AppLogger.Error(q.LogFilename).Println("[Common] OpenFile: ", q.Err)
		Result = CommonResStruct{
			Status:     "error",
			Err:        fmt.Errorf(q.Err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
		return Result
	}
	q.Err = os.Remove(q.Name)
	if q.Err != nil {
		q.AppLogger.Error(q.LogFilename).Println("[Common] Remove: ", q.Err)
	}
	Result = CommonResStruct{
		Status:   "success",
		Filename: q.RandomName,
		PathName: q.ApiName,
	}

	return Result
}

func UploadFile(w http.ResponseWriter, r *http.Request, appLogger *createfile.LogDir, logFilename string) (Result *http.Request, err error) {
	selectedFile := ""
	apiName := ""
	isNeedPathName := true

	switch r.FormValue("path") {
	case "member_card":
		selectedFile = "member_card/"
	case "side_menu":
		selectedFile = "side_menu/"
	case "splash":
		selectedFile = "splash/"
	case "ui":
		selectedFile = "ui/"
	default:
		selectedFile = ""
	}

	commonParam := CommonStruct{
		W:            w,
		R:            r,
		AppLogger:    appLogger,
		LogFilename:  logFilename,
		SelectedFile: selectedFile,
		ApiName:      apiName,
		FormFile:     "file",
	}

	result := Common(commonParam)

	if result.Status != "success" {
		http.Error(w, fmt.Sprint(result.Err), result.StatusCode)
		r.Form.Add("status", result.Status)
		return r, result.Err
	}

	r.Form.Add("status", result.Status)
	r.Form.Add("filename", result.Filename)
	if isNeedPathName {
		r.Form.Add("path_name", result.PathName)
	}
	return r, nil
}
