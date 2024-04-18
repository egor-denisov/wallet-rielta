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
        "/wallet": {
            "post": {
                "description": "Создает новый кошелек с уникальным ID. Идентификатор генерируется сервером.\n\nСозданный кошелек должен иметь сумму 100.0 у.е. на балансе",
                "tags": [
                    "Wallet"
                ],
                "summary": "Создание кошелька",
                "responses": {
                    "200": {
                        "description": "Кошелек создан",
                        "schema": {
                            "$ref": "#/definitions/entity.Wallet"
                        }
                    },
                    "500": {
                        "description": "Не удалось создать кошелек"
                    },
                    "504": {
                        "description": "Время ожидания вышло"
                    }
                }
            }
        },
        "/wallet/{walletId}": {
            "get": {
                "tags": [
                    "Wallet"
                ],
                "summary": "Получение текущего состояния кошелька",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID кошелька",
                        "name": "walletId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Wallet"
                        }
                    },
                    "404": {
                        "description": "Указанный кошелек не найден"
                    },
                    "500": {
                        "description": "Не удалось выполнить запрос"
                    },
                    "504": {
                        "description": "Время ожидания вышло"
                    }
                }
            }
        },
        "/wallet/{walletId}/history": {
            "get": {
                "description": "Возвращает историю транзакций по указанному кошельку.",
                "tags": [
                    "Wallet"
                ],
                "summary": "Получение историй входящих и исходящих транзакций",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID кошелька",
                        "name": "walletId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "История транзакций получена",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Transaction"
                            }
                        }
                    },
                    "404": {
                        "description": "Указанный кошелек не найден"
                    },
                    "500": {
                        "description": "Не удалось выполнить запрос"
                    },
                    "504": {
                        "description": "Время ожидания вышло"
                    }
                }
            }
        },
        "/wallet/{walletId}/send": {
            "post": {
                "tags": [
                    "Wallet"
                ],
                "summary": "Перевод средств с одного кошелька на другой",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID кошелька",
                        "name": "walletId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Запрос перевода средств",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.transactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Перевод успешно проведен"
                    },
                    "400": {
                        "description": "Ошибка в пользовательском запросе"
                    },
                    "404": {
                        "description": "Исходящий кошелек не найден"
                    },
                    "500": {
                        "description": "Ошибка перевода"
                    },
                    "504": {
                        "description": "Время ожидания вышло"
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Transaction": {
            "description": "Денежный перевод.",
            "type": "object",
            "required": [
                "amount",
                "from",
                "time",
                "to"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 30
                },
                "from": {
                    "type": "string",
                    "example": "5b53700ed469fa6a09ea72bb78f36fd9"
                },
                "time": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-02-04T17:25:35.448Z"
                },
                "to": {
                    "type": "string",
                    "example": "eb376add88bf8e70f80787266a0801d5"
                }
            }
        },
        "entity.Wallet": {
            "description": "Состояние кошелька.",
            "type": "object",
            "required": [
                "balance",
                "id"
            ],
            "properties": {
                "balance": {
                    "type": "integer",
                    "example": 100
                },
                "id": {
                    "type": "string",
                    "example": "5b53700ed469fa6a09ea72bb78f36fd9"
                }
            }
        },
        "v1.transactionRequest": {
            "description": "Запрос перевода средств.",
            "type": "object",
            "required": [
                "amount",
                "to"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 100
                },
                "to": {
                    "type": "string",
                    "example": "eb376add88bf8e70f80787266a0801d5"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Wallet",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
