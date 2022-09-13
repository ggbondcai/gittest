package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	host := "dba-sg.uss.shopee.io/"          //s3 domain
	ak := "23498786"                         //appid
	sk := "hFqtMItPAkGSBuafrLGdUKmohoKDdVVV" //secret

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		Endpoint:         aws.String(host),
		Region:           aws.String("default"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess := session.New(s3Config)
	svc := s3.New(sess)
	bucketList, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("get bucket list fail: %v", err)
	}
	fmt.Println(bucketList)
}
