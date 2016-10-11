package obj

import (
	"fmt"
	"io"
	"github.com/ncw/swift"
)

type swiftClient struct {
	//TODO: figure out what to store, generally bucket name, etc.
	client Connection
	container string
}

func newSwiftClient(userName string, apiKey string, authUrl string, container string) (*swiftClient, error) {
	//TODO: inputs should be necessary information to start a new session with swift
	//TODO: start a new session and return pointer to swiftClient with context
	client := &swift.Connection{UserName: userName, 
                                    ApiKey: apiKey, 
				    AuthUrl: authUrl}
	
	err := c.Authenticate()
	if err != nil {
		panic(err)
	}

	return &swiftClient{client, container}, nil
}

func (c *swiftClient) Writer(name string) (io.WriteCloser, error) {
	return newBackoffWriteCloser(c, newWriter(c, name)), nil
}

func (c *swiftClient) Walk(name string, fn func(name string) error) error {
	//TODO: implement walk
}

func (c *swiftClient) Reader(name string, offset uint64, size uint64) (io.ReadCloser, error) {
	byteRange := ""
	if size == 0 {
		byteRange = fmt.Sprintf("bytes=%d-", offset)
	} else {
		// we substract 1 one from the right bound because http byte ranges are
		// inclusive rather than clopen
		byteRange = fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)
	}

	// TODO: implement read opject in swift


	getObjectOutput, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(name),
		Range:  aws.String(byteRange),
	})
	if err != nil {
		return nil, err
	}
	return newBackoffReadCloser(c, getObjectOutput.Body), nil
}

func (c *swiftClient) Delete(name string) error {
	//TODO: implement delete object 'name' in swift 
	err := c.client.ObjectDelete(c.container, name)
	return err
}

func (c *swiftClient) Exists(name string) bool {
	//TODO: implement 'name' Exists
	_, _, err := c.client.Object(c.container, name)
	return err == swift.ObjectNotFound
}

func (c *swiftClient) IsRetryable(err error) bool {
	//TODO: check if an operation should be retried based on the error
  	//basically, if error code is service unavaible, retry, otherwise, don't
 	awsErr, ok := err.(awserr.Error) //TODO: figure out what this line does
	if !ok {
		return false
	}
	for _, c := range []string{
		storagegateway.ErrorCodeServiceUnavailable,
		storagegateway.ErrorCodeInternalError,
		storagegateway.ErrorCodeGatewayInternalError,
	} {
		if c == awsErr.Code() {
			return true
		}
	}
	return false
}

func (c *swiftClient) IsIgnorable(err error) bool {
	//TODO: check if there are any ignorable errors (probably not)
	return false
}

func (c *swiftClient) IsNotExist(err error) bool {
	//TODO: are there any non existants errors in swift?
	return err == swift.ContainerNotFound || err ==  swift.ObjectNotFound
}

type swiftWriter struct {
	errChan chan error
	pipe    *io.PipeWriter
}

func newWriter(client *swiftClient, name string) *swiftWriter {
	//TODO: create new swift writer object
	//TODO: fix return type
	obj, headers, err := client.client.Object(client.client, name)
	file, err := client.client.ObjectCreate(client.client.container, name, true, obj.Hash,  obj.ContentType, header) 
	return file
}

func (w *swiftWriter) Write(p []byte) (int, error) {
	//TODO: implement
}

func (w *swiftWriter) Close() error {
	//TODO: implement
}
