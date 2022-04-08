package cos

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type (
	Cos struct {
		Bucket      string
		AccessKey   string
		SecretKey   string
		Region      string
		Source      string
		Target      string
		StripPrefix string
		Endpoint    string
		Include     string
		Exclude     string
	}

	Ext struct {
		AutoTime   bool
		TimeFormat string
		Debug      bool
		Pause      bool
		Proxy      string
	}

	Plugin struct {
		Cos Cos
		Ext Ext
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	logrus.Info("start upload...")
	p.upload()
	logrus.Info("upload success")
	return nil
}

func (p Plugin) upload() {
	op := &UploadOptions{
		// StorageClass: storageClass,
		// RateLimiting: rateLimiting,
		// PartSize:     partSize,
		ThreadNum: 10,
	}
	if p.Cos.Endpoint == "" {
		p.Cos.Endpoint = fmt.Sprintf("cos.%s.myqcloud.com", p.Cos.Region)
	}
	c := CreateClient(p.Cos.AccessKey, p.Cos.SecretKey, p.Cos.Endpoint, p.Cos.Bucket)
	if p.Ext.AutoTime {
		p.Cos.Target = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.Cos.Target, "/"), time.Now().Format(p.Ext.TimeFormat))
	}
	MultiUpload(c, p.Cos.Source, p.Cos.Bucket, p.Cos.Target, p.Cos.Include, p.Cos.Exclude, op)
}

// 根据函数参数创建客户端
func CreateClient(secretID, secretKey, endpoint, bucketIDName string) *cos.Client {
	return cos.NewClient(CreateURL(bucketIDName, "https", endpoint), &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
}

func GenBucketURL(bucketIDName string, protocol string, endpoint string) string {
	b := fmt.Sprintf("%s://%s.%s", protocol, bucketIDName, endpoint)
	return b
}

func GenServiceURL(protocol string, endpoint string) string {
	s := fmt.Sprintf("%s://%s", protocol, endpoint)
	return s
}

func GenCiURL(bucketIDName string, protocol string, endpoint string) string {
	c := fmt.Sprintf("%s://%s.%s", protocol, bucketIDName, endpoint)
	return c
}

// 根据函数参数生成URL
func CreateURL(idName string, protocol string, endpoint string) *cos.BaseURL {
	b := GenBucketURL(idName, protocol, endpoint)
	s := GenServiceURL(protocol, endpoint)
	c := GenCiURL(idName, protocol, endpoint)

	bucketURL, _ := url.Parse(b)
	serviceURL, _ := url.Parse(s)
	ciURL, _ := url.Parse(c)

	return &cos.BaseURL{
		BucketURL:  bucketURL,
		ServiceURL: serviceURL,
		CIURL:      ciURL,
	}
}

type UploadOptions struct {
	StorageClass string
	RateLimiting float32
	PartSize     int64
	ThreadNum    int
}

func UploadPathFixed(localPath string, targetPath string) (string, string) {
	// eg:~/example/123.txt => cos://bucket/path/123.txt
	// 0. ~/example/123.txt => cos://bucket
	if targetPath == "" {
		pathList := strings.Split(localPath, "/")
		fileName := pathList[len(pathList)-1]
		targetPath = fileName
	}
	// 1. ~/example/123.txt => cos://bucket/path/
	s, err := os.Stat(localPath)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	if s.IsDir() {
		fileNames := strings.Split(localPath, "/")
		fileName := fileNames[len(fileNames)-1]
		targetPath = targetPath + fileName
	}
	// 2. 123.txt => cos://bucket/path/
	if !filepath.IsAbs(localPath) {
		dirPath, err := os.Getwd()
		if err != nil {
			logrus.Fatalln(err)
			os.Exit(1)
		}
		localPath = dirPath + "/" + localPath
	}
	return localPath, targetPath
}
func SingleUpload(c *cos.Client, sourcePath, bucketName, targetPath string, op *UploadOptions) {
	opt := &cos.MultiUploadOptions{
		OptIni: &cos.InitiateMultipartUploadOptions{
			ACLHeaderOptions: &cos.ACLHeaderOptions{
				XCosACL:              "",
				XCosGrantRead:        "",
				XCosGrantWrite:       "",
				XCosGrantFullControl: "",
				XCosGrantReadACP:     "",
				XCosGrantWriteACP:    "",
			},
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				CacheControl:             "",
				ContentDisposition:       "",
				ContentEncoding:          "",
				ContentType:              "",
				ContentMD5:               "",
				ContentLength:            0,
				ContentLanguage:          "",
				Expect:                   "",
				Expires:                  "",
				XCosContentSHA1:          "",
				XCosMetaXXX:              nil,
				XCosStorageClass:         op.StorageClass,
				XCosServerSideEncryption: "",
				XCosSSECustomerAglo:      "",
				XCosSSECustomerKey:       "",
				XCosSSECustomerKeyMD5:    "",
				XOptionHeader:            nil,
				XCosTrafficLimit:         (int)(op.RateLimiting * 1024 * 1024 * 8),
				// Listener:                 &CosListener{},
			},
		},
		PartSize:       op.PartSize,
		ThreadPoolSize: op.ThreadNum,
		CheckPoint:     true,
		// EnableVerification: false,
	}
	sourcePath, targetPath = UploadPathFixed(sourcePath, targetPath)
	logrus.Infof("Upload %s => cos://%s/%s\n", sourcePath, bucketName, targetPath)
	_, _, err := c.Object.Upload(context.Background(), targetPath, sourcePath, opt)
	if err != nil {
		logrus.Errorf("upload: %s,err: %v", sourcePath, err)
	}
}

func MultiUpload(c *cos.Client, sourceDir, bucketName, targetDir, include, exclude string, op *UploadOptions) {
	if sourceDir != "" && (sourceDir[len(sourceDir)-1] != '/' && sourceDir[len(sourceDir)-1] != '\\') {
		sourceDir += "/"
	}
	if targetDir != "" && targetDir[len(targetDir)-1] != '/' {
		targetDir += "/"
	}

	files := GetLocalFilesListRecursive(sourceDir, include, exclude)

	for _, f := range files {
		sourcePath := sourceDir + f
		targetPath := targetDir + f
		SingleUpload(c, sourcePath, bucketName, targetPath, op)
	}
}

func GetLocalFilesListRecursive(sourcePath string, include string, exclude string) (files []string) {
	// bfs遍历文件夹
	var dirs []string
	dirs = append(dirs, sourcePath)
	for len(dirs) > 0 {
		dirName := dirs[0]
		dirs = dirs[1:]

		fileInfos, err := ioutil.ReadDir(dirName)
		if err != nil {
			logrus.Fatalln(err)
			os.Exit(1)
		}

		for _, f := range fileInfos {
			fileName := dirName + "/" + f.Name()
			if f.IsDir() {
				dirs = append(dirs, fileName)
			} else {
				fileName = fileName[len(sourcePath)+1:]
				files = append(files, fileName)
			}
		}
	}

	if len(include) > 0 {
		files = MatchPattern(files, include, true)
	}
	if len(exclude) > 0 {
		files = MatchPattern(files, exclude, false)
	}

	return files
}

func MatchPattern(strs []string, pattern string, include bool) []string {
	res := make([]string, 0)
	re := regexp.MustCompile(pattern)
	for _, s := range strs {
		match := re.MatchString(s)
		if !include {
			match = !match
		}
		if match {
			res = append(res, s)
		}
	}
	return res
}
