# ecommerce

This project is an e-commerce platform built using microservices architecture. It includes services for managing users, products, purchases, carts, and inventory. The services communicate asynchronously using Kafka.

## Project Structure


## Prerequisites

Make sure you have Docker and Docker Compose installed on your system.

## Setup and Run

Follow these steps to start the e-commerce platform:

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/ecommerce.git
   cd ecommerce
docker-compose up --build


###User Service
Create User:

Method: POST
URL: http://localhost:8001/users
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}

##Get Users:

Method: GET
URL: http://localhost:8001/users

###Create Product:

Method: POST
URL: http://localhost:8002/products
{
  "name": "Product A",
  "description": "Description of Product A",
  "price": 24.99
}

###Get Products:
URL: http://localhost:8002/products

##Create Purchase

Method: POST
URL: http://localhost:8003/purchases
{
  "user_id": 1,
  "product_id": 1,
  "quantity": 2,
  "total_price": 49.99
}

###Get Purchases
URL: http://localhost:8003/purchases

##Cart Service

Method: POST
URL: http://localhost:8004/cart

{
  "user_id": 1,
  "product_id": 1,
  "quantity": 3
}
###Get Cart for User
Method: GET
URL: http://localhost:8004/cart/{user_id}



This `README.md` file provides an overview of the project, instructions for setting up and running the services, details on available endpoints, and information on contributing to the project. Make sure to replace `yourusername` with your actual GitHub username if you plan to publish this repository.
