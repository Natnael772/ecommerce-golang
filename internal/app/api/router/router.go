package router

import (
	"ecommerce-app/internal/domain/address"
	"ecommerce-app/internal/domain/auth"
	"ecommerce-app/internal/domain/cart"
	"ecommerce-app/internal/domain/cartitem"
	"ecommerce-app/internal/domain/category"
	"ecommerce-app/internal/domain/coupon"
	"ecommerce-app/internal/domain/inventory"
	"ecommerce-app/internal/domain/order"
	"ecommerce-app/internal/domain/payment"
	"ecommerce-app/internal/domain/product"
	"ecommerce-app/internal/domain/review"
	"ecommerce-app/internal/domain/user"
	"ecommerce-app/internal/infra/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)


func NewRouter(pool *pgxpool.Pool) chi.Router {
	q:= db.NewQueries(pool)
	r := chi.NewRouter()

	// Enable CORS middleware if needed
    r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
    AllowCredentials: false,
	}))


	// Product domain setup
	productRepo := product.NewRepository(q)
	productSvc := product.NewService(productRepo)
	productRoutes := product.Routes(productSvc)

	// User domain setup
	userRepo := user.NewRepository(q)
	userSvc := user.NewService(userRepo)
	userRoutes := user.Routes(userSvc)

	// Category domain setup
	categoryRepo := category.NewRepository(q)
	categorySvc := category.NewService(categoryRepo)
	categoryRoutes := category.Routes(categorySvc)

	// Review domain setup
	reviewRepo := review.NewRepository(q)
	reviewSvc := review.NewService(reviewRepo)
	reviewRoutes := review.Routes(reviewSvc)

	// Coupon domain setup
	couponRepo := coupon.NewRepository(q)
	couponSvc := coupon.NewService(couponRepo)
	couponRoutes := coupon.Routes(couponSvc)

	// Cart domain setup
	cartRepo := cart.NewRepository(q)
	cartSvc := cart.NewService(cartRepo)
	cartRoutes := cart.Routes(cartSvc)

	// CartItem domain setup
	cartItemRepo := cartitem.NewRepository(q)
	cartItemSvc := cartitem.NewService(cartItemRepo, cartSvc)
	cartItemRoutes := cartitem.Routes(cartItemSvc)

	// Address domain setup
	addressRepo := address.NewRepository(q)
	addressSvc := address.NewService(addressRepo)
	addressRoutes := address.Routes(addressSvc)

	// Order domain setup
	orderRepo := order.NewRepository(q)
	orderSvc := order.NewService(orderRepo, productSvc)
	orderRoutes := order.Routes(orderSvc)

	// Payment domain setup
	paymentRepo := payment.NewPaymentRepository(q)
	paymentSvc := payment.NewPaymentService(paymentRepo, orderSvc)
	paymentRoutes := payment.Routes(paymentSvc)

	// Auth domain setup
	authRepo := auth.NewRepository(q)
	authSvc := auth.NewService(authRepo)
	authRoutes := auth.Routes(authSvc)

	// Inventory domain setup
	inventoryRepo := inventory.NewRepository(q)
	inventorySvc := inventory.NewService(inventoryRepo)
	inventoryRoutes := inventory.Routes(inventorySvc)

	// Mount domain routes
	r.Mount("/users", userRoutes)
	r.Mount("/products", productRoutes)
	r.Mount("/categories", categoryRoutes)
	r.Mount("/reviews", reviewRoutes)
	r.Mount("/coupons", couponRoutes)
	r.Mount("/carts", cartRoutes)
	r.Mount("/cart-items", cartItemRoutes)
	r.Mount("/addresses", addressRoutes)
	r.Mount("/orders", orderRoutes)
	r.Mount("/payments", paymentRoutes)
	r.Mount("/auth", authRoutes)
	r.Mount("/inventories", inventoryRoutes)

	return r
}
