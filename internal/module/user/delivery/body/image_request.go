package body

import "mime/multipart"

type ImageRequest struct {
	Img multipart.File `form:"file"`
}
