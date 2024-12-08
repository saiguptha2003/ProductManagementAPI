package handlers

import (
	"ProductManagement/cache"
	"ProductManagement/db"
	"ProductManagement/queue"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)
func CreateProduct(c *gin.Context) {
	type ProductInput struct {
		UserID            int      `json:"user_id" binding:"required"`
		ProductName       string   `json:"product_name" binding:"required"`
		ProductDescription string   `json:"product_description"`
		ProductImages     []string `json:"product_images"`
		ProductPrice      float64  `json:"product_price" binding:"required"`
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productImages := strings.Join(input.ProductImages, ",")

	res, err := db.DB.Exec(
		`INSERT INTO products (user_id, product_name, product_description, product_images, product_price) 
         VALUES (?, ?, ?, ?, ?)`,
		input.UserID, input.ProductName, input.ProductDescription, productImages, input.ProductPrice,
	)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	productID, err := res.LastInsertId()
	if err != nil {

		log.Printf("Error fetching last inserted ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	message := map[string]interface{}{
		"product_id":   productID,
		"image_urls":   input.ProductImages,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process images"})
		return
	}

	if err := queue.Publish(messageBytes); err != nil {
		log.Printf("Error publishing message to queue: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process images"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully. Images are being processed asynchronously."})
}


// func GetProductByID(c *gin.Context) {
// 	id := c.Param("id")
// 	var (
// 		productName        string
// 		productDescription string
// 		productImages      string
// 		productPrice       float64
// 		compressedImages   sql.NullString
// 	)
// err := db.DB.QueryRow(
// 		`SELECT product_name, product_description, product_images, product_price, compressed_product_images 
//          FROM products WHERE id = ?`, id,
// 	).Scan(&productName, &productDescription, &productImages, &productPrice, &compressedImages)
// 	if err != nil {
// 		log.Printf("Error fetching product: %v", err)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"product_name":        productName,
// 		"product_description": productDescription,
// 		"product_images":      strings.Split(productImages, ","),
// 		"product_price":       productPrice,
// 		"compressed_images":   compressedImages.String,
// 	})
// }
func GetProductByID(c *gin.Context) {
    id := c.Param("id")
    cacheKey := "product_" + id

    // Check Redis cache
    cachedProduct, err := cache.Get(cacheKey)
    if err == nil {
        var product map[string]interface{}
        json.Unmarshal([]byte(cachedProduct), &product)
        c.JSON(http.StatusOK, product)
        return
    }

    var (
        productName        string
        productDescription string
        productImages      string
        productPrice       float64
        compressedImages   sql.NullString
    )
    err = db.DB.QueryRow(
        `SELECT product_name, product_description, product_images, product_price, compressed_product_images 
         FROM products WHERE id = ?`, id,
    ).Scan(&productName, &productDescription, &productImages, &productPrice, &compressedImages)
    if err != nil {
        log.Printf("Error fetching product: %v", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    product := map[string]interface{}{
        "product_name":        productName,
        "product_description": productDescription,
        "product_images":      strings.Split(productImages, ","),
        "product_price":       productPrice,
        "compressed_images":   compressedImages.String,
    }

    // Cache the result
    productBytes, _ := json.Marshal(product)
    cache.Set(cacheKey, string(productBytes), time.Hour)

    c.JSON(http.StatusOK, product)
}

// func GetProductsHandler(c *gin.Context) {
// 	query := c.Request.URL.Query()

// 	userID := query.Get("user_id")
// 	minPrice := query.Get("min_price")
// 	maxPrice := query.Get("max_price")
// 	productName := query.Get("product_name")

// 	if userID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
// 		return
// 	}

// 	sqlQuery := `SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price 
//                  FROM products WHERE user_id = ?`
// 	params := []interface{}{userID}

// 	if minPrice != "" {
// 		sqlQuery += " AND product_price >= ?"
// 		price, err := strconv.ParseFloat(minPrice, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price value"})
// 			return
// 		}
// 		params = append(params, price)
// 	}

// 	if maxPrice != "" {
// 		sqlQuery += " AND product_price <= ?"
// 		price, err := strconv.ParseFloat(maxPrice, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price value"})
// 			return
// 		}
// 		params = append(params, price)
// 	}

// 	if productName != "" {
// 		sqlQuery += " AND product_name LIKE ?"
// 		params = append(params, "%"+productName+"%")
// 	}

// 	rows, err := db.DB.Query(sqlQuery, params...)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
// 		return
// 	}
// 	defer rows.Close()

// 	var products []map[string]interface{}
// 	for rows.Next() {
// 		var id int
// 		var userID int
// 		var productName string
// 		var productDescription string
// 		var productImages string
// 		var compressedImages sql.NullString
// 		var productPrice float64

// 		err := rows.Scan(&id, &userID, &productName, &productDescription, &productImages, &compressedImages, &productPrice)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
// 			return
// 		}

// 		product := map[string]interface{}{
// 			"id":                 id,
// 			"user_id":            userID,
// 			"product_name":       productName,
// 			"product_description": productDescription,
// 			"product_images":     strings.Split(productImages, ","),
// 			"compressed_images":  compressedImages.String,
// 			"product_price":      productPrice,
// 		}
// 		products = append(products, product)
// 	}

// 	c.JSON(http.StatusOK, products)
// }
func GetProductsHandler(c *gin.Context) {
    query := c.Request.URL.Query()

    userID := query.Get("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
        return
    }

    cacheKey := "products_user_" + userID
    cachedProducts, err := cache.Get(cacheKey)
    if err == nil {
        var products []map[string]interface{}
        json.Unmarshal([]byte(cachedProducts), &products)
        c.JSON(http.StatusOK, products)
        return
    }

    sqlQuery := `SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price 
                 FROM products WHERE user_id = ?`
    params := []interface{}{userID}

    minPrice := query.Get("min_price")
    if minPrice != "" {
        sqlQuery += " AND product_price >= ?"
        price, err := strconv.ParseFloat(minPrice, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price value"})
            return
        }
        params = append(params, price)
    }

    maxPrice := query.Get("max_price")
    if maxPrice != "" {
        sqlQuery += " AND product_price <= ?"
        price, err := strconv.ParseFloat(maxPrice, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price value"})
            return
        }
        params = append(params, price)
    }

    productName := query.Get("product_name")
    if productName != "" {
        sqlQuery += " AND product_name LIKE ?"
        params = append(params, "%"+productName+"%")
    }

    rows, err := db.DB.Query(sqlQuery, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
        return
    }
    defer rows.Close()

    var products []map[string]interface{}
    for rows.Next() {
        var id int
        var userID int
        var productName string
        var productDescription string
        var productImages string
        var compressedImages sql.NullString
        var productPrice float64

        err := rows.Scan(&id, &userID, &productName, &productDescription, &productImages, &compressedImages, &productPrice)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
            return
        }

        product := map[string]interface{}{
            "id":                 id,
            "user_id":            userID,
            "product_name":       productName,
            "product_description": productDescription,
            "product_images":     strings.Split(productImages, ","),
            "compressed_images":  compressedImages.String,
            "product_price":      productPrice,
        }
        products = append(products, product)
    }

    // Cache the result
    productsBytes, _ := json.Marshal(products)
    cache.Set(cacheKey, string(productsBytes), time.Hour)

    c.JSON(http.StatusOK, products)
}
