package form

import "github.com/labstack/echo/v4"

func ProcessFormData[F interface{}](c echo.Context, ptr *F) error {
	if err := c.Bind(ptr); err != nil {
		return err
	}
	if err := c.Validate(ptr); err != nil {
		return err
	}

	return nil
}
