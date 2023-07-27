package pdf

import (
	"bytes"
	"image/jpeg"
	"time"

	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

func GetPages(pdf []byte) (pages [][]byte, err error) {
	pool := single_threaded.Init(single_threaded.Config{})
	defer pool.Close()

	instance, err := pool.GetInstance(time.Second * 1)
	if err != nil {
		return
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdf,
	})
	if err != nil {
		return
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return
	}

	for i := 0; i < pageCount.PageCount; i++ {
		pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI: 200,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    i,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		var pageBuf bytes.Buffer
		err = jpeg.Encode(&pageBuf, pageRender.Result.Image, nil)
		if err != nil {
			return nil, err
		}
		pages = append(pages, pageBuf.Bytes())
	}

	return
}
