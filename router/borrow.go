package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostBorrowingEquipment POST /items/:id/borrowing
func PostBorrowingEquipment(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostBorrowingEquipmentReturn POST /items/:id/borrowing/return
func PostBorrowingEquipmentReturn(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostBorrowings POST /items/:id/owners/:ownershipid/borrowing
func PostBorrowings(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// GetBorrowingsId GET /items/:id/owners/:ownershipid/borrowing/:borrowingid
func GetBorrowingsId(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostBorrowingsReply POST /items/:id/owners/:ownershipid/borrowing/:borrowingid/reply
func PostBorrowingsReply(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}

// PostBorrowingsReturn POST /items/:id/owners/:ownershipid/borrowing/:borrowingid/return
func PostBorrowingsReturn(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not Implemented")
}
