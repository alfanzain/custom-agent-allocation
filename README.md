# Qiscus Custom Agent Allocation Service

This repository contains the source code for a custom agent allocation service built using Go, PostgreSQL, and Redis.  This service manages the allocation of agents to customers based on availability and load.

## Main Idea

1. **Allocate Agent Webhook:**
    - Enqueue the customer request into the Redis queue.
    - Check if there are any pending requests in the queue.

2. **Mark as Resolved Webhook:**
    - Decrement the current load of the agent who resolved the customer's issue.
    - Check if there are any pending requests in the queue.

3. **Queue Check:**
    - If the queue is not empty:
        1. Call the Qiscus API (`https://omnichannel.qiscus.com/api/v1/admin/service/allocate_agent`) to get an available agent.
        2. Check if the agent exists in the database:
            - If the agent does not exist, insert a new record into the database.
            - If the agent exists, update the agent's `current_load` in the database.
        3. If the allocated agent's `data.agent.count` equals the `MAX_COUNT`, wait for a short period and repeat step 1.
        4. If the allocated agent's `data.agent.count` is less than `MAX_COUNT`, assign the agent to the customer from the head of the queue.


## Technologies Used

* **Go** 
* **PostgreSQL** Used for persistent storage of agent data.
* **Redis** Used for queuing incoming customer requests.
* **Echo Framework** Used for building the web server.
* **GORM** ORM for interacting with PostgreSQL.
* **go-redis** Client for interacting with Redis.


## Architecture

The service consists of the following components:

* **API Handlers:**  Handles incoming webhooks from Qiscus for allocating and resolving customer interactions.
* **Services:** Contains business logic for agent management, queue management, and Qiscus API interaction.
* **Models:** Defines the data structures for agents.
* **Database (PostgreSQL):** Stores agent information, including their current load and maximum capacity.
* **Queue (Redis):** Manages a queue of customer requests waiting for agent assignment.
* **Polling Mechanism (Redis):** Continuously monitors the queue for new requests (logging purpose only).


## Setup

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:** Create a `.env` file based on the `.env.example` file, populating it with your database credentials, Redis connection details, and Qiscus API keys.

4. **Run database migrations:**  (Assuming you have a migration tool set up)

5. **Start the service:**
   ```bash
   go run main.go
   ```

## Usage

The service exposes two webhook endpoints:

* `/allocate-agent/webhook`:  Triggered when a new customer interaction is initiated.
* `/mark-as-solved/webhook`: Triggered when a customer interaction is resolved.


