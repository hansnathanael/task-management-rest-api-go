Task Management REST API (Go)

A RESTful API built using Go (net/http) to manage tasks with a clean and
structured backend architecture.

Features: 
- CRUD operations 
- RESTful API design 
- JSON handling 
- API Key Authentication 
- Layered architecture 
- Unit testing

Tech Stack: 
- Go (Golang) 
- net/http 
- MySQL 
- OpenAPI 3.0 
- httprouter

How to Run: 
1. Clone repo 
2. Setup database 
3. Run: go run main.go

API Endpoints: 
GET /tasks 
GET /tasks/{id} 
POST /tasks 
PUT /tasks/{id}
DELETE /tasks/{id}

Sample Request: 
{ 
	“title”: Self Learning - Go”, 
	“description”: “Learn how to make Go program with RESTful API based”,
	“status”: “Done” 
}

Sample Response: 
{ 
	“code”: 200, 
	“status”: “OK”, 
	“data”: { 
			“id”: 1,
			“title”: “Self Learning - Go”, 
			“description”: “Learn how to make Go program with RESTful API based”, 
			“status”: “Done” 
		} 
}

Author: Hans Nathanael
