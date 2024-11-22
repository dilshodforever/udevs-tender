# Tender Management Application

This project is a backend service designed to manage tenders, enabling clients to register tenders, submit bids, and manage bidder information. It provides robust features to handle the entire tender lifecycle while ensuring high performance and reliability.

## Key Features and Implementations

### 1. Tender Registration and Management
- **Tender Creation:** Clients can create tenders with details like title, description, opening and closing dates, and relevant requirements.
- **Tender Lifecycle Management:** Includes functionalities to open, and close tenders based on their lifecycle stages, ensuring smooth operations.

### 2. Bid Management
- **Bid Submission:** Registered users can submit bids for active tenders, providing necessary details for evaluation.
- **Bid Tracking:** Users can monitor the status of their submitted bids in real time.

### 3. Performance Enhancements
- **Rate Limiting:** Prevents abuse by limiting the number of actions (such as bid submissions or tender creations) a user can perform within a specific time frame. Implemented using Redis.
- **Caching:** Frequently accessed data, such as tender lists, is cached using Redis to improve response times and reduce database load.

### 4. Security Features
- **Authentication:** Secure user authentication using modern encryption standards.
- **Data Validation:** Input validation is enforced across all endpoints to ensure data integrity and security.

### 5. Documentation
- **Swagger API Documentation:** Provides interactive and detailed API documentation for seamless integration with frontend or external systems.
- 
## Technical Stack

- **Programming Language:** Golang (chosen for its scalability and efficiency).
- **Framework:** Gin (a lightweight and performant web framework).
- **Database:** PostgreSQL (used for robust and structured data storage).
- **Caching:** Redis (used for session management, rate limiting, and caching frequently accessed data).

## Setup and Usage

### Prerequisites
- **Go:** Install Go version 1.18 or higher.
- **Docker:** Ensure Docker and Docker Compose are installed on your system.

### Installation and Setup

1. **Clone the repository:**
   ```bash
   git clone [repository URL]
   cd [repository directory]
	```

## Start the Application with Docker Compose

### To start the database:
```bash
make run_db
```

### To run the application:
```bash
make run
```

This command will:
- Build and start all required Docker containers.
- Set up the database and caching services.
- Start the backend service.

### Access the application:
The backend API will be available at `http://localhost:[PORT]`, where `PORT` is defined in your `.env` file.

---

## Development Workflow


### Stop and Clean Services:
```bash
make stop
```

---

## Key Configurations and Tools Used

- **Rate Limiting:** Redis is utilized to restrict excessive actions from clients, such as multiple bid submissions within a short period.
- **Caching:** Redis caches commonly used data like tender lists to reduce database load and speed up responses.
- **Error Logging:** A custom logger captures application errors, providing detailed insights for debugging and monitoring.
- **Session Management:** Redis handles session tokens for efficient and secure user authentication.
- **Swagger:** Comprehensive API documentation is automatically generated for easy exploration of available endpoints.
- **WebSockets:** WebSockets are used for real-time notifications to clients. When a notification is created, it is pushed to the corresponding user over an active WebSocket connection.
- **RabbitMQ:** RabbitMQ handles message queuing for asynchronous processing of notifications. It ensures that notifications are reliably delivered to users.
---

## Contribution Guidelines

We welcome contributions to improve the project! Here’s how you can contribute:

1. **Fork the repository:** Create a personal copy of the project repository on GitHub.
2. **Create a feature branch:** Develop your changes in a dedicated branch.
	```bash
	git checkout -b feature/your-feature-name
	```
3. **Commit your changes:** Ensure your commit messages are descriptive and concise.
	```bash
	git commit -m "Add feature: your-feature-name"
	```
4. **Push to your fork:** Upload your changes to your GitHub fork.
	```bash
	git push origin feature/your-feature-name
	```
5. **Open a pull request:** Submit a pull request to the main repository for review.


## Conclusion
This Tender Management Application provides a comprehensive, scalable solution for managing tenders and bids. The use of Golang ensures that the system is both efficient and highly performant, while Redis optimizes the application for speed and reliability. We hope this project will provide a strong foundation for further enhancements and real-world applications.

We appreciate your interest in the project and encourage you to contribute with any improvements or new features. Together, we can build a more powerful system to streamline tender management processes and support users in making smarter, more informed decisions.