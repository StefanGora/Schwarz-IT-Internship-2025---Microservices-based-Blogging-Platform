# 📰 Microservice Blog Application

This project was developed during my internship at **Schwarz IT**. It is a **microservice-based blog application** designed to demonstrate scalable architecture, modular service separation, and efficient communication between backend components.

The application is composed of multiple services that work together to provide authentication, blog management, and a user-friendly frontend interface.

During my internship, **my main contributions** included:
- 🐳 **Deploying the entire application** using Docker, including creating individual **Dockerfiles** for some microservice and a centralized **Docker Compose** setup.
- 🔐 **Developing and testing the Auth Service**, implementing automated integration tests and gRPC-based API communication.
- 💻 **Building and integrating the Frontend**, connecting it to the backend microservices via REST and gRPC interfaces.

## 🧩 Tech Stack

| Layer | Technology | Description |
|--------|-------------|--------------|
| **Backend** | **Go (Golang)** | Implements the business logic for both authentication and blog services |
| **Frontend** | **Vue.js** | Provides the user interface for managing and viewing blog posts |
| **Databases** | **PostgreSQL**, **MongoDB** | PostgreSQL for structured data (Auth), MongoDB for blog content and flexibility |
| **APIs** | **gRPC**, **REST** | gRPC for internal Auth API communication, REST API for Blog operations |
| **Containerization** | **Docker & Docker Compose** | Automates setup and deployment of all services locally |
| **Version Control** | **Git (Azure DevOps)** | Used for source control and project collaboration |

## 🐳 Prerequisites

Before running the project, ensure that you have the following installed on your local machine:

- [Docker](https://www.docker.com/get-started)
- [Go (Golang)](https://go.dev/doc/install) — required for building backend microservices
- [Node.js](https://nodejs.org/en) — required for the frontend
- [npm](https://www.npmjs.com/) — Node package manager

---

## 🚀 Running the Project

### After cloning the repository run the following commands


```bash
docker compose up --build
```

### This will automatically set up and run:
- Auth Service (Go + PostgreSQL)
- Blog Service (Go + PostgreSQL + MongoDB)

### Navigate to the frontend directory:

```bash
cd frontend/frontend-service
```

### Then install dependencies and start the development server:

```bash
npm install
```

```bash
npm run dev
```
## 🔧 Key Contributions  

### 1. Deployment  

My first assignment in this project was to **create a Dockerfile for the Auth Service**. I implemented a **multi-stage build** using a Golang base image for compiling the application and a minimal Scratch image for the final runtime, ensuring a lightweight and efficient container.  

Following this, I was tasked with creating the **Docker Compose configuration** for the Auth Service and its PostgreSQL database. In later iterations of the project, I expanded the compose file to include the **Blog Service** and **MongoDB dependencies**.  

To improve the reliability of the backend services, I implemented **health checks** to verify that database connections were successfully established before starting service execution. This ensured that services would not fail silently due to unavailable dependencies, improving overall system stability during deployment.

### 2. Authentication Service  

#### 2.1 CRUD Operations and Database Layer  

This service uses **PostgreSQL** to store user data. For the full implementation, you can check the `database.go` file with the following command:  

```bash
cd auth-service/internal/db/
```

In this task, I designed the table schema for the users, stored in the file **schemas.go**, which acts as a manifest file for the database.

After that, I created a Go structure:
```bash
type Database struct {
	Host     string
	User     string
	Port     int
	Password string
	Dbname   string
	ConnPool *pgxpool.Pool //pointer to a connection pool
}
```
For the database connection, I used the **pgx** Go package to create a connection pool, ensuring thread safety.

Finally, I implemented methods attached to the Database structure (leveraging Go’s struct method feature) to:

- Configure the database connection
- Set up and initialize the database
- Define full CRUD operations for the users table

#### 2.2 Integration Testing  

Another assignment on the Auth Service was to create **integration tests** once the database and API layer were fully implemented. To see the full implementation of the tests, you can check the `integration_tests.go` file using:  

```bash
cd auth-service/cmd
```

For the integration tests, I used Testcontainers to simulate an isolated Docker environment and Testify to run the test cases.

The process simplified functions with the following steps:

1. Create a test function.
2. Create a Docker network using Testcontainers.
3. Set up a PostgreSQL container and database connection.
4. Set up the Auth container.
5. Extract the PostgreSQL container IP from the network.
6. Inject the IP into the Auth container environment as the DB host.
7. Set up a gRPC client.
8. Run subtests specifying gRPC data for different test cases.


### 3. Frontend  

The frontend of the application was developed using **Vue.js with TypeScript**. My contribution included building all frontend pages, designing the project structure, and integrating with backend microservices via **Axios**.  

#### Service Structure  

- `assets/`: Static visual assets for the application.  
  - `css/`: Global stylesheets, variables, and fonts.  
  - `svg/`: SVG icons and other vector-based images.  

- `components/`: Reusable Vue components that are not full pages.  

- `pages/`: Primary view components, representing the final web pages of the application. Each file typically corresponds to a specific route defined in the router.  

- `router/`: Vue Router configuration files, defining the application's URL routes and mapping them to page components.  

- `models/` & `types/`: TypeScript definitions.  
  - `models/`: Interface definitions, storing data structures from API calls.  
  - `types/`: General TypeScript types, enums, and utility types used throughout the application.  

#### Key Features Implemented  
- Fully functional page navigation and routing.  
- API integration with Auth and Blog services using Axios.  
- Dynamic data rendering based on backend responses.
