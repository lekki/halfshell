package halfshell

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/youtube/vitess/go/pools"
	"golang.org/x/net/context"
	"time"
	"fmt"
)


//var s3Config  aws.Config =

var S3Pool   *pools.ResourcePool


type S3Service struct {
	s *s3.S3
}

func (s3Service *S3Service) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return s3Service.s.GetObject(input)
}

func (s *S3Service) Close() {

}

func InitS3Pool() {
	S3Pool = NewS3Pool(30, 100, 0)
}

func S3Factory() pools.Factory {
	return func() (pools.Resource, error) {
		s3Service := new(S3Service)
		s3Service.s = s3.New(&aws.Config{Region: aws.String("us-east-1"), LogLevel: 0})

		return s3Service, nil
	}
}

func NewS3Pool(capacity int, maxCapacity int, idleTimout time.Duration) *pools.ResourcePool {
	return pools.NewResourcePool(S3Factory(), capacity, maxCapacity, idleTimout)
}

func CloseS3Pool() {
	S3Pool.Close()
}

func S3GetFromPool() (*S3Service, error) {

	start := time.Now()
	ctx := context.TODO()
	resource, err := S3Pool.Get(ctx)
	fmt.Printf("getting connection took: %v \n", time.Since(start))

	return resource.(*S3Service), err
}

func S3ReleaseToPool(s3 *S3Service) {
	S3Pool.Put(s3)
}
