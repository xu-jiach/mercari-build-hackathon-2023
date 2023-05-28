package db

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xu-jiach/mecari-build-hackathon-2023/backend/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	UpdateBalance(ctx context.Context, id int64, balance int64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

func (r *UserDBRepository) AddUser(ctx context.Context, user domain.User) (int64, error) {
	// Start a new transaction
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password); err != nil {
		tx.Rollback()
		return 0, echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}

	// Retrieve the ID of the last inserted row
	row := r.QueryRowContext(ctx, "SELECT id FROM users WHERE rowid = LAST_INSERT_ROWID()")

	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Password, &user.Balance)
}

func (r *UserDBRepository) UpdateBalance(ctx context.Context, id int64, balance int64) error {
	if _, err := r.ExecContext(ctx, "UPDATE users SET balance = ? WHERE id = ?", balance, id); err != nil {
		return err
	}
	return nil
}

type ItemRepository interface {
	AddItem(ctx context.Context, item domain.Item) (domain.Item, error)
	EditItem(ctx context.Context, item domain.Item) (domain.Item, error)
	AddCategory(ctx context.Context, category domain.Category) (domain.Category, error)
	GetItem(ctx context.Context, id int32) (domain.Item, error)
	GetItemImage(ctx context.Context, id int32) ([]byte, error)
	GetOnSaleItems(ctx context.Context) ([]domain.Item, error)
	GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error)
	GetCategory(ctx context.Context, id int64) (domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetItemByKeyword(ctx context.Context, keyword string) ([]domain.Item, error)
	UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error
	GetItemsByCategory(ctx context.Context, categoryID int64) ([]domain.Item, error) // for category search page
}

type ItemDBRepository struct {
	*sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &ItemDBRepository{DB: db}
}

func (r *ItemDBRepository) AddItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return domain.Item{}, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO items (name, price, description, category_id, seller_id, image, status) VALUES (?, ?, ?, ?, ?, ?, ?)", item.Name, item.Price, item.Description, item.CategoryID, item.UserID, item.Image, item.Status); err != nil {
		tx.Rollback()
		return domain.Item{}, echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}

	row := r.QueryRowContext(ctx, "SELECT * FROM items WHERE rowid = LAST_INSERT_ROWID()")

	var res domain.Item
	return res, row.Scan(&res.ID, &res.Name, &res.Price, &res.Description, &res.CategoryID, &res.UserID, &res.Image, &res.Status, &res.CreatedAt, &res.UpdatedAt)
}

// Create an Edit Method
func (r *ItemDBRepository) EditItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	// start a new transaction
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return domain.Item{}, err
	}

	if _, err := r.ExecContext(ctx, "UPDATE items SET name = ?, price = ?, description = ?, category_id = ?, image = ?, status = ? WHERE id = ?", item.Name, item.Price, item.Description, item.CategoryID, item.Image, item.Status, item.ID); err != nil {
		tx.Rollback()
		return domain.Item{}, echo.NewHTTPError(http.StatusConflict, err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Item{}, err
	}

	return item, nil
}

func (r *ItemDBRepository) GetItem(ctx context.Context, id int32) (domain.Item, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM items WHERE id = ?", id)

	var item domain.Item
	return item, row.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt)
}

func (r *ItemDBRepository) GetItemImage(ctx context.Context, id int32) ([]byte, error) {
	row := r.QueryRowContext(ctx, "SELECT image FROM items WHERE id = ?", id)

	var image []byte
	if err := row.Scan(&image); err != nil {
		return nil, err
	}

	return image, nil
}

func (r *ItemDBRepository) GetOnSaleItems(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE status = ? ORDER BY updated_at desc", domain.ItemStatusOnSale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE seller_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error {
	if _, err := r.ExecContext(ctx, "UPDATE items SET status = ? WHERE id = ?", status, id); err != nil {
		return err
	}
	return nil
}

func (r *ItemDBRepository) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM category WHERE id = ?", id)

	var cat domain.Category
	return cat, row.Scan(&cat.ID, &cat.Name)
}

func (r *ItemDBRepository) AddCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	var newCategory domain.Category

	err := r.QueryRowContext(ctx, "INSERT INTO category (name) VALUES (?) RETURNING *", category.Name).Scan(&newCategory.ID, &newCategory.Name)
	if err != nil {
		return domain.Category{}, err
	}

	return newCategory, nil
}

func (r *ItemDBRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *ItemDBRepository) GetItemByKeyword(ctx context.Context, keyword string) ([]domain.Item, error) {
	pattern := "%" + keyword + "%"
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE name LIKE ?", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) GetCategoryByName(ctx context.Context, name string) (domain.Category, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM category WHERE name = ?", name)

	var cat domain.Category
	return cat, row.Scan(&cat.ID, &cat.Name)
}

// categories id page method
func (r *ItemDBRepository) GetItemsByCategory(ctx context.Context, categoryID int64) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE category_id = ?", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type OnsitePurchaseRepository interface {
	AddOnsitePurchase(ctx context.Context, purchase domain.OnsitePurchase) error
	ValidatePassword(ctx context.Context, itemID int64, password string) (bool, error)
}

type OnsitePurchaseDBRepository struct {
	*sql.DB
}

func NewOnsitePurchaseRepository(db *sql.DB) OnsitePurchaseRepository {
	return &OnsitePurchaseDBRepository{DB: db}
}

func (r *OnsitePurchaseDBRepository) AddOnsitePurchase(ctx context.Context, purchase domain.OnsitePurchase) error {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO onsite_purchase (item_id, seller_id, password) VALUES (?, ?, ?)",
		purchase.ItemID, purchase.SellerID, purchase.Password); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}

	return nil
}

func (r *OnsitePurchaseDBRepository) ValidatePassword(ctx context.Context, itemID int64, password string) (bool, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	row := r.QueryRowContext(ctx, "SELECT * FROM onsite_purchase WHERE item_id = ?", itemID)

	var purchase domain.OnsitePurchase
	if err := row.Scan(&purchase.ID, &purchase.ItemID, &purchase.SellerID, &purchase.Password); err != nil {
		return false, err
	}

	if purchase.Password == string(hash) {
		return true, nil
	} else {
		return false, nil
	}
}
