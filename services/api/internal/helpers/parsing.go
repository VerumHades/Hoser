package helpers

import (
	"common/pkg/domain/entities/listing"
	"common/pkg/shared"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func ParseSearchQuery(c echo.Context) listing.SearchQuery {
	return listing.SearchQuery{
		Text:      c.QueryParam("q"),
		CPU:       parseNumericRange(c, "cpu_min", "cpu_max"),
		RAMBytes:  parseNumericRange(c, "ram_min", "ram_max"),
		DiskBytes: parseNumericRange(c, "disk_min", "disk_max"),
		Price:     parseNumericRange(c, "price_min", "price_max"),
		CreatedAt: parseDateRange(c, "date_from", "date_to"),
		AuthorID:  parseAuthorID(c.QueryParam("author_id")),
	}
}

func parseNumericRange(c echo.Context, minKey, maxKey string) listing.NumericRange {
	return listing.NumericRange{
		Min: parseOptionalInt64(c.QueryParam(minKey)),
		Max: parseOptionalInt64(c.QueryParam(maxKey)),
	}
}

func parseDateRange(c echo.Context, fromKey, toKey string) listing.DateRange {
	return listing.DateRange{
		From: parseOptionalTime(c.QueryParam(fromKey)),
		To:   parseOptionalTime(c.QueryParam(toKey)),
	}
}

func parseOptionalInt64(value string) *int64 {
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseOptionalTime(value string) *time.Time {
	if value == "" {
		return nil
	}
	// Assuming RFC3339 format (2026-04-18T17:14:16Z)
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseAuthorID(value string) *shared.UserID {
	if value == "" {
		return nil
	}
	id := shared.UserID(value)
	return &id
}
