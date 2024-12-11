package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chud-lori/go-echo-temp/adapters/transport"
	"github.com/chud-lori/go-echo-temp/domain/ports"
    "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	ports.UserService
}

func NewUserController(service ports.UserService) *UserController {
	return &UserController{UserService: service}
}

func GetPayload(ctx echo.Context, result interface{}) error {
	if err := ctx.Bind(result); err != nil {
        return err
    }
    return nil
}

func WriteResponse(writer http.ResponseWriter, response interface{}, httpCode int64) {

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(int(httpCode))
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)

	if err != nil {
		panic(err)
	}
}

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

func (controller *UserController) Create(ctx echo.Context) error {
	userRequest := transport.UserRequest{}

    if err := GetPayload(ctx, &userRequest); err != nil {
        return ctx.JSON(http.StatusBadRequest, WebResponse{
            Message: "Invalid request Payload",
            Status: 0,
            Data: nil,
        })
    }

	userResponse, err := controller.UserService.Save(ctx.Request().Context(), &userRequest)

	if err != nil {
		fmt.Println("Error create controller")
		panic(err)
	}

	response := WebResponse{
		Message: "success save user",
		Status:  1,
		Data:    userResponse,
	}

    return ctx.JSON(http.StatusCreated, response)
}

func (controller *UserController) Update(ctx echo.Context) error {
	userRequest := transport.UserRequest{}

    if err := GetPayload(ctx, &userRequest); err != nil {
        return ctx.JSON(http.StatusBadRequest, WebResponse{
            Message: "Invalid request Payload",
            Status: 0,
            Data: nil,
        })
    }

	userResponse, err := controller.UserService.Update(ctx.Request().Context(), &userRequest)

	if err != nil {
		fmt.Println("Error update controller")
        return ctx.JSON(http.StatusInternalServerError, WebResponse{
			Message: "Failed to update user",
			Status:  0,
			Data:    nil,
		})
	}

	response := WebResponse{
		Message: "success update user",
		Status:  1,
		Data:    userResponse,
	}

    return ctx.JSON(http.StatusOK, response)

}

func (controller *UserController) Delete(ctx echo.Context) error {
	userId := ctx.Param("userId")

	err := controller.UserService.Delete(ctx.Request().Context(), userId)

	if err != nil {
		fmt.Println("Error delete controller")
        return ctx.JSON(http.StatusInternalServerError, WebResponse{
            Message: "Failed Delete user",
            Status: 0,
            Data: nil,
        })
	}

	response := WebResponse{
		Message: "success delete user",
		Status:  1,
		Data:    "sucess",
	}

    return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) FindById(ctx echo.Context) error {
	userId := ctx.Param("userId")

	user, err := c.UserService.FindById(ctx.Request().Context(), userId)

	if err != nil {
		logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)
		logger.Info("Error find by id controller: ", err)

        return ctx.JSON(http.StatusNotFound, WebResponse{
			Message: "Failed get user id",
			Status:  0,
			Data:    nil,
		})
	}

	response := WebResponse{
		Message: "success get user by id",
		Status:  1,
		Data:    &user,
	}

    return ctx.JSON(http.StatusOK, response)
}

func (controller *UserController) FindAll(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)

	users, err := controller.UserService.FindAll(ctx.Request().Context())

	if err != nil {
		logger.Info("Error Find All users: ", err)
        return ctx.JSON(http.StatusInternalServerError, WebResponse{
            Message: "Failed GET ALL users",
            Status: 0,
            Data: nil,
        })

	}

	response := WebResponse{
		Message: "success get all users",
		Status:  1,
		Data:    users,
	}

    return ctx.JSON(http.StatusOK, response)
}
