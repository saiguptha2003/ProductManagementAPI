# Product Management API

##### This is a Product Management API built using the Gin web framework in Go. The API allows users to create, fetch, and manage product details. It integrates with a relational database and message queues to handle product-related operations and asynchronous image processing.
---
## Features
##### Create Product: Add new products to the database with details such as name, description, images, and price.
##### Get Product by ID: Retrieve detailed information about a specific product.
##### Fetch All Products: Query and filter products based on user ID, price range, and product name.
##### Redis Cache Integration: Optimize product fetching by caching frequently accessed data.
---
## Technologies Used
##### Go: Programming language used for implementation.
##### Gin: Web framework for building RESTful APIs.
##### sqlite: Database used to store product information.
##### Redis: Cache layer for improved read performance.
##### RabbitMQ (Queue): Handles asynchronous processing of image-related tasks.
---
## Installation Setup

### Clone the Repository
```bash
git clone https://github.com/saiguptha2003/ProductManagementAPI
cd ProductManagementAPI
```

### Install Dependencies Ensure you have Go installed
```bash
go mod tidy
```
### Setup Sqlite3 Database
```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_images TEXT,
    compressed_product_images TEXT,
    product_price DECIMAL(10, 2) NOT NULL
);
```
### RabbitMQ Installation 
```bash
# latest RabbitMQ 4.0.x
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4.0-management
```
##### port address : 15672
##### Host address : localhost
### Redis Installation and Setup
#### Install Redis on Windows Use Redis on Windows for development
#### Redis is not officially supported on Windows. However, you can install Redis on Windows for development by following the instructions below.

#### To install Redis on Windows, you'll first need to enable WSL2 (Windows Subsystem for Linux). WSL2 lets you run Linux binaries natively on Windows. For this method to work, you'll need to be running Windows 10 version 2004 and higher or Windows 11.

### Install or enable WSL2
#### Microsoft provides detailed instructions for installing WSL. Follow these instructions, and take note of the default Linux distribution it installs. This guide assumes Ubuntu.

```bash
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

sudo apt-get update
sudo apt-get install redis
```
#### Lastly, start the Redis server like so:
```bash
sudo service redis-server start
```
##### port address : 6379
##### Host address : localhost

---
## Run Application
```bash
go run main.go
```
##### port address :8080
##### hOST address : localhost
---
## API TESTING

#### Open PostMan
##### API ENDPOINTS
###### POST /products
###### request Body
```json
{
  "user_id": 1,
  "product_name": "bottle",
  "product_description": "red color",
  "product_images": ["https://unsplash.com/photos/a-white-bottle-with-a-black-cap-OUjR8lrGccs", "https://unsplash.com/photos/a-person-holding-a-bottle-with-a-string-attached-to-it--qAVQodEMpA"],
  "product_price": 99.99
}
```
##### Response
```json 
{
  "message": "Product created successfully. Images are being processed asynchronously."
}
```
##### GET /products/:id

#####  GET /products

###### Query Parameters

###### user_id (required): The ID of the user.
###### min_price (optional): Minimum product price.
###### max_price (optional): Maximum product price.
###### product_name (optional): Filter by product name (partial match).

---
## Special Feactures

#### Image Storage and Processing
##### Amazon S3 Integration:
Product images are securely stored in an Amazon S3 bucket, ensuring scalability, reliability, and accessibility. This approach minimizes server storage dependency and enhances the system's performance.
##### Asynchronous Mechanism:
###### The API uses asynchronous processing to handle image uploads. This ensures that the user experience remains smooth, as image processing tasks are offloaded to a queue system (RabbitMQ). The product creation response is immediate, while the image upload and processing continue in the background.
---

## OUTPUTS:
![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/addProducts.png)

![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/awsS3Storage.png)
![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/createProductLog.png)
![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/getProductByID.png)
![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/parameterGetProduct.png)
![Alt text](https://github.com/saiguptha2003/ProductManagementAPI/blob/main/outputs/rabbitMQWorking.png)

---
## Cache Strategies Used
#### Lazy Loading (Cache-Aside) used for get Product with ID
##### On a request, the application first checks the Redis cache.
##### If the data is found (cache hit), it is returned directly from the cache.
##### If not found (cache miss), the data is fetched from the SQLite database, stored in the Redis cache for future requests, and then returned.

#### Selective Caching used for get Product with parameters
##### Product information based on user preferences (user_id) and filters like product_name or price_range was selectively cached to ensure high query performance for repeated requests.

---
## Compromised Features

#### Testing : problem with module. tried to resolve problems but i have endterm exam tomorrow.
#### Postgres SQL : Due to low system specifications and no free space i used sqlite3 database for development

