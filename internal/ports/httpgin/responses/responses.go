package responses

import "github.com/gin-gonic/gin"

func Error(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func Ok() *gin.H {
	return &gin.H{
		"message": "success",
	}
}

func Link(link string) *gin.H {
	return &gin.H{
		"link": link,
	}
}

func Data(data any) *gin.H {
	return &gin.H{
		"data":  data,
		"error": nil,
	}
}
