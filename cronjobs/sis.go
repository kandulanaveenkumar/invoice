package cronjobs

import (
	"io/fs"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/sftp"
	"bitbucket.org/radarventures/forwarder-shipments/services/sis"

	"bitbucket.org/radarventures/forwarder-adapters/utils/sftpclient"
	"bitbucket.org/radarventures/forwarder-shipments/config"

	"go.uber.org/zap"
)

var ftpClient sftpclient.ISFTPClient
var sftpClient sftp.ISFTP
var clientErr error

func InitialiseSIS(ctx *context.Context) {

	if config.Get().EnableSFTP == "true" {
		sftpClient, clientErr = sftp.NewConnection(config.Get().SFTPHost, config.Get().SFTPUser, config.Get().SFTPPass, config.Get().SFTPPort)
		if sftpClient == nil {
			ctx.Log.Error("TERR", zap.Error(clientErr))
			panic(clientErr)
		}
	} else {
		ftpClient, clientErr = sftpclient.NewConnection(config.Get().SISFTPHost, config.Get().SISFTPUser, config.Get().SISFTPPass, config.Get().SISFTPPort)
		if ftpClient == nil {
			ctx.Log.Error("TERR", zap.Error(clientErr))
			panic(clientErr)
		}
	}

}

func Run(ctx *context.Context) {

	ctx.Log.Info("Running the Job")
	ctx.Log.Info("Ticking")

	InitialiseSIS(ctx)

	var fileInfo []fs.FileInfo
	var err error

	if config.Get().EnableSFTP == "true" {
		if sftpClient == nil {
            ctx.Log.Error("SFTP client is nil")
            return
        }
		fileInfo, err = sftpClient.ReadDir(config.Get().SFTPPathIn)
		if err != nil {
			ctx.Log.Error("Toking", zap.Error(err))
		}
	} else {
		fileInfo, err = ftpClient.ReadDir(config.Get().SISFTPPathIn)
		if err != nil {
			ctx.Log.Error("Toking", zap.Error(err))
		}
	}

	for _, file := range fileInfo {

		ctx.Log.Info("File", zap.Any("Fname", file.Name()))

		if file.IsDir() {
			continue
		}

		if config.Get().EnableSFTP == "true" {
			ReadFile(ctx, config.Get().SFTPPathIn+file.Name())
		} else {
			ReadFile(ctx, config.Get().SISFTPPathIn+file.Name())
		}
	}

}

func ReadFile(ctx *context.Context, fname string) {

	if config.Get().EnableSFTP == "true" {
		data, err := sftpClient.Get(fname)
		if err != nil {
			ctx.Log.Error("error getting file name", zap.Error(err))
			return
		}
		if err = CreateNewBookingReq(ctx, data, fname); err == nil {
			sftpClient.Delete(fname)
		}
	} else {
		data, err := ftpClient.Get(fname)
		if err != nil {
			ctx.Log.Error("error getting file name", zap.Error(err))
			return
		}
		if err = CreateNewBookingReq(ctx, data, fname); err == nil {
			ftpClient.Delete(fname)
		}
	}

}

func CreateNewBookingReq(ctx *context.Context, data []byte, file_name string) error {

	return sis.NewSISService().Save(ctx, file_name, data)

}
