package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	contentTypeHeader = "Content-Type"
	textPlain         = "text/plain"
	gauge             = "gauge"
	counter           = "counter"
)

func ValidateUpdateURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodPost {
			//err := fmt.Errorf("Invalid http method! Expected %s, Got %s. ", http.MethodPost, ctx.Request.Method)
			ctx.AbortWithStatus(http.StatusNotFound)
			//log.Println(err)
			return
		}

		//if ctx.GetHeader(contentTypeHeader) != textPlain {
		//	err := fmt.Errorf("Invalid content type! Expected %s, Got %s. ", textPlain, ctx.GetHeader(contentTypeHeader))
		//	ctx.AbortWithStatus(http.StatusBadRequest)
		//	log.Println(err)
		//	return
		//}

		path := ctx.Request.URL.RequestURI()
		path = strings.Trim(path, "/")
		params := strings.Split(path, "/")

		if len(params) == 2 && params[1] != gauge && params[1] != counter {
			//err := fmt.Errorf("Invalid metric type and url path!\n"+
			//	"Expected url params length = %d, Got %d. ", 4, len(params))
			ctx.AbortWithStatus(http.StatusBadRequest)
			//log.Println(err)
			return
		}

		if len(params) != 4 {
			//err := fmt.Errorf("Invalid url path! Expected url params length = %d, Got %d. ", 4, len(params))
			ctx.AbortWithStatus(http.StatusNotFound)
			//log.Println(err)
			return
		}

		switch params[1] {
		case gauge:
			_, err := strconv.ParseFloat(params[3], 64)
			if err != nil {
				//err = fmt.Errorf("Invalid gauge value = %v ! Expected float64. ", params[3])
				ctx.AbortWithStatus(http.StatusBadRequest)
				//log.Println(err)
				return
			}
		case counter:
			_, err := strconv.ParseInt(params[3], 10, 64)
			if err != nil {
				//err = fmt.Errorf("Invalid counter value = %v ! Expected int64. ", params[3])
				ctx.AbortWithStatus(http.StatusBadRequest)
				//log.Println(err)
				return
			}
		default:
			//err := errors.New("Invalid metric type! ")
			ctx.AbortWithStatus(http.StatusBadRequest)
			//log.Println(err)
			return
		}

		ctx.Next()
	}
}
