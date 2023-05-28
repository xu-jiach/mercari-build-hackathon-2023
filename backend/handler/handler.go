// TODO: change password requirment if have time

package handler

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
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
	Name     string `json:"name"`
	Password string `json:"password"`
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

type getItemPasswordResponse struct {
	Password string `json:"password"`
}

type getCategoriesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type sellRequest struct {
	ItemID int32 `json:"item_id"`
}

type addItemRequest struct {
	Name         string `form:"name"`
	CategoryID   int64  `form:"category_id"`
	Price        int64  `form:"price"`
	Description  string `form:"description"`
	ItemPassword string `form:"item_password"`
}

type addItemResponse struct {
	ID int64 `json:"id"`
}

type editItemRequest struct {
	ID          int32  `form:"id"`
	Name        string `form:"name"`
	CategoryID  int64  `form:"category_id"`
	Price       int64  `form:"price"`
	Description string `form:"description"`
}

type editItemResponse struct {
	ID int64 `json:"id"`
}

type addBalanceRequest struct {
	Balance int64 `json:"balance"`
}

type getBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type loginRequest struct {
	UserID   int64  `json:"user_id"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type addCategoryRequest struct {
	Name string `json:"name"`
}

type addCategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GPT3Request struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type GPT3Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type onsitePurchaseRequest struct {
	ItemID   int32  `json:"item_id"`
	Password string `json:"password"`
}

type generateDescriptionRequest struct {
	Name       string `json:"name"`
	CategoryID int64  `json:"categoryID"`
}

type generateDescriptionResponse struct {
	Description string `json:"description"`
}

type Handler struct {
	DB                 *sql.DB
	UserRepo           db.UserRepository
	ItemRepo           db.ItemRepository
	OnsitePurchaseRepo db.OnsitePurchaseRepository
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
	// TODO: validation
	// http.StatusBadRequest(400)
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//	Validation
	// Pending to change back to the original approach if have time
	if len(req.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name is invalid")
	}
	if len(req.Password) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "password is invalid")
	}

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

	//	Validation
	// Pending to change back to the original approach if have time
	if len(req.Password) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "password is invalid")
	}

	user, err := h.UserRepo.GetUser(ctx, req.UserID)
	if err != nil {
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
	ctx := c.Request().Context()

	req := new(addItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// validation
	// if file.Size > 1<<20 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "image size must be less than 1MB")
	// }
	// if file.Size == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "image must not be empty")
	// }
	// if req.Price <= 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "price must be greater than 0")
	// }
	// validation
	// if req.Name == "" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "name must not be empty")
	// }
	// if req.Description == "" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "description must not be empty")
	// }
	// end of validation
	// end of validation

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	// // check the file type
	// // Read the first 512 bytes to determine the file type
	// buffer := make([]byte, 512)
	// _, err = src.Read(buffer)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// // Reset the reader back to the start of the file
	// _, err = src.Seek(0, io.SeekStart)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// // Detect the content type of the file
	// contentType := http.DetectContentType(buffer)
	// if !strings.HasPrefix(contentType, "image/") {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "uploaded file is not an image")
	// }

	var dest []byte
	blob := bytes.NewBuffer(dest)

	// separate the copy operation into a goroutine
	errCh := make(chan error)

	go func() {
		if _, err := io.Copy(blob, src); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	// Check if the category exists
	categoryCh := make(chan error)
	go func() {
		_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				categoryCh <- errors.New("Category does not exist")
			}
			categoryCh <- err
		}
		close(categoryCh)
	}()

	// We must ensure the copy operation has finished before we can use the blob
	if err = <-errCh; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("copy operation failed: %v", err))
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

	err = h.OnsitePurchaseRepo.AddOnsitePurchase(ctx, domain.OnsitePurchase{
		ItemID:   item.ID,
		SellerID: userID,
		Password: req.ItemPassword,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, addItemResponse{ID: int64(item.ID)})
}

func (h *Handler) EditItem(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(editItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	itemIdParam := c.Param("itemID")
	itemId, err := strconv.ParseInt(itemIdParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID")
	}

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	req.ID = int32(itemId)
	existingItem, err := h.ItemRepo.GetItem(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if existingItem.UserID != userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "User is not the owner of the item")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// validation
	// if req.Price <= 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "price must be greater than 0")
	// }
	// if req.Name == "" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "name must not be empty")
	// }
	// if req.Description == "" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "description must not be empty")
	// }
	// end of validation

	// validation
	// if file.Size > 1<<20 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "image size must be less than 1MB")
	// }
	// if file.Size == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "image must not be empty")
	// }

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	// check the file type
	// Read the first 512 bytes to determine the file type
	// buffer := make([]byte, 512)
	// _, err = src.Read(buffer)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// // Reset the reader back to the start of the file
	// _, err = src.Seek(0, io.SeekStart)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// // Detect the content type of the file
	// contentType := http.DetectContentType(buffer)
	// if !strings.HasPrefix(contentType, "image/") {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "uploaded file is not an image")
	// }

	var dest []byte
	blob := bytes.NewBuffer(dest)

	// separate the copy operation into a goroutine
	errCh := make(chan error)

	go func() {
		if _, err := io.Copy(blob, src); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	// if file.Header.Get("Content-Type") != "image/png" && file.Header.Get("Content-Type") != "image/jpeg" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "image must be png or jpeg")
	// }

	// Check if the category exists
	//
	categoryCh := make(chan error)
	go func() {
		_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				categoryCh <- errors.New("Category does not exist")
			}
			categoryCh <- err
		}
		close(categoryCh)
	}()

	if err = <-errCh; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy image file")
	}

	if err = <-categoryCh; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	item, err := h.ItemRepo.EditItem(c.Request().Context(), domain.Item{
		ID:          req.ID,
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

	return c.JSON(http.StatusOK, editItemResponse{ID: int64(item.ID)})
}
func (h *Handler) Sell(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(sellRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	UserID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	item, err := h.ItemRepo.GetItem(ctx, req.ItemID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusPreconditionFailed, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: check req.UserID and item.UserID
	// http.StatusPreconditionFailed(412)
	if item.UserID != UserID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "cannot sell other user's item")
	}
	// TODO: only update when status is initial
	// http.StatusPreconditionFailed(412)
	if item.Status != domain.ItemStatusInitial {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "invalid status. Has been sold or on sale")
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
			return echo.NewHTTPError(http.StatusNotFound, err)
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
			return echo.NewHTTPError(http.StatusNotFound, err)
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

// GetItemPassword returns item password.
// It is separeted from GetItem not to affect benchmark.
func (h *Handler) GetItemPassword(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 64)

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user")
	}

	pass, err := h.OnsitePurchaseRepo.GetItemPassword(ctx, userID, int32(itemID))
	if err != nil {
		c.Logger().Printf("failed to get item password: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, getItemPasswordResponse{Password: pass})

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
			return echo.NewHTTPError(http.StatusNotFound, "No items found for this user")
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
			return echo.NewHTTPError(http.StatusNotFound, "Categories not found")
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
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid itemID")
	}
	if itemID <= 0 || itemID > math.MaxInt32 {
		return echo.NewHTTPError(http.StatusBadRequest, "ItemID out of range")
	}

	data, err := h.ItemRepo.GetItemImage(ctx, int32(itemID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	contentType := http.DetectContentType(data[:512])
	if !strings.HasPrefix(contentType, "image/") {
		return echo.NewHTTPError(http.StatusNotFound, "uploaded file is not an image")
	}
	return c.Blob(http.StatusOK, contentType, data)
}

func (h *Handler) AddBalance(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addBalanceRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if req.Balance <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "balance must be positive")
	}

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	// TODO: not found handling
	// http.StatusPreconditionFailed(412)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
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
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
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

	// Return error if the itemID is out of range
	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 64)
	if err != nil || itemID > math.MaxInt32 || itemID < 0 {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid itemID")
	}

	// Get the item from the database.
	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusPreconditionFailed, "Item not found.")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	// Prevent the user from buying their own items.
	if item.UserID == userID {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "You cannot buy your own item.")
	}

	// If the item is not on sale, return a 412 error.
	if item.Status != domain.ItemStatusOnSale {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Item is not on sale")
	}

	// Get the user from the database.
	user, err := h.UserRepo.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusPreconditionFailed, "User not found")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Check if user has enough balance
	if user.Balance < item.Price {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Insufficient balance")
	}

	// Continue with the status update if the item is on sale and user has enough balance to finish the transactions.
	if err := h.ItemRepo.UpdateItemStatus(ctx, int32(itemID), domain.ItemStatusSoldOut); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance-item.Price); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	sellerID := item.UserID

	seller, err := h.UserRepo.GetUser(ctx, sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusPreconditionFailed, "Seller not found")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	if err := h.UserRepo.UpdateBalance(ctx, sellerID, seller.Balance+item.Price); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) OnsitePurchase(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(onsitePurchaseRequest)

	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	// Return error if the itemID is out of range
	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 64)
	if err != nil || itemID > math.MaxInt32 || itemID < 0 {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid itemID")
	}

	// Get the item from the database.
	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusPreconditionFailed, "Item not found.")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	// Prevent the user from buying their own items.
	if item.UserID == userID {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "You cannot buy your own item.")
	}

	// If the item is not on sale, return a 412 error.
	if item.Status != domain.ItemStatusOnSale {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Item is not on sale")
	}

	// Get the user from the database.
	user, err := h.UserRepo.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusPreconditionFailed, "User not found")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Check if user has enough balance
	if user.Balance < item.Price {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Insufficient balance")
	}

	isValid, err := h.OnsitePurchaseRepo.ValidatePassword(ctx, itemID, req.Password)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}
	if !isValid {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusPreconditionFailed, "Invalid password")
	}

	// Continue with the status update if the item is on sale and user has enough balance to finish the transactions.
	if err := h.ItemRepo.UpdateItemStatus(ctx, int32(itemID), domain.ItemStatusSoldOut); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance-item.Price); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, "successful")
}

// Search API
// Search Item By Keyword
func (h *Handler) SearchItemByKeyword(c echo.Context) error {
	ctx := c.Request().Context()

	// Retrieve the keyword from query parameters
	keyword := c.QueryParam("name")
	if keyword == "" {
		// Keyword is required
		c.Logger().Error("Keyword is required")
		return echo.NewHTTPError(http.StatusBadRequest, "Keyword is required")
	}

	// Call your repository method
	items, err := h.ItemRepo.GetItemByKeyword(ctx, keyword)
	if err != nil {
		c.Logger().Error(err)
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// return the response
	var res []getUserItemsResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getUserItemsResponse{ID: item.ID, Name: item.Name, Price: item.Price, CategoryName: cat.Name})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

// SearchItemAndInfoByKeyword Almost equivalent to SearchItemByKeyword.
// Returns []getItemResponse
// Kurumi created this not to disturb the bench test.
func (h *Handler) SearchItemAndInfoByKeyword(c echo.Context) error {
	ctx := c.Request().Context()

	// Retrieve the keyword from query parameters
	keyword := c.QueryParam("name")
	if keyword == "" {
		// Keyword is required
		c.Logger().Error("Keyword is required")
		return echo.NewHTTPError(http.StatusBadRequest, "Keyword is required")
	}

	// Call your repository method
	items, err := h.ItemRepo.GetItemByKeyword(ctx, keyword)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// return the response
	var res []getItemResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getItemResponse{ID: item.ID, Name: item.Name, CategoryID: item.CategoryID,
					CategoryName: cat.Name, Price: item.Price,
					Description: item.Description, Status: item.Status})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GenerateDescription(c echo.Context) error {
	fmt.Println("Received a request to generate a description...")

	ctx := c.Request().Context()
	req := new(generateDescriptionRequest)
	if err := c.Bind(req); err != nil {
		fmt.Printf("Error in Bind: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var categoryName string
	// read category name from GetCategory function
	category, err := h.ItemRepo.GetCategory(ctx, req.CategoryID) // Replace 123 with the desired category ID
	if err != nil {
		fmt.Printf("Error in GetCategory: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	fmt.Printf("Category is %v\n", category)
	categoryName = category.Name
	// Use the categoryName
	fmt.Printf("Category name is %s\n", categoryName)

	description, err := generateDescriptionWithGPT3(req.Name, categoryName)
	if err != nil {
		fmt.Printf("Error in generateDescriptionWithGPT3: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	fmt.Println("Description generated successfully, sending response...")
	return c.JSON(http.StatusOK, generateDescriptionResponse{Description: description})
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

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// AddCategory API
func (h *Handler) AddCategory(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addCategoryRequest) // Define your request struct for adding category
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Check if category already exists
	_, err := h.ItemRepo.GetCategoryByName(ctx, req.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	} else {
		// Category already exists
		return echo.NewHTTPError(http.StatusBadRequest, "Category already exists")
	}

	// If category does not exist, proceed to create it
	category, err := h.ItemRepo.AddCategory(ctx, domain.Category{
		Name: req.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, addCategoryResponse{ID: int64(category.ID)})
}

// search by category api

func (h *Handler) GetItemsByCategory(c echo.Context) error {
	ctx := c.Request().Context()

	// get category id
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// Invalid category ID
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid category ID")
	}

	// repo call
	items, err := h.ItemRepo.GetItemsByCategory(ctx, categoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

// GPT 3 API
func generateDescriptionWithGPT3(name string, categoryName string) (string, error) {

	fmt.Println("Starting the description generation process...")

	// set up the client
	client := openai.NewClient("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")
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     "text-davinci-003",
		MaxTokens: 35,
		Prompt:    fmt.Sprintf("Imagine you are writing a product description for an online store. The product is named '%s' and it belongs to category %s. Please provide a compelling and attractive description that is around 20 words long.", name, categoryName),
	}
	response, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Error in createCompletion: %s\n", err.Error())
		return "", err
	}

	if len(response.Choices) == 0 {
		err = errors.New("Received no choices from the completion request")
		fmt.Println(err)
		return "", err
	}

	// Use the text from the first choice
	fmt.Println("Description generation was successful.")
	return response.Choices[0].Text, nil
}
