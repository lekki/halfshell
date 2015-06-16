// Copyright (c) 2014 Oyster
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package halfshell

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"fmt"
)

var S3Service *s3.S3 = s3.New(&aws.Config{Region: "us-east-1", Logger: os.Stdout, LogLevel: 0})

const (
	ImageSourceTypeS3 ImageSourceType = "s3"
)

type S3ImageSource struct {
	Config *SourceConfig
	Logger *Logger
}

func NewS3ImageSourceWithConfig(config *SourceConfig) ImageSource {
	fmt.Printf("%v", config)

	return &S3ImageSource{
		Config: config,
		Logger: NewLogger("source.s3.%s", config.Name),
	}
}

func (s *S3ImageSource) GetImage(request *ImageSourceOptions) (*Image, error) {


	if s.Config.LocalCache {

		content, err := CacheRead(request.Path)
		if content != nil && err == nil {
			return content, nil
		}
	}

	params := &s3.GetObjectInput{Bucket: aws.String(s.Config.S3Bucket), Key: aws.String(request.Path)}
	resp, err := S3Service.GetObject(params)

	if awsErr, ok := err.(awserr.Error); ok {
		// Generic AWS Error with Code, Message, and original error (if any)
		fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		if reqErr, ok := err.(awserr.RequestFailure); ok {
					// A service error occurred
			s.Logger.Warnf("Error on fetching %v %v %v %v %v", request.Path, reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			return nil, err
		} else {
			s.Logger.Warnf("%v %v", err.Error(), request.Path)
			return nil, err
		}
	}

	image, err := NewImageFromBuffer(resp.Body)
	if err != nil {
		s.Logger.Warnf("Unable to create image from response  (url=%v)", request.Path)
		return nil, err
	}

	if s.Config.LocalCache {
		CacheWrite(request.Path, image)
	}

	return image, nil
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSourceWithConfig)
}
