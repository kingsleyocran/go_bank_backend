package db

import (
	"context"
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
		//the keyword go starts a go routine
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
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
	//existed := make(map[int]bool)

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
	}

}
