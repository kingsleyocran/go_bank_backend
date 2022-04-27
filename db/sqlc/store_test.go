package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	//Create a new store
	store := NewStore(testDB)

	//To create efficient unit test we need to create new accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	log.Println(">> before:", account1.Balance, account2.Balance)

	/*
		From my experience, writing database transaction is something we must always be very careful with.
		It can be easy to write, but can also easily become a nightmare if we don’t handle the concurrency carefully.
		So the best way to make sure that our transaction works well is to run it with several concurrent go routines.
		Let’s say I want to run n = 5 concurrent transfer transactions, and each of them will transfer an amount of 10 from account 1 to account 2.
		So I will use a simple for loop with n iterations:
	*/

	n := 5
	amount := int64(10)

	// we use make(chan datatype) to define a channel to pass our result from concurrency
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {

		//Printing the tx name is good for debugging errors like deadlock
		txName := fmt.Sprintf("tx %d", i+1)

		//You can use context.WithValue() to get thr tx name
		ctx := context.WithValue(context.Background(), txKey, txName)

		//the keyword go starts a go routine
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			//Now, inside the go routine, we can send err to the errs channel using this arrow operator <-.
			//The channel should be on the left, and data to be sent should be on the right of the arrow operator.
			//Therefore we store any [err] to the [errs] channel and any [result] to the [results] channel
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {

		//To receive the error from the channel, we use the same arrow operator,
		//but this time, the channel is on the right of the arrow, and the variable to store the received data is on the left.
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//Now the tx is done wen can check if the transfer exists by getting an error. We don't neccesarily need to check if the values match
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries from the fromEntry side
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//Check entries from the toEntry side
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		//the amount and diff can't be equal but instead because we are increasing by amount,
		//diff of any instance say diff3 divided by amount should not leave a reminder
		// 1 * amount, 2 * amount, 3 * amount, ..., n * amount
		require.True(t, diff1%amount == 0)

		//Wherever we divide diff/amount, the value at any instance should be in range [1,5]
		//since we are running the concurrently for 5 times
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		//We also need to make sure the current instance say diff4 doesnt have k=1,2 or 3
		//so we use existed map to create a map[int]bool and assign at each point the current k = true.
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	log.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	//Now after n transactions, the balance of account 1 must decrease by n * amount.
	//So we require the updatedAccount1.Balance to equal to that value. amount is of type int64,
	//so we need to convert n to int64 before doing the multiplication.
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

//This test is to check for two concurrent transaction involving the same pair of accounts(account table form simple_bank_backend).
//This usually calls for a potential deadlock
//To avoid this make sure your statements always acquire locks in the same order even in reverse.

func TestTransferTxDeadlock(t *testing.T) {
	//Create a new store
	store := NewStore(testDB)

	//To create efficient unit test we need to create new accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	log.Println(">> before:", account1.Balance, account2.Balance)

	/*
		Here, let’s say we’re gonna run n = 10 concurrent transactions.
		The idea is to have 5 transactions that send money from account 1 to account 2, and another 5
		transactions that send money in reverse direction, from account 2 to account 1.
	*/

	n := 10
	amount := int64(10)

	// we use make(chan datatype) to define a channel to pass our result from concurrency
	errs := make(chan error)
	//We don't need the results since we have checked that in TestTransferTx
	//results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		//Printing the tx name is good for debugging errors like deadlock
		txName := fmt.Sprintf("tx %d", i+1)

		//You can use context.WithValue() to get thr tx name
		ctx := context.WithValue(context.Background(), txKey, txName)

		fromAccountID := account1.ID
		toAccountID := account2.ID

		//At odd number and even number instances accounts change
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check for only erros no results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	log.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	//Now after n transactions, the balance of account 1 must be equal to that of account2
	//since the other 5 additional tx reverses the changes in balance
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
