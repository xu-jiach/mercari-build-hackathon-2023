package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	// newly added
	"github.com/go-playground/validator/v10"

	// new added one ended here
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/xu-jiach/mecari-build-hackathon-2023/backend/db"
	"github.com/xu-jiach/mecari-build-hackathon-2023/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	logFile = getEnv("LOGFILE", "access.log")
)

type JwtCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type InitializeResponse struct {
	Message string `json:"message"`
}

type registerRequest struct {
	// name and password are set be aplhanumeric only
	Name     string `json:"name" validate:"required,min=4,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=20,password"`
}

type registerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type getUserItemsResponse struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	CategoryName string `json:"category_name"`
}

type getOnSaleItemsResponse struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	CategoryName string `json:"category_name"`
}

type getItemResponse struct {
	ID           int32             `json:"id"`
	Name         string            `json:"name"`
	CategoryID   int64             `json:"category_id"`
	CategoryName string            `json:"category_name"`
	UserID       int64             `json:"user_id"`
	Price        int64             `json:"price"`
	Description  string            `json:"description"`
	Status       domain.ItemStatus `json:"status"`
}

type getCategoriesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type sellRequest struct {
	UserID int64 `json:"user_id"`
	ItemID int32 `json:"item_id"`
}

type addItemRequest struct {
	Name         string `form:"name" validate:"required,min=1,max=35"`
	CategoryID   int64  `form:"category_id" validate:"number"`
	CategoryName string `form:"category_name" validate:"min=0,max=35"`
	Price        int64  `form:"price" validate:"required,number,gte=0,lte=99999999"`
	Description  string `form:"description" validate:"required,min=1,max=2555"`
}

type editItemRequest struct {
	ItemID       int32  `form:"item_id" `
	Name         string `form:"name" validate:"required,min=1,max=35"`
	CategoryID   int64  `form:"category_id" validate:"number"`
	CategoryName string `form:"category_name" validate:"min=0,max=35"`
	Price        int64  `form:"price" validate:"required,number,gte=0,lte=99999999"`
	Description  string `form:"description" validate:"required,min=1,max=2555"`
}

type addItemResponse struct {
	ID int64 `json:"id"`
}

type editItemResponse struct {
	ID int64 `json:"id"`
}

type addBalanceRequest struct {
	Balance int64 `json:"balance" validate:"number,gte=0"`
}

type getBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type loginRequest struct {
	UserID   int64  `json:"user_id" validate:"required,number"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Handler struct {
	DB       *sql.DB
	UserRepo db.UserRepository
	ItemRepo db.ItemRepository
}

func GetSecret() string {
	if secret := os.Getenv("SECRET"); secret != "" {
		return secret
	}
	return "secret-key"
}

func (h *Handler) Initialize(c echo.Context) error {
	err := os.Truncate(logFile, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to truncate access log"))
	}

	err = db.Initialize(c.Request().Context(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to initialize"))
	}

	return c.JSON(http.StatusOK, InitializeResponse{Message: "Success"})
}

func (h *Handler) AccessLog(c echo.Context) error {
	return c.File(logFile)
}

func (h *Handler) Register(c echo.Context) error {
	// TODO: validation -- done
	// http.StatusBadRequest(400)
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validation
	validate := validator.New()

	// validate the password
	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate the name
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		// Initialize an empty slice for error messages
		ErrMsgs := make([]string, len(validationErrors))

		// Loop through the validation errors, mapping each to a user-friendly message
		for i, e := range validationErrors {
			ErrMsgs[i] = mapErrorMessage(e)
		}

		// Return an HTTP error messages
		return echo.NewHTTPError(http.StatusBadRequest, ErrMsgs)
	}

	// end of validation

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	userID, err := h.UserRepo.AddUser(c.Request().Context(), domain.User{Name: req.Name, Password: string(hash)})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, registerResponse{ID: userID, Name: req.Name})
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	// TODO: validation
	// http.StatusBadRequest(400)
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validation -- want to change the validation instead of calling the external function

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		// Initialize an empty slice for error messages
		ErrMsgs := make([]string, len(validationErrors))

		// Loop through the validation errors, mapping each to a user-friendly message
		for i, e := range validationErrors {
			ErrMsgs[i] = mapErrorMessage(e)
		}

		// Return an HTTP error messages
		return echo.NewHTTPError(http.StatusBadRequest, ErrMsgs)
	}

	// end of validation

	user, err := h.UserRepo.GetUser(ctx, req.UserID)
	if err != nil {
		// add another error msg when a id not found
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		req.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(GetSecret()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, loginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: encodedToken,
	})
}

func (h *Handler) AddItem(c echo.Context) error {
	// TODO: validation
	// http.StatusBadRequest(400)
	ctx := c.Request().Context()

	req := new(addItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validation -- want to change the validation instead of calling the external function

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		// Initialize an empty slice for error messages
		ErrMsgs := make([]string, len(validationErrors))

		// Loop through the validation errors, mapping each to a user-friendly message
		for i, e := range validationErrors {
			ErrMsgs[i] = mapErrorMessage(e)
		}

		// Return an HTTP error messages
		return echo.NewHTTPError(http.StatusBadRequest, ErrMsgs)
	}

	// end of validation

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	var dest []byte
	blob := bytes.NewBuffer(dest)
	// TODO: pass very big file
	// http.StatusBadRequest(400)

	// passing error if the image is bigger than 1MB
	const MaxSize = 1 << 20 // 1MB
	if file.Size > MaxSize {
		return echo.NewHTTPError(http.StatusBadRequest, "file size exceeds limit")
	}

	// passing error if the image is not jpeg or png
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		return echo.NewHTTPError(http.StatusBadRequest, "file type must be jpeg or png")
	}

	if _, err := io.Copy(blob, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// check if the category exists
	_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create a new category with the provided name if doesn't
			var category domain.Category
			category, err = h.ItemRepo.AddCategory(ctx, domain.Category{
				Name: req.CategoryName,
			})
			if err != nil {
				// Handle error creating category
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			// Update req.CategoryID with the new category's ID
			req.CategoryID = category.ID
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	item, err := h.ItemRepo.AddItem(c.Request().Context(), domain.Item{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		UserID:      userID,
		Price:       req.Price,
		Description: req.Description,
		Image:       blob.Bytes(),
		Status:      domain.ItemStatusInitial,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, addItemResponse{ID: int64(item.ID)})
}

// EditItem edits an item
func (h *Handler) EditItem(c echo.Context) error {
	ctx := c.Request().Context()

	itemIdParam := c.Param("itemID")
	log.Println("itemIdParam:", itemIdParam)
	itemId, err := strconv.ParseInt(itemIdParam, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID")
	}

	req := new(editItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	req.ItemID = int32(itemId)

	// Validation
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		// Initialize an empty slice for error messages
		ErrMsgs := make([]string, len(validationErrors))

		// Loop through the validation errors, mapping each to a user-friendly message
		for i, e := range validationErrors {
			ErrMsgs[i] = mapErrorMessage(e)
		}

		// Return an HTTP error messages
		return echo.NewHTTPError(http.StatusBadRequest, ErrMsgs)
	}

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	// Check if the user is the owner of the item
	existingItem, err := h.ItemRepo.GetItem(ctx, req.ItemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if existingItem.UserID != userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "User is not the owner of the item")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	var dest []byte
	blob := bytes.NewBuffer(dest)

	// Similar checks as in AddItem
	const MaxSize = 1 << 20 // 1MB
	if file.Size > MaxSize {
		return echo.NewHTTPError(http.StatusBadRequest, "file size exceeds limit")
	}

	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		return echo.NewHTTPError(http.StatusBadRequest, "file type must be jpeg or png")
	}

	if _, err := io.Copy(blob, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Check if the category exists
	_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Assuming req has an ID field
	item, err := h.ItemRepo.EditItem(c.Request().Context(), domain.Item{
		ID:          req.ItemID,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		UserID:      userID,
		Price:       req.Price,
		Description: req.Description,
		Image:       blob.Bytes(),
		Status:      domain.ItemStatusInitial,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, editItemResponse{ID: int64(item.ID)}) // Assume there is a similar structure as addItemResponse
}

func (h *Handler) Sell(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(sellRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// get the current logged in user information to req
	var err error
	req.UserID, err = getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	item, err := h.ItemRepo.GetItem(ctx, req.ItemID)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// no need to check user 404 because it is loggged in already

	// TODO: check req.UserID and item.UserID
	// http.StatusPreconditionFailed(412)
	if req.UserID != item.UserID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "the item is not yours")
	}

	// TODO: only update when status is initial
	// http.StatusPreconditionFailed(412)
	if item.Status != domain.ItemStatusInitial {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "the item is in the state of initial. It has been on sale or sold.")
	}
	if err := h.ItemRepo.UpdateItemStatus(ctx, item.ID, domain.ItemStatusOnSale); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) GetOnSaleItems(c echo.Context) error {
	ctx := c.Request().Context()

	items, err := h.ItemRepo.GetOnSaleItems(ctx)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []getOnSaleItemsResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getOnSaleItemsResponse{ID: item.ID, Name: item.Name, Price: item.Price, CategoryName: cat.Name})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetItem(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	category, err := h.ItemRepo.GetCategory(ctx, item.CategoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, getItemResponse{
		ID:           item.ID,
		Name:         item.Name,
		CategoryID:   item.CategoryID,
		CategoryName: category.Name,
		UserID:       item.UserID,
		Price:        item.Price,
		Description:  item.Description,
		Status:       item.Status,
	})
}

func (h *Handler) GetUserItems(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid userID type")
	}

	items, err := h.ItemRepo.GetItemsByUserID(ctx, userID)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "No items listed for this user")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []getUserItemsResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getUserItemsResponse{ID: item.ID, Name: item.Name, Price: item.Price, CategoryName: cat.Name})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

	cats, err := h.ItemRepo.GetCategories(ctx)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]getCategoriesResponse, len(cats))
	for i, cat := range cats {
		res[i] = getCategoriesResponse{ID: cat.ID, Name: cat.Name}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetImage(c echo.Context) error {
	ctx := c.Request().Context()

	// TODO: overflow
	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 64)
	if err != nil || itemID > math.MaxInt32 || itemID < math.MinInt32 {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid or out of range itemID")
	}

	// オーバーフローしていると。ここのint32(itemID)がバグって正常に処理ができないはず
	data, err := h.ItemRepo.GetItemImage(ctx, int32(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// decode content type from
	contentType := http.DetectContentType(data)

	if contentType != "image/jpeg" && contentType != "image/png" {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid image type")
	}

	return c.Blob(http.StatusOK, contentType, data)
	// TODO: might need to change it to accept both jpeg and png
}

func (h *Handler) AddBalance(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addBalanceRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate the request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Checking if the balance to be added is negative
	if req.Balance < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot add negative balance")
	}

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance+req.Balance); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) GetBalance(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, getBalanceResponse{Balance: user.Balance})
}

func (h *Handler) Purchase(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	// TODO: overflow
	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 64)
	if err != nil || itemID > math.MaxInt32 || itemID < math.MinInt32 {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid or out of range itemID")
	}

	// TODO: update only when item status is on sale
	// http.StatusPreconditionFailed(412)
	// move this part upward for early check before the status update
	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "Item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: not to buy own items. 自身の商品を買おうとしていたら、http.StatusPreconditionFailed(412)
	// Prevent the user from buying their own items.
	if item.UserID == userID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Cannot purchase own item")
	}

	// If the item is not on sale, return a 412 error.
	if item.Status != domain.ItemStatusOnSale {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Item is not on sale")
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: balance consistency
	// Check if user has enough balance
	// move it before it change the status
	if user.Balance < item.Price {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Not enough balance")
	}

	// Continue with the status update if the item is on sale and user has enough balance to finished the transactions.
	if err := h.ItemRepo.UpdateItemStatus(ctx, int32(itemID), domain.ItemStatusSoldOut); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// オーバーフローしていると。ここのint32(itemID)がバグって正常に処理ができないはず
	if err := h.ItemRepo.UpdateItemStatus(ctx, int32(itemID), domain.ItemStatusSoldOut); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance-item.Price); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sellerID := item.UserID

	seller, err := h.UserRepo.GetUser(ctx, sellerID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "Seller not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := h.UserRepo.UpdateBalance(ctx, sellerID, seller.Balance+item.Price); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "successful")
}

func getUserID(c echo.Context) (int64, error) {
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return -1, fmt.Errorf("invalid token")
	}
	claims := user.Claims.(*JwtCustomClaims)
	if claims == nil {
		return -1, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

// Search Item By Keyword
func (h *Handler) SearchItemByKeyword(c echo.Context) error {
	ctx := c.Request().Context()

	// Retrieve the keyword from query parameters
	keyword := c.QueryParam("name")
	if keyword == "" {
		// Keyword is required
		return echo.NewHTTPError(http.StatusBadRequest, "Keyword is required")
	}

	// Call your repository method
	items, err := h.ItemRepo.GetItemByKeyword(ctx, keyword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return the response
	var res []getUserItemsResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getUserItemsResponse{ID: item.ID, Name: item.Name, Price: item.Price, CategoryName: cat.Name})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Error returns when registering and logining
func mapErrorMessage(e validator.FieldError) string {
	var ErrMsg string

	switch e.Tag() {
	case "required":
		ErrMsg = fmt.Sprintf("%s is required", e.Field())
	case "min":
		ErrMsg = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		ErrMsg = fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
	case "alphanum":
		ErrMsg = fmt.Sprintf("%s must only contain alphanumeric characters", e.Field())
	case "gte":
		ErrMsg = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		ErrMsg = fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "password":
		ErrMsg = ("The password needs to be 6-20 characters long, and contain at least two groups of the following: uppercase letters, lowercase letters, numbers, and symbols")
	case "number":
		ErrMsg = fmt.Sprintf("%s must be a number", e.Field())
	default:
		ErrMsg = fmt.Sprintf("%s is not valid", e.Field())
	}

	return ErrMsg
}

// Customize validator function for password
func passwordValidator(fl validator.FieldLevel) bool {
	num := `[0-9]`
	az := `[a-z]`
	AZ := `[A-Z]`
	special := `[!@#\$%\^&\*\(\)\\_\+\-=\[\]\{\};':",.<>\/\?\\|]`
	pwd := fl.Field().String()

	// check if the length of the password meet the requirment
	if len(pwd) < 6 || len(pwd) > 32 {
		return false
	}

	// set a counter
	// the password needs to contains at least 1 characters from the two groups among the four
	count := 0
	if m, _ := regexp.MatchString(num, pwd); m {
		count++
	}
	if m, _ := regexp.MatchString(az, pwd); m {
		count++
	}
	if m, _ := regexp.MatchString(AZ, pwd); m {
		count++
	}
	if m, _ := regexp.MatchString(special, pwd); m {
		count++
	}

	if count < 2 {
		return false
	}

	return true

}
