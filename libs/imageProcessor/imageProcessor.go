package imageProcessor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"io"
	"io/ioutil"
	"os/exec"
)

const maxAttachmentSize = 20 * 1024 * 1024
const folder = "./tmp/"

type Client interface {
	NewImageName() (fullName string, name string)
	Validate(r io.Reader) error
	Resize(srcImage string, maxWidth int, maxHeight int) (dstImage string, err error)
}

type defaultClient struct {
}

func GetDefault() Client {
	return defaultClient{}
}

func (dc defaultClient) NewImageName() (fullName string, name string) {
	rand.Seed(time.Now().UnixNano())
	name = uuid.New().String() + ".jpg"
	fullName = folder + name
	return
}

func (dc defaultClient) Validate(fileHandle io.Reader) error {

	reader := io.LimitReader(fileHandle, maxAttachmentSize)
	attachmentFileBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	fType, err := filetype.Image(attachmentFileBytes)

	if err != nil {
		return err
	}

	if fType.MIME.Value != "image/jpeg" {
		return errors.New("wrong image type")
	}

	return nil
}

// Resize function scales source image and stores it to file system for caching purposes. This prevents scaling to happen
// more than once for same image and maxWidth and maxHeight parameters.
func (dc defaultClient) Resize(srcImage string, maxWidth int, maxHeight int) (dstImage string, err error) {

	// We use imagemagick resize parameter as part of destination file name.
	// Examples of different resize combinations:
	// "{maxWidth}" if only maxWidth is set
	// "x{maxHeight}" if only maxHeight is set
	// "{maxWidth}x{maxHeight}" if both are set
	// e.g. filename-100x100.jpg is a picture having max height and width set to 100

	var scaling = ""

	if maxWidth > 0 {
		scaling += fmt.Sprintf("%d", maxWidth)
	}
	if maxHeight > 0 {
		scaling += fmt.Sprintf("x%d", maxHeight)
	}

	srcImageWithPath := folder + srcImage
	if _, err = os.Stat(srcImageWithPath); err != nil {
		// source image is not found
		return
	}

	if scaling != "" {
		nameSplit := strings.Split(srcImage, ".")
		dstImage = fmt.Sprintf("%s%s-%s.jpg", folder, nameSplit[0], scaling)
	} else {
		// No scaling needed so source and destination is the same
		dstImage = srcImageWithPath
	}

	if _, err = os.Stat(dstImage); err == nil {
		// File already exists
		return
	}

	var fDest *os.File
	if fDest, err = os.Create(dstImage); err != nil {
		return
	}
	defer fDest.Close()

	var fSrc *os.File
	if fSrc, err = os.Open(srcImageWithPath); err != nil {
		return
	}
	defer fSrc.Close()

	reader := io.LimitReader(fSrc, maxAttachmentSize)
	attachmentFileBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return
	}

	resizeParameter := []string{"-", "-resize", scaling, "-"}

	cmd := exec.Command("convert", resizeParameter...)
	cmd.Stdout = bufio.NewWriter(fDest)
	cmd.Stdin = bytes.NewBuffer(attachmentFileBytes)

	err = cmd.Run()

	return
}
