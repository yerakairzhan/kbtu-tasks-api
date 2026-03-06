# User Service -- Pagination, Filtering and Friends API

Backend service written in **Go** that demonstrates:

-   SQL pagination
-   dynamic filtering
-   sorting
-   many-to-many user friendships
-   common friends queries
-   REST API handlers

The service works with a **Users database** and exposes HTTP endpoints
to query and filter users.

------------------------------------------------------------------------

# Features

## Pagination

Users can be retrieved using **limit & offset pagination**.

SQL uses:

    LIMIT
    OFFSET
    ORDER BY

Example logic:

``` go
offset := (page - 1) * pageSize
```

Response format:

``` json
{
  "data": [],
  "totalCount": 100,
  "page": 1,
  "pageSize": 10
}
```

------------------------------------------------------------------------

# User Model

``` go
type User struct {
    ID        UUID      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Gender    string    `json:"gender"`
    BirthDate time.Time `json:"birthDate"`
}
```

------------------------------------------------------------------------

# Database Schema

## Users table

Contains at least **20 user records**.

Required fields:

-   id
-   name
-   email
-   gender
-   birth_date

------------------------------------------------------------------------

## Friendships table

Many-to-many relationship between users.

``` sql
CREATE TABLE user_friends (
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, friend_id)
);
```

Constraint:

-   a user **cannot be friends with himself**

------------------------------------------------------------------------

# Core Functions

## GetPaginatedUsers

Retrieves users with:

-   pagination
-   dynamic filtering
-   sorting

Supports filtering by:

-   ID
-   Name
-   Email
-   Gender
-   Birth date

Sorting examples:

    order_by=id
    order_by=name
    order_by=email

Default sorting is applied if `order_by` is not provided.

------------------------------------------------------------------------

## GetCommonFriends

Returns common friends between two users.

Requirements:

-   single SQL query
-   avoid **N+1 problem**
-   use **JOIN**

Concept:

    User A friends
    INTERSECT
    User B friends

------------------------------------------------------------------------

# REST API

The service exposes at least **two HTTP endpoints**.

## Get users

    GET /users

Query parameters:

    page
    page_size
    order_by
    name
    email
    gender
    birth_date

Example:

    GET /users?page=1&page_size=10&order_by=name

------------------------------------------------------------------------

## Get common friends

    GET /users/common-friends

Parameters:

    user1_id
    user2_id

Example:

    GET /users/common-friends?user1_id=1&user2_id=5

------------------------------------------------------------------------

# Demo Requirements

The backend service must:

-   run on port **8080**
-   return JSON responses
-   correctly process filtering, pagination, and sorting
-   avoid N+1 database queries

Demo in Postman should show:

1.  pagination flow with limit & offset
2.  filtering by ID, Name, Email
3.  filtering by three random fields
4.  pagination + filtering + order_by
5.  common friends query

Database must contain:

-   **20 users**
-   friendship relationships
-   at least **two users sharing three common friends**

------------------------------------------------------------------------

# Optional Features

## Easy -- "The Baby Gopher"

-   pagination
-   filtering
-   ordering
-   REST endpoints

------------------------------------------------------------------------

## Medium -- "The Gopher-at-Work"

### Soft deletes

Users can be marked as deleted without removing them from the database.

------------------------------------------------------------------------

## Advanced -- "The Gopher Wizard"

### Cursor Pagination

Pagination using cursors instead of offset.

### Query Validation

Safe query builder.

Allowed operators:

    =
    <
    >
    ILIKE

Reject invalid filters before SQL execution.

------------------------------------------------------------------------

### Transactional Friendship Consistency

Friendship must be **bidirectional**.

    A → B
    B → A

Use:

-   SQL transactions
-   or database triggers

If one insert fails → rollback.

------------------------------------------------------------------------

### Friend Recommendation System

    GetFriendRecommendations(userID)

Logic:

    friends of friends
    exclude existing friends
    exclude the user himself

Sort by number of mutual friends.

Constraint:

-   avoid **N+1 queries**

------------------------------------------------------------------------

# Running the Project

    go run main.go

Service runs on:

    http://localhost:8080

------------------------------------------------------------------------

# Deliverables

  Category                           Points
  ---------------------------------- --------------
  Project flow                       0.5
  Pagination + filtering + sorting   0.75
  Common friends function            0.75
  HTTP handlers                      0.5
  Demo video                         0.5
  **Total**                          **3 points**

------------------------------------------------------------------------

# Deadline

    15 March 2026
    23:59
