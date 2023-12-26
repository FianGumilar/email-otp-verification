package userService

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/FianGumilar/email-otp-verification/constans"
	"github.com/FianGumilar/email-otp-verification/models"
	"github.com/FianGumilar/email-otp-verification/services"
	. "github.com/FianGumilar/email-otp-verification/utils"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	Service services.UsecaseService
}

func NewUserService(Service services.UsecaseService) userService {
	return userService{
		Service: Service,
	}
}

func (svc userService) Register(ctx echo.Context) error {
	var result models.Response
	request := new(models.RegisterReq)

	if err := BindValidateStruct(ctx, request, "Register"); err != nil {
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusOK, result)
	}

	// Check if username is already registered
	_, exists := svc.Service.UserRepo.IsUsernameExistsByIndex(request.Username)
	if exists {
		log.Println("Error Register - IsUsernameExistsByIndex : ", "Username Already Exists")
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Username Already Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	// Check if email is already registered
	_, exists = svc.Service.UserRepo.IsEmailExistsByIndex(request.Email)
	if exists {
		log.Println("Error Register - IsEmailExistsByIndex : ", "Email Already Exists")
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Email Already Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	// Check len Password
	if len(request.Password) < 8 {
		log.Println("Error Register : ", "Password minimum length is 8 characters")
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Password minimum length is 8 characters", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)

	_ = models.User{
		FullName:  request.FullName,
		Phone:     request.Phone,
		Email:     request.Email,
		Username:  strings.ToUpper(request.Username),
		Password:  string(hashedPassword),
		CreatedAt: TimeStampNow(),
		UpdatedAt: TimeStampNow(),
	}

	// _, err := svc.Service.UserRepo.AddUser(user)
	// if err != nil {
	// 	log.Println("Error AddUser - AddUser : ", err.Error())
	// 	result = ResponseJSON(constans.FALSE_VALUE, constans.SYSTEM_ERROR_CODE, "Failed Add User", nil)
	// 	return ctx.JSON(http.StatusOK, result)
	// }

	otpCode := GenerateRandomNumber(4)
	referenceID := GenerateRandomString(16)

	log.Printf("your-otp: %s", otpCode)

	_ = svc.Service.CacheRepo.SetCache("otp:"+referenceID, []byte(otpCode))
	_ = svc.Service.CacheRepo.SetCache("user-ref:"+referenceID, []byte(request.Email))

	response := models.RegisterRes{
		ReferenceID: referenceID,
	}

	log.Println("Reponse Register ", "Success Register")
	result = ResponseJSON(constans.TRUE_VALUE, constans.SUCCESS_CODE, constans.EMPTY_VALUE, response)
	return ctx.JSON(http.StatusOK, result)
}

func (svc userService) ValidateOtp(ctx echo.Context) error {
	var result models.Response
	request := new(models.ValidateOtpReq)

	if err := BindValidateStruct(ctx, request, "ValidateOTP"); err != nil {
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusOK, result)
	}

	val, err := svc.Service.CacheRepo.GetCache("otp:" + request.ReferenceID)
	if err != nil {
		log.Println("Error ValidateOTP - GetCache : ", err.Error())
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Failed Get Cache", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	otp := string(val)

	if request.OTP != otp {
		log.Println("Error ValidateOTP - GetCache : Invalid OTP")
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Invalid OTP", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	val, err = svc.Service.CacheRepo.GetCache("user-ref:" + request.ReferenceID)
	if err != nil {
		log.Println("Error ValidateOTP - GetCache : ", err.Error())
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Invalid OTP", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	user, exists := svc.Service.UserRepo.IsEmailExistsByIndex(string(val))
	if !exists {
		log.Println("Error ValidateOTP - IsEmailExistsByIndex : Email Not Exists")
		result = ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, "Email Not Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	user.EmailVerifiedAt = time.Now()

	err = svc.Service.UserRepo.EditUser(user)
	if err != nil {
		log.Println("Error EditUser - EditUser : ", err.Error())
		result = ResponseJSON(constans.FALSE_VALUE, constans.SYSTEM_ERROR_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	log.Println("Reponse ValidateOTP", "Succes ValidateOTP")
	result = ResponseJSON(constans.TRUE_VALUE, constans.SUCCESS_CODE, constans.EMPTY_VALUE, "Success Validate OTP")
	return ctx.JSON(http.StatusOK, result)
}
