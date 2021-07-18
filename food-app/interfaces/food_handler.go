package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/infrastructure/auth"
	"food-app/interfaces/fileupload"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Food struct {
	foodApp    application.FoodAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

func NewFood(
	fApp application.FoodAppInterface,
	uApp application.UserAppInterface,
	fd fileupload.UploadFileInterface,
	rd auth.AuthInterface,
	tk auth.TokenInterface,
) *Food {
	return &Food{
		foodApp:    fApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
	}
}

func (fo *Food) SaveFood(c *gin.Context) {
	// 首先是用户验证
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	// 从 Redis 中查询 metadata
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	// 前端需要用 map 来展示错误
	var saveFoodError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	// 初始化 food 用于验证：防止 payload 为空或是不合法的数据类型
	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	saveFoodError = emptyFood.Validate("")
	if len(saveFoodError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	// 检查用户是否存在
	_, err = fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	// 处理接受到的图片
	// file, err := c.FormFile("food_image")
	// if err != nil {
	// 	saveFoodError["invalid_file"] = "a valid file is required"
	// 	c.JSON(http.StatusUnprocessableEntity, saveFoodError)
	// 	return
	// }
	// uploadedFile, err := fo.fileUpload.UploadFile(file)
	// if err != nil {
	// 	saveFoodError["upload_err"] = err.Error()
	// 	c.JSON(http.StatusUnprocessableEntity, saveFoodError)
	// 	return
	// }

	var food = entity.Food{}
	food.UserID = userId
	food.Title = title
	food.Description = description
	// food.FoodImage = uploadedFile
	savedFood, saveErr := fo.foodApp.SaveFood(&food)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedFood)
}

func (fo *Food) UpdateFood(c *gin.Context) {
	// 首先是用户验证
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	// 从 Redis 中查询 metadata
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	// 前端需要用 map 来展示错误
	var updateFoodError = make(map[string]string)

	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}

	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	updateFoodError = emptyFood.Validate("update")
	if len(updateFoodError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateFoodError)
		return
	}
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	food, err := fo.foodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	if user.ID != food.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this food")
		return
	}
	// 因为此请求为更新数据，一个新的图片有可能会被更新。
	// 如果没有提供图片那么错误就会出现。我们需要忽视此类的错误。
	// file, _ := c.FormFile("food_image")
	// if file != nil {
	// 	food.FoodImage, err = fo.fileUpload.UploadFile(file)
	// 	//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
	// 	food.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage
	// 	if err != nil {
	// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
	// 			"upload_err": err.Error(),
	// 		})
	// 		return
	// 	}
	// }

	food.Title = title
	food.Description = description
	food.UpdatedAt = time.Now()
	updatedFood, dbUpdateErr := fo.foodApp.UpdateFood(food)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedFood)
}

func (fo *Food) GetFoodAndCreator(c *gin.Context) {
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	food, err := fo.foodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(food.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	foodAndUser := map[string]interface{}{
		"food":    food,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, foodAndUser)
}

func (fo *Food) DeleteFood(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.foodApp.DeleteFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "food deleted")
}
