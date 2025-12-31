package util

import (
	"common/pkg/shared"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NextBatchFunc defines the function signature to fetch the next batch of items.
type GetBatchFunc[T any, CursorType any] func(ctx context.Context, request shared.BatchRequest[CursorType]) (items []T, nextCursor CursorType, err error)
type ViewConversionFunction[ElementType any, ElementViewType any] func(elements []ElementType) (views []ElementViewType)

func HandleBatchRequest[ElementType any, ElementViewType any, CursorType any](
	c echo.Context,
	batchFunction GetBatchFunc[ElementType, CursorType],
	viewConversionFunction ViewConversionFunction[ElementType, ElementViewType],
) error {
	ctx := c.Request().Context()
	batchSize := 50

	var cursor CursorType
	if encoded := c.QueryParam("cursor"); encoded != "" {
		if decoded, err := DecodeCursor[CursorType](encoded); err == nil {
			cursor = decoded
		}
	}

	batchRequest := shared.BatchRequest[CursorType]{
		Cursor:       cursor,
		MaxBatchSize: batchSize,
	}

	items, nextCursor, err := batchFunction(ctx, batchRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch library")
	}

	apiListings := viewConversionFunction(items)
	encodedCursor, _ := EncodeCursor(nextCursor)

	resp := PaginatedResponse[ElementViewType, CursorType]{
		Items:  apiListings,
		Cursor: encodedCursor,
	}

	return c.JSON(http.StatusOK, resp)
}

type ConversionFunction[InputType any, OutputType any] func(elemenet InputType) (convertedElement OutputType)

func MapList[InputType any, OutputType any](
	elements []InputType,
	convert ConversionFunction[InputType, OutputType],
) (convertedElements []OutputType) {
	outputElements := make([]OutputType, len(elements))
	for i, element := range elements {
		outputElements[i] = convert(element)
	}
	return outputElements
}
