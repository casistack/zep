// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/api/v1/sessions/{sessionId}/memory": {
            "get": {
                "description": "get memory by session id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memory"
                ],
                "summary": "Returns a memory (latest summary and list of messages) for a given session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "session_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Last N messages. Overrides memory_window configuration",
                        "name": "lastn",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Memory"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "add memory messages by session id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memory"
                ],
                "summary": "Add memory messages to a given session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "session_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Memory messages",
                        "name": "memoryMessages",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Memory"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete memory messages by session id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memory"
                ],
                "summary": "Delete memory messages for a given session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "session_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    }
                }
            }
        },
        "/api/v1/sessions/{sessionId}/search": {
            "post": {
                "description": "search memory messages by session id and query",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "search"
                ],
                "summary": "Search memory messages for a given session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "session_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit the number of results returned",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "description": "Search query",
                        "name": "searchPayload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SearchPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SearchResult"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Memory": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Message"
                    }
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": true
                },
                "summary": {
                    "$ref": "#/definitions/models.Summary"
                }
            }
        },
        "models.Message": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": true
                },
                "role": {
                    "type": "string"
                },
                "token_count": {
                    "type": "integer"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "models.SearchPayload": {
            "type": "object",
            "properties": {
                "meta": {
                    "description": "reserved for future use",
                    "type": "object",
                    "additionalProperties": true
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.SearchResult": {
            "type": "object",
            "properties": {
                "dist": {
                    "type": "number"
                },
                "message": {
                    "$ref": "#/definitions/models.Message"
                },
                "meta": {
                    "description": "reserved for future use",
                    "type": "object",
                    "additionalProperties": true
                },
                "summary": {
                    "description": "reserved for future use",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.Summary"
                        }
                    ]
                }
            }
        },
        "models.Summary": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": true
                },
                "recent_message_uuid": {
                    "description": "The most recent message UUID that was used to generate this summary",
                    "type": "string"
                },
                "token_count": {
                    "type": "integer"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "server.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
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
	BasePath:         "/apt/v1",
	Schemes:          []string{"http", "https"},
	Title:            "zep Long-term Memory API",
	Description:      "zep stores, manages, enriches, and searches long-term memory for conversational AI applications",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
