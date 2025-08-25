package handlers

import (
    "context"
    "net/http"
	"log"

    "github.com/gin-gonic/gin"
    "github.com/leo140803/finance-app-backend/config"
    "github.com/leo140803/finance-app-backend/models"
    "github.com/lengzuo/supa/utils/enum"
)

func GetTransactions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var transactions []models.Transaction

	// bisa ditambah filter pakai query param
	err := config.SupaClient.DB.From("transactions").Select("*").
		Order("date", enum.OrderAsc).
		Eq("user_id", userID.(string)).Execute(context.Background(), &transactions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func CreateTransaction(c *gin.Context) {
	// Get user ID dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var tx models.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Inject user id
	tx.UserID = userID.(string)
	tx.ID = ""
	tx.CreatedAt = ""

	// 1️⃣ Ambil saldo terkini dari tabel accounts
	var accounts []models.Account
	err := config.SupaClient.DB.From("accounts").
		Select("*").
		Eq("id", tx.AccountID).
		Execute(context.Background(), &accounts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account: " + err.Error()})
		return
	}
	if len(accounts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	lastBalance := accounts[0].InitialBalance

	// 2️⃣ Hitung saldo baru & isi balance_after di transaksi
	newBalance := 0.0
	if tx.Type == "INCOME" {
		newBalance = lastBalance + tx.Amount
	} else {
		newBalance = lastBalance - tx.Amount
	}
	tx.BalanceAfter = newBalance
	log.Printf("💰 CreateTransaction: account=%s lastBalance=%.2f amount=%.2f type=%s newBalance=%.2f",
	tx.AccountID, lastBalance, tx.Amount, tx.Type, newBalance)

	// 3️⃣ Insert transaksi baru
	var result interface{}
	err = config.SupaClient.DB.From("transactions").
		Insert(tx).
		Execute(context.Background(), &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	// 4️⃣ Update saldo di tabel accounts
	updateData := models.UpdateAccountBalance{InitialBalance: newBalance}
	err = config.SupaClient.DB.
		From("accounts").
		Update(updateData).
		Eq("id", tx.AccountID).
		Execute(context.Background(), nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance: " + err.Error()})
		return
	}

	// 5️⃣ Response
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Transaction created successfully",
		"account_id":    tx.AccountID,
		"amount":        tx.Amount,
		"type":          tx.Type,
		"user_id":       tx.UserID,
		"balance_after": tx.BalanceAfter,
	})
}

func UpdateTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	// Ambil transaksi lama
	var oldTx []models.Transaction
	err := config.SupaClient.DB.From("transactions").
		Select("*").
		Eq("id", txID).
		Eq("user_id", userID.(string)).
		Execute(context.Background(), &oldTx)

	if err != nil || len(oldTx) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	prevTx := oldTx[0]

	// Bind data baru
	var newTx models.Transaction
	if err := c.ShouldBindJSON(&newTx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	newTx.ID = txID
	newTx.UserID = userID.(string)

	ctx := context.Background()

	// -----------------------------
	// Step 1: Revert saldo account lama
	// -----------------------------
	var oldAcc []models.Account
	err = config.SupaClient.DB.From("accounts").
		Select("*").
		Eq("id", prevTx.AccountID).
		Eq("user_id", userID.(string)).
		Execute(ctx, &oldAcc)
	if err != nil || len(oldAcc) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Old account not found"})
		return
	}
	accountOld := oldAcc[0]

	if prevTx.Type == "INCOME" {
		accountOld.InitialBalance -= prevTx.Amount
	} else {
		accountOld.InitialBalance += prevTx.Amount
	}

	// Update saldo account lama
	err = config.SupaClient.DB.From("accounts").
		Update(models.UpdateAccountBalance{InitialBalance: accountOld.InitialBalance}).
		Eq("id", prevTx.AccountID).
		Execute(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old account balance"})
		return
	}

	// -----------------------------
	// Step 2: Apply transaksi baru ke account baru
	// -----------------------------
	var newAcc []models.Account
	err = config.SupaClient.DB.From("accounts").
		Select("*").
		Eq("id", newTx.AccountID).
		Eq("user_id", userID.(string)).
		Execute(ctx, &newAcc)
	if err != nil || len(newAcc) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "New account not found"})
		return
	}
	accountNew := newAcc[0]

	if newTx.Type == "INCOME" {
		accountNew.InitialBalance += newTx.Amount
	} else {
		accountNew.InitialBalance -= newTx.Amount
	}
	newTx.BalanceAfter = accountNew.InitialBalance

	// Update transaksi
	var result interface{}
	err = config.SupaClient.DB.From("transactions").
		Update(newTx).
		Eq("id", txID).
		Execute(ctx, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update saldo account baru
	err = config.SupaClient.DB.From("accounts").
		Update(models.UpdateAccountBalance{InitialBalance: accountNew.InitialBalance}).
		Eq("id", newTx.AccountID).
		Execute(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update new account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}


func DeleteTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	// Ambil transaksi lama
	var oldTx []models.Transaction
	err := config.SupaClient.DB.From("transactions").
		Select("*").
		Eq("id", txID).
		Eq("user_id", userID.(string)).
		Execute(context.Background(), &oldTx)

	if err != nil || len(oldTx) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	tx := oldTx[0]

	// Ambil account
	var accounts []models.Account
	err = config.SupaClient.DB.From("accounts").
		Select("*").
		Eq("id", tx.AccountID).
		Execute(context.Background(), &accounts)

	if err != nil || len(accounts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	account := accounts[0]
	lastBalance := account.InitialBalance

	// Reverse pengaruh transaksi
	if tx.Type == "INCOME" {
		lastBalance -= tx.Amount
	} else {
		lastBalance += tx.Amount
	}

	// Hapus transaksi
	var result interface{}
	err = config.SupaClient.DB.From("transactions").
		Delete().
		Eq("id", txID).
		Execute(context.Background(), &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	// Update saldo account
	updateData := models.UpdateAccountBalance{InitialBalance: lastBalance}
	err = config.SupaClient.DB.From("accounts").
		Update(updateData).
		Eq("id", tx.AccountID).
		Execute(context.Background(), nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}




