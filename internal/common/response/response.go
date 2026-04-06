package response

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type ErrorResponse struct {
	Code    apperror.ErrorCode `json:"code,omitempty"`
	Details any                `json:"details,omitempty"`
}

type Meta struct {
	Total      int64 `json:"total"`
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalPages int32 `json:"total_pages"`
}

type Links struct {
	Self string  `json:"self"`
	Next *string `json:"next,omitempty"`
	Prev *string `json:"prev,omitempty"`
}

type APIResponse struct {
	Message   string         `json:"message,omitempty"`
	Data      any            `json:"data,omitempty"`
	Meta      *Meta          `json:"meta,omitempty"`
	Links     *Links         `json:"links,omitempty"`
	Error     *ErrorResponse `json:"error,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
}

func Success(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, APIResponse{
		Data:      data,
		RequestID: middleware.GetRequestID(c),
	})
}

func Paginated(c *gin.Context, statusCode int, data any, meta *Meta, links *Links) {
	c.JSON(statusCode, APIResponse{
		Data:      data,
		Meta:      meta,
		Links:     links,
		RequestID: middleware.GetRequestID(c),
	})
}

func BuildPaginationLinks(c *gin.Context, page, limit int, totalPages int32, allowedKeys map[string]struct{}) *Links {
	selfURL := buildPageURL(c, page, limit, allowedKeys)

	links := &Links{
		Self: selfURL,
	}

	if page > 1 && totalPages > 0 {
		prevURL := buildPageURL(c, page-1, limit, allowedKeys)
		links.Prev = &prevURL
	}

	if totalPages > 0 && int32(page) < totalPages {
		nextURL := buildPageURL(c, page+1, limit, allowedKeys)
		links.Next = &nextURL
	}

	return links
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.Error); ok {
		apperror.SetErrorContext(c, appErr)

		errRes := APIResponse{
			Message:   appErr.Message,
			RequestID: middleware.GetRequestID(c),
			Error: &ErrorResponse{
				Code: appErr.Code,
			},
		}

		if appErr.Code == apperror.ErrCodeValidation && appErr.Details != nil {
			errRes.Error.Details = appErr.Details
		}

		c.JSON(apperror.HttpStatusCodeFromCode(appErr.Code), errRes)
		return
	}

	internalErr := &apperror.Error{
		Status:  http.StatusInternalServerError,
		Code:    apperror.ErrCodeInternal,
		Message: "internal server error",
		Err:     err,
	}

	apperror.SetErrorContext(c, internalErr)

	c.JSON(http.StatusInternalServerError, APIResponse{
		Message:   internalErr.Message,
		RequestID: middleware.GetRequestID(c),
		Error: &ErrorResponse{
			Code: internalErr.Code,
		},
	})
}

func buildPageURL(c *gin.Context, page, limit int, allowedKeys map[string]struct{}) string {
	raw := c.Request.URL.Query()
	query := url.Values{}

	for k, vals := range raw {
		if _, ok := allowedKeys[k]; !ok {
			continue
		}
		for _, v := range vals {
			query.Add(k, v)
		}
	}

	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))

	path := c.Request.URL.Path
	encoded := query.Encode()
	if encoded == "" {
		return path
	}

	return path + "?" + encoded
}
