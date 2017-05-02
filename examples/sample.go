package main

import (
	"fmt"
	"os"

	".."
)

func main() {
	cli := linenotify.New(os.Getenv("TOKEN"))

	// If you specified imageThumbnail, imageFullsize and imageFile, imageFile takes precedence.
	err := cli.Post("Hello, World!", &linenotify.Option{
		ImageFile:        "local.png",
		StickerPackageId: 1,
		StickerId:        1,
	})
	if err != nil {
		fmt.Println(err)
	}

	// sticker list: https://devdocs.line.me/files/sticker_list.pdf
	err := cli.Post("Hello, World!", &linenotify.Option{
		ImageThumbnail:   "https://upload.wikimedia.org/wikipedia/commons/thumb/4/41/LINE_logo.svg/300px-LINE_logo.svg.png",
		ImageFullsize:    "https://upload.wikimedia.org/wikipedia/commons/thumb/4/41/LINE_logo.svg/300px-LINE_logo.svg.png",
		StickerPackageId: 1,
		StickerId:        1,
	})
	if err != nil {
		fmt.Println(err)
	}
}
