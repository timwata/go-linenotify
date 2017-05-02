package linenotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const ENDPOINT = "https://notify-api.line.me/api/notify"

type Client struct {
	token string
}

type Option struct {
	ImageThumbnail   string
	ImageFullsize    string
	ImageFile        string
	StickerPackageId int
	StickerId        int
}

func New(token string) *Client {
	return &Client{
		token: token,
	}
}

func (c Client) Post(msg string, opt *Option) error {
	var err error
	var req *http.Request

	if opt != nil && opt.ImageFile != "" {
		req, err = createMultiPartRequest(msg, opt)
	} else {
		req, err = createFormRequest(msg, opt)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to post message")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusOK {
		return nil
	}

	resp := new(struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	})

	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return errors.Wrap(err, "failed to post message")
	}

	return fmt.Errorf("%d: %s", resp.Status, resp.Message)
}

func createMultiPartRequest(msg string, opt *Option) (*http.Request, error) {
	fd, err := os.Open(opt.ImageFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("imageFile", opt.ImageFile)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, fd)
	if err != nil {
		return nil, err
	}

	fw, err = w.CreateFormField("message")
	if err != nil {
		return nil, err
	}
	fw.Write([]byte(msg))

	if opt.StickerPackageId > 0 {
		fw, _ = w.CreateFormField("stickerPackageId")
		fw.Write([]byte(strconv.Itoa(opt.StickerPackageId)))
	}
	if opt.StickerId > 0 {
		fw, _ = w.CreateFormField("stickerId")
		fw.Write([]byte(strconv.Itoa(opt.StickerId)))
	}

	w.Close()

	req, err := http.NewRequest("POST", ENDPOINT, &b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func createFormRequest(msg string, opt *Option) (*http.Request, error) {
	data := url.Values{"message": {msg}}

	if opt != nil {
		if opt.ImageThumbnail != "" {
			data.Add("imageThumbnail", opt.ImageThumbnail)
		}
		if opt.ImageFullsize != "" {
			data.Add("imageFullsize", opt.ImageFullsize)
		}
		if opt.StickerPackageId > 0 {
			data.Add("stickerPackageId", strconv.Itoa(opt.StickerPackageId))
		}
		if opt.StickerId > 0 {
			data.Add("stickerId", strconv.Itoa(opt.StickerId))
		}
	}

	req, err := http.NewRequest("POST", ENDPOINT, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
