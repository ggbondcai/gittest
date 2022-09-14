package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

func main() {
	ctx := context.Background()
	uploadCommand(ctx)
	uploadCmd := uploadCommand(ctx)
	fileName := "a.txt" // txt文件路径
	data, errRead := ioutil.ReadFile(fileName)
	if errRead != nil {
		fmt.Println("文件读取失败！")
	}
	uploadCmd.Stdin = bytes.NewReader(data)
	err := RunInSequence(uploadCmd)
	if err != nil {
		fmt.Printf(err.Error())
	}

}
func uploadCommand(ctx context.Context) *exec.Cmd {
	uploadArgs := []string{
		"put",
		"--md5",
		fmt.Sprintf("--storage=s3"),
		fmt.Sprintf("--s3-endpoint=https://proxy.uss.s3.test.sz.shopee.io"),
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
		println("end %s without error", cmd.Path)
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
