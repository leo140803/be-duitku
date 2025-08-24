# Finance App Backend

A Go backend application for managing personal finances, now integrated with Supabase.

## Features

- Account management (create, list accounts)
- Supabase integration for database operations
- RESTful API endpoints

## Setup

### Prerequisites

- Go 1.24.6 or higher
- Supabase project

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Supabase Configuration
SUPABASE_PROJECT_ID=your-project-id
SUPABASE_ANON_KEY=your-anon-key-here

# Server Configuration
PORT=8080
```

### Getting Supabase Credentials

1. Go to your Supabase project dashboard
2. Navigate to Settings > API
3. Copy the Project ID (not the full URL) and anon/public key
4. Add them to your `.env` file

**Note:** Use only the project ID (e.g., `tiqqkmlbntrbpczlocmt`) for `SUPABASE_PROJECT_ID`, not the full URL.

### Database Schema

Make sure you have the following tables in your Supabase database:

```sql
-- Accounts table
CREATE TABLE accounts (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    initial_balance DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Categories table
CREATE TABLE categories (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Transactions table
CREATE TABLE transactions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id TEXT NOT NULL,
    account_id UUID REFERENCES accounts(id),
    category_id UUID REFERENCES categories(id),
    date DATE NOT NULL,
    description TEXT,
    amount DECIMAL(10,2) NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    balance_after DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/main.go
```

The server will start on the port specified in your `PORT` environment variable (default: 8080).

## API Endpoints

### Authentication (Public)

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user

#### Register Request Body

```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

#### Login Request Body

```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

#### Auth Response

```json
{
    "user": {
        "id": "user-uuid",
        "email": "user@example.com",
        "created_at": "2024-01-15T10:30:00Z"
    },
    "access_token": "jwt-token-here",
    "refresh_token": "refresh-token-here"
}
```

### Protected Endpoints

All endpoints below require authentication. Include the access token in the Authorization header:
```
Authorization: Bearer your-access-token-here
```

#### User Profile

- `GET /api/auth/profile` - Get current user profile

### Accounts

- `GET /api/accounts` - Get all accounts for current user
- `POST /api/accounts` - Create a new account

#### Create Account Request Body

```json
{
    "name": "Main Account",
    "initial_balance": 1000.00
}
```

### Categories

- `GET /api/categories` - Get all categories for current user
- `POST /api/categories` - Create a new category

#### Create Category Request Body

```json
{
    "name": "Food & Dining"
}
```

### Transactions

- `GET /api/transactions` - Get all transactions for current user (ordered by date ascending)
- `POST /api/transactions` - Create a new transaction

#### Create Transaction Request Body

```json
{
    "account_id": "account-uuid",
    "category_id": "category-uuid",
    "date": "2024-01-15",
    "description": "Lunch at restaurant",
    "amount": 25.50,
    "type": "EXPENSE"
}
```

**Note:** The `balance_after` field is automatically calculated by the API based on the previous transaction balance.

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/lengzuo/supa` - Supabase Go client
- `github.com/joho/godotenv` - Environment variable loading

## Migration from Direct PostgreSQL

This application has been migrated from using direct PostgreSQL connections to Supabase. The main changes include:

1. Replaced `pgx` driver with `github.com/lengzuo/supa`
2. Updated database operations to use Supabase client
3. Changed environment variables from `DATABASE_URL` to `SUPABASE_PROJECT_ID` and `SUPABASE_ANON_KEY`
4. Simplified database queries using Supabase's query builder
