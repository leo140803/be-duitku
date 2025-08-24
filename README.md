#  Duitku Finance App Backend

A robust Go backend application for managing personal finances, built with modern architecture and Supabase integration.

## ✨ Features

-  **Authentication & Authorization** - JWT-based auth with refresh tokens
- 💰 **Account Management** - Create and manage multiple financial accounts
- 📊 **Category Management** - Organize transactions with custom categories
- 💳 **Transaction Tracking** - Comprehensive income/expense management
- 🗄️ **Supabase Integration** - Modern database with real-time capabilities
- 🚀 **RESTful API** - Clean, documented endpoints
- 🔒 **Middleware Protection** - Secure route handling

## 🛠️ Tech Stack

- **Language**: Go 1.24.6+
- **Framework**: Gin (Web framework)
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT tokens
- **Environment**: Go modules

## 📋 Prerequisites

- Go 1.24.6 or higher
- Supabase project
- Git

##  Quick Start

### 1. Clone Repository
```bash
git clone <your-repo-url>
cd finance-app-backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Environment Setup
Create a `.env` file in the root directory:

```env
# Supabase Configuration
SUPABASE_PROJECT_ID=your-project-id
SUPABASE_ANON_KEY=your-anon-key-here

# Server Configuration
PORT=8080
```

### 4. Get Supabase Credentials
1. Go to [Supabase Dashboard](https://app.supabase.com)
2. Select your project
3. Navigate to **Settings** → **API**
4. Copy the **Project ID** and **anon/public key**
5. Add them to your `.env` file

**⚠️ Important**: Use only the project ID (e.g., `tiqqkmlbntrbpczlocmt`) for `SUPABASE_PROJECT_ID`, not the full URL.

### 5. Run Application
```bash
# Development mode
go run cmd/main.go

# Build and run
go build -o main cmd/main.go
./main
```

The server will start on the port specified in your `PORT` environment variable (default: 8080).

## ️ Database Schema

Run these SQL commands in your Supabase SQL editor:

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

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

-- Create indexes for better performance
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_categories_user_id ON categories(user_id);
```

##  API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/login` | Login user |
| `POST` | `/api/auth/logout` | Logout user |

#### Request/Response Examples

**Register User**
```bash
POST /api/auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

**Login User**
```bash
POST /api/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

**Auth Response**
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

**📝 Note**: The `balance_after` field is automatically calculated by the API based on the previous transaction balance.

## 📦 Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gin-gonic/gin` | Latest | Web framework |
| `github.com/lengzuo/supa` | Latest | Supabase Go client |
| `github.com/joho/godotenv` | Latest | Environment variable loading |

## 🔄 Migration from Direct PostgreSQL

This application has been migrated from using direct PostgreSQL connections to Supabase. The main changes include:

1. ✅ Replaced `pgx` driver with `github.com/lengzuo/supa`
2. ✅ Updated database operations to use Supabase client
3. ✅ Changed environment variables from `DATABASE_URL` to `SUPABASE_PROJECT_ID` and `SUPABASE_ANON_KEY`
4. ✅ Simplified database queries using Supabase's query builder

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](../../issues) page
2. Create a new issue with detailed information
3. Contact the development team

---

**Made with ❤️ by the Duitku Team**
