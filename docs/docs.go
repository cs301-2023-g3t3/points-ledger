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
        "/accounts": {
            "get": {
                "description": "Retrieves a list of point accounts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "points"
                ],
                "summary": "Get all Point Accounts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PointsAccount"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/accounts/paginate": {
            "get": {
                "description": "Retrieves a list of point accounts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "points"
                ],
                "summary": "Get all Point Accounts by Pagination",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "size",
                        "name": "size",
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
                                "$ref": "#/definitions/models.PointsAccount"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/accounts/{id}": {
            "get": {
                "description": "Retrieve a list of Points Account By UserID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "points"
                ],
                "summary": "Get Points Account by UserId",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UserID",
                        "name": "UserID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PointsAccount"
                            }
                        }
                    },
                    "400": {
                        "description": "Id cannot be empty",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Points Account not found with Id",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update Points By Id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "points"
                ],
                "summary": "Update Points by Id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID",
                        "name": "ID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Points adjusted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid JSON body",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "404": {
                        "description": "User not found with Id",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check the health of the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Get Health",
                "responses": {
                    "200": {
                        "description": "Sucess"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "models.PointsAccount": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
