package repository

import (
	"context"
	"rest_api_pks/internal/models"

	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, image_url, name, price, description, specifications, quantity, is_favorite, in_cart FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Title, &product.ImageURL, &product.Name, &product.Price, &product.Description, &product.Specifications, &product.Quantity, &product.IsFavorite, &product.InCart)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx, "SELECT id, title, image_url, name, price, description, specifications, quantity, is_favorite, in_cart FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Title, &product.ImageURL, &product.Name, &product.Price, &product.Description, &product.Specifications, &product.Quantity, &product.IsFavorite, &product.InCart)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO products (title, image_url, name, price, description, specifications, quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, product.Title, product.ImageURL, product.Name, product.Price, product.Description, product.Specifications, product.Quantity)
	return err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := r.db.Exec(ctx, `
		UPDATE products SET title = $1, image_url = $2, name = $3, price = $4, description = $5, specifications = $6, quantity = $7
		WHERE id = $8
	`, product.Title, product.ImageURL, product.Name, product.Price, product.Description, product.Specifications, product.Quantity, product.ID)
	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) UpdateProductQuantity(ctx context.Context, id int, quantity int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET quantity = $1 WHERE id = $2", quantity, id)
	return err
}

func (r *ProductRepository) ToggleFavorite(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET is_favorite = NOT is_favorite WHERE id = $1", id)
	return err
}

func (r *ProductRepository) ToggleCart(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET in_cart = NOT in_cart WHERE id = $1", id)
	return err
}

func (r *ProductRepository) GetFilteredProducts(ctx context.Context, searchQuery string, minPrice, maxPrice float64, sortBy, sortOrder string) ([]models.Product, error) {
	query := `
		SELECT id, title, image_url, name, price, description, specifications, quantity, is_favorite, in_cart
		FROM products
		WHERE ($1 = '' OR name ILIKE $1 OR title ILIKE $1)
		AND ($2 = 0 OR price >= $2)
		AND ($3 = 0 OR price <= $3)
	`

	args := []interface{}{"%" + searchQuery + "%", minPrice, maxPrice}

	// Добавление сортировки
	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", sortBy)
		if sortOrder == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Title, &product.ImageURL, &product.Name, &product.Price, &product.Description, &product.Specifications, &product.Quantity, &product.IsFavorite, &product.InCart)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO orders (user_id, total_price, order_date, products)
        VALUES ($1, $2, $3, $4)
    `, order.UserID, order.TotalPrice, order.OrderDate, order.Products)
	return err
}

func (r *ProductRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, total_price, order_date, products
        FROM orders
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.OrderDate, &order.Products)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *ProductRepository) UpdateProductInCartStatus(ctx context.Context, productID int, inCart bool) error {
	_, err := r.db.Exec(ctx, `
        UPDATE products
        SET in_cart = $1
        WHERE id = $2
    `, inCart, productID)
	return err
}

func (r *ProductRepository) ClearCartAfterOrder(ctx context.Context, productIDs []int) error {
	// Обновляем состояние товаров в корзине
	for _, productID := range productIDs {
		_, err := r.db.Exec(ctx, `
            UPDATE products
            SET in_cart = false, quantity = 0
            WHERE id = $1
        `, productID)
		if err != nil {
			return err
		}
	}
	return nil
}
