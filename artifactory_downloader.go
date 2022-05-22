package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	auth2 "github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/sirupsen/logrus"
	"os"
)

type rtDownloader struct {
	downloader
	manager *manager.Downloader
}

type rtRepoObject struct {
	Repo string
	File string
}

func getDownloadParams() services.DownloadParams {
	params := services.NewDownloadParams()
	params.Pattern = "repo/*/*.tgz"
	params.Target = "target/path/"
	// Filter the downloaded files by properties.
	params.Props = "key1=val1;key2=val2"
	params.Recursive = true
	params.IncludeDirs = false
	params.Flat = false
	params.Explode = false
	params.Symlink = true
	params.ValidateSymlink = false
	params.Exclusions = []string{"(.*)a.zip"}

	return params
}

func createServiceConfig(rtDetails auth2.ServiceDetails) (config.Config, error) {
	svcCfg, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetContext(context.TODO()).
		Build()

	return svcCfg, err
}

func getRTDetails() auth2.ServiceDetails {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl("")
	//rtDetails.SetApiKey("apikey")
	rtDetails.SetAccessToken("")

	return rtDetails
}

func createRTManager() (artifactory.ArtifactoryServicesManager, error) {
	rtDetails := getRTDetails()
	serviceConfig, err := createServiceConfig(rtDetails)
	if err != nil {
		return nil, fmt.Errorf("failed creating a service config : %v", err)
	}

	rtManager, err := artifactory.New(serviceConfig)

	return rtManager, fmt.Errorf("failed creating a new artifactory client : %v", err)
}

func (downloader rtDownloader) Download(remotePath, outputPath string) (err error) {
	ctx := context.TODO()
	//bucketObject, err := parseS3ResourceFromARN(remotePath)
	//if err != nil {
	//	return
	//}

	file, err := os.Create(outputPath)
	if err != nil {
		logrus.Warnf("Could not create file '%s' for writing", outputPath)
		return
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	//parameters := &s3.GetObjectInput{
	//	Bucket: aws.String(bucketObject.Bucket),
	//	Key:    aws.String(bucketObject.File),
	//}

	//numBytes, err := downloader.manager.Download(ctx, file, parameters)

	//totalDownloaded, totalFailed, err := rtManager.DownloadFiles(params)

	if err != nil {
		logrus.Warnf("Could not download file '%s' from Artifactory repo '%s': %v", bucketObject.File, bucketObject.Bucket, err)
		return
	}
	logrus.Debugf("Downloaded %d bytes from Artifactory", numBytes)

	return
}

//func (downloader rtDownloader) RemoteChecksum(remotePath string) (string, error) {
//	hashRemotePath := fmt.Sprintf("%s.md5", remotePath)
//
//	dir, err := ioutil.TempDir("", "*")
//	if err != nil {
//		logrus.Error("Cannot create temporary file to store remote checksum")
//		return "", err
//	}
//	defer os.RemoveAll(dir)
//	hashFile := filepath.Join(dir, "md5Hash")
//
//	err = downloader.Download(hashRemotePath, hashFile)
//	if err != nil {
//		logrus.Infof("MD5 sum not reachable. %v", err)
//		return "", nil
//	}
//
//	logrus.Infof("Found MD5 sum at: %s", hashRemotePath)
//
//	content, err := ioutil.ReadFile(hashFile)
//	if err != nil {
//		logrus.Warn("Error reading remote checksum")
//		return "", err
//	}
//
//	return string(content), nil
//}
