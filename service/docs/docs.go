// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Diego Rodriguez Mancini",
            "email": "diegorodriguezmancini@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/conesearch": {
            "get": {
                "description": "Search for objects in a given region using ra, dec and radius",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conesearch"
                ],
                "summary": "Search for objects in a given region",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Right ascension in degrees",
                        "name": "ra",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Declination in degrees",
                        "name": "dec",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Radius in degrees",
                        "name": "radius",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Catalog to search in",
                        "name": "catalog",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Number of neighbors to return",
                        "name": "nneighbor",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Mastercat"
                            }
                        }
                    },
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/conesearch.ValidationError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/metadata": {
            "get": {
                "description": "Search for metadata by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "metadata"
                ],
                "summary": "Search for metadata by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID to search for",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Catalog to search in",
                        "name": "catalog",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/repository.AllwiseMetadata"
                        }
                    },
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/metadata.ValidationError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "conesearch.ValidationError": {
            "type": "object",
            "properties": {
                "errValue": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "metadata.ValidationError": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "repository.AllwiseMetadata": {
            "type": "object",
            "properties": {
                "h_m_2mass": {
                    "type": "number"
                },
                "h_msig_2mass": {
                    "type": "number"
                },
                "j_m_2mass": {
                    "type": "number"
                },
                "j_msig_2mass": {
                    "type": "number"
                },
                "k_m_2mass": {
                    "type": "number"
                },
                "k_msig_2mass": {
                    "type": "number"
                },
                "source_id": {
                    "type": "string"
                },
                "w1mpro": {
                    "type": "number"
                },
                "w1sigmpro": {
                    "type": "number"
                },
                "w2mpro": {
                    "type": "number"
                },
                "w2sigmpro": {
                    "type": "number"
                },
                "w3mpro": {
                    "type": "number"
                },
                "w3sigmpro": {
                    "type": "number"
                },
                "w4mpro": {
                    "type": "number"
                },
                "w4sigmpro": {
                    "type": "number"
                }
            }
        },
        "repository.Mastercat": {
            "type": "object",
            "properties": {
                "cat": {
                    "type": "string"
                },
                "dec": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "ipix": {
                    "type": "integer"
                },
                "ra": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "CrossWave HTTP API",
	Description:      "API for the CrossWave Xmatch service. This service allows to search for objects in a given region and to retrieve metadata from the catalogs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
