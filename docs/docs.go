// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/orders": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает список всех заказов в системе",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin-orders"
                ],
                "summary": "Все заказы (только для админа)",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка получения заказов",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/orders/export": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает CSV-файл со всеми заказами (только для админа)",
                "produces": [
                    "text/csv"
                ],
                "tags": [
                    "admin-orders"
                ],
                "summary": "Экспорт заказов в CSV",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Ошибка экспорта",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/products": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Создаёт новый товар (доступно только администратору)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "admin-products"
                ],
                "summary": "Добавить новый товар",
                "parameters": [
                    {
                        "description": "Данные товара",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/chechnya-product_internal_models.Product"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Товар добавлен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Невалидный JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Нет доступа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка добавления товара",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/products/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Обновляет данные товара по ID (доступно только администратору)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "admin-products"
                ],
                "summary": "Обновить товар",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID товара",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Новые данные товара",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/chechnya-product_internal_models.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Товар обновлён",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный ID или JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка обновления товара",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Удаляет товар по ID (доступно только администратору)",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "admin-products"
                ],
                "summary": "Удалить товар",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID товара",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Товар удалён",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Нет доступа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка удаления товара",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cart": {
            "get": {
                "description": "Возвращает список товаров в корзине по user_id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Получить корзину пользователя",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный user_id",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка получения корзины",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Добавляет определённое количество товара в корзину пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Добавить товар в корзину",
                "parameters": [
                    {
                        "description": "Данные для добавления в корзину",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.AddToCartRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Добавлено в корзину",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Невалидный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка добавления в корзину",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cart/{product_id}": {
            "put": {
                "description": "Обновляет количество определённого товара в корзине текущего пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Обновить количество товара в корзине",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID товара",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Новое количество (quantity)",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "integer"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Количество обновлено",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Невалидный JSON или ошибка",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет указанный товар из корзины текущего пользователя",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Удалить товар из корзины",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID товара",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Товар удалён из корзины",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка удаления",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Возвращает JWT токен при успешном входе",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Авторизация",
                "parameters": [
                    {
                        "description": "Данные входа",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT токен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Невалидный JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неверный логин или пароль",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает ID, имя пользователя и роль из JWT",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Профиль текущего пользователя",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/order": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Создаёт заказ на основе содержимого корзины пользователя",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Оформить заказ",
                "responses": {
                    "200": {
                        "description": "Заказ успешно оформлен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка: корзина пуста или товара не хватает",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает список всех заказов текущего пользователя",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "История заказов пользователя",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка получения заказов",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/products": {
            "get": {
                "description": "Возвращает товары с фильтрацией по поиску и категории",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Получить товары",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Поиск по названию или описанию",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Категория товара",
                        "name": "category",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/chechnya-product_internal_models.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка получения товаров",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "get": {
                "description": "Возвращает один товар по его ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Получить товар по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID товара",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chechnya-product_internal_models.Product"
                        }
                    },
                    "404": {
                        "description": "Товар не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Регистрирует нового пользователя по логину и паролю",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Данные регистрации",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_handlers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь зарегистрирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Невалидный JSON или пользователь уже существует",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "chechnya-product_internal_models.Product": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "stock": {
                    "type": "integer"
                }
            }
        },
        "internal_handlers.AddToCartRequest": {
            "type": "object",
            "properties": {
                "product_id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                },
                "user_id": {
                    "description": "временно напрямую",
                    "type": "integer"
                }
            }
        },
        "internal_handlers.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "internal_handlers.RegisterRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{"http"},
	Title:            "Chechnya Product API",
	Description:      "Онлайн-магазин для продажи продуктов",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
