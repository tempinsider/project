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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/messages/list": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/messages.ListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/messages.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/messages.Error"
                        }
                    }
                }
            }
        },
        "/services/toggle": {
            "post": {
                "description": "toggle worker status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Worker Service"
                ],
                "summary": "toggle worker status",
                "parameters": [
                    {
                        "description": "Toggle Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.ToggleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.ToggleResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/service.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "messages.Error": {
            "type": "object",
            "properties": {
                "ERROR": {}
            }
        },
        "messages.ListResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/messages.Message"
                    }
                }
            }
        },
        "messages.Message": {
            "type": "object",
            "properties": {
                "external_message_id": {
                    "type": "string"
                },
                "sent_at": {
                    "type": "string"
                }
            }
        },
        "service.Error": {
            "type": "object",
            "properties": {
                "ERROR": {}
            }
        },
        "service.ServiceStatus": {
            "type": "string",
            "enum": [
                "STOPPED",
                "WORKING"
            ],
            "x-enum-varnames": [
                "ServiceStatusStopped",
                "ServiceStatusWorking"
            ]
        },
        "service.ToggleRequest": {
            "type": "object"
        },
        "service.ToggleResponse": {
            "type": "object",
            "properties": {
                "SERVICE_STATUS": {
                    "$ref": "#/definitions/service.ServiceStatus"
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
	Title:            "Sample Messaging Worker",
	Description:      "This is a sample messaging worker.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
