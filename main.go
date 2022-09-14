package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
)

func main() {
	//fileName2 := "a.txt" // txt文件路径
	//data, err_read := ioutil.ReadFile(fileName2)
	//if err_read != nil {
	//	fmt.Println("读取失败")
	//}
	bucket := "appinfraentrytask"
	host := "http://proxy.uss.s3.test.sz.shopee.io/" //s3 domain
	ak := "52633284"                                 //appid
	sk := "afsZqzjLWuzftIwKldTldtkoacMbZRil"         //secret
	fileName := "backup"                             //upload file name
	s3Config := &aws.Config{
		Credentials:          credentials.NewStaticCredentials(ak, sk, ""),
		Endpoint:             aws.String(host),
		Region:               aws.String("default"),
		DisableSSL:           aws.Bool(true),
		S3ForcePathStyle:     aws.Bool(true),
		S3Disable100Continue: aws.Bool(true),
	}
	sess := session.New(s3Config)
	//backupCmd := backupCommand(ctx)
	//backupStdout, err := backupCmd.StdoutPipe()
	//if err != nil {
	//	fmt.Printf("backup 失败1")
	//	return
	//}
	//err = RunInSequence(backupCommand(ctx))
	//if err != nil {
	//	fmt.Println("backup 失败2")
	//}
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   GetIoreader(),
	})
	if err != nil {
		// Print the error and exit.
		fmt.Printf("Unable to upload %q to %q, %v", fileName, bucket, err)
		return
	}
	fmt.Printf("Upload %q to %q success", fileName, bucket)
}
func backupCommand(ctx context.Context) *exec.Cmd {
	backupArgs := []string{
		"--backup",
		"--stream=xbstream",
		fmt.Sprintf("--host=localhost"),
		fmt.Sprintf("--port=3306"),
		fmt.Sprintf("--user=root"),
		fmt.Sprintf("--password=123456"),
		fmt.Sprintf("--target-dir=/temp"),
		fmt.Sprintf("--extra-lsndir=/temp"),
		fmt.Sprintf("--parallel=8"),
		"--compress",
		"--compress-threads=8",
		"--read-buffer-size=2G",
		"--encrypt-chunk-size=2G",
	}
	return exec.CommandContext(ctx, "xtrabackup", backupArgs...)
}
func uploadCommand(ctx context.Context) *exec.Cmd {
	uploadArgs := []string{
		"put",
		"--md5",
		fmt.Sprintf("--storage=s3"),
		fmt.Sprintf("--s3-endpoint=proxy.uss.s3.test.sz.shopee.io"),
		fmt.Sprintf("--s3-access-key=52633284"),
		fmt.Sprintf("--s3-secret-key=afsZqzjLWuzftIwKldTldtkoacMbZRil"),
		fmt.Sprintf("--s3-bucket=appinfraentrytask"),
		fmt.Sprintf("--parallel=8"),
		"caimingyang18",
	}

	return exec.CommandContext(ctx, "xbcloud", uploadArgs...)
}
func RunInSequence(cmds ...*exec.Cmd) error {
	for _, cmd := range cmds {
		var stdErr bytes.Buffer
		cmd.Stderr = &stdErr

		cmdString := PrintCmd(cmd)
		println("start %s", cmdString)
		if err := cmd.Run(); err != nil {
			fmt.Printf("execute command failed %v, command: %s, stdErr: %s", err, cmdString, stdErr.String())
			return err
		}
		println("end  without error", cmd.Path)
	}
	return nil
}
func PrintCmd(cmd *exec.Cmd) string {
	s := cmd.String()
	s = pwRegex.ReplaceAllString(s, "$1***$3")
	s = secretRegex.ReplaceAllString(s, "$1***$3")
	return s
}

var (
	pwRegex     = regexp.MustCompile(`(.*password[^=]*=)(\w+)(\s?.*)$`)
	secretRegex = regexp.MustCompile(`(.*secret[^=]*=)(\w+)(\s?.*)$`)
)

func GetIoreader() io.Reader {
	ctx := context.Background()
	backupCmd := backupCommand(ctx)
	backupStdout, err := backupCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("backup 失败1")
		return nil
	}
	err = RunInSequence(backupCommand(ctx))
	fmt.Println(ioutil.ReadAll(backupStdout))
	return backupStdout
}
