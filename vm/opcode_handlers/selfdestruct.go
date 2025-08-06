package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SelfDestructOpCode struct{}

func (*SelfDestructOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("selfdestruct op code requires the execution context to be set")
	}

	if v.StateProvider == nil {
		return fmt.Errorf("selfdestruct op code requires state provider to be set")
	}

	// SELFDESTRUCT requires one value on the stack (the beneficiary address)
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Check for static call context - SELFDESTRUCT not allowed in static calls
	frame := v.CurrentFrame()
	if frame != nil && frame.IsStatic {
		return vm.ErrStaticCallStateChange
	}

	// Pop beneficiary address from stack
	beneficiaryInt, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Convert to 20-byte address
	var beneficiary [20]byte
	beneficiaryBytes := beneficiaryInt.Bytes()
	if len(beneficiaryBytes) <= 20 {
		copy(beneficiary[20-len(beneficiaryBytes):], beneficiaryBytes)
	} else {
		copy(beneficiary[:], beneficiaryBytes[len(beneficiaryBytes)-20:])
	}

	currentAddr := v.Context.Address
	currentBalance := v.StateProvider.GetBalance(currentAddr)

	// EIP-6780: Check if the contract was created in the same transaction
	createdInTransaction := v.IsAccountCreatedInTransaction(currentAddr)

	if createdInTransaction {
		// Original SELFDESTRUCT behavior - delete the account completely

		// Transfer balance to beneficiary (even if it's the same address, which burns ether)
		if !currentBalance.IsZero() {
			if beneficiary != currentAddr {
				// Transfer to different address
				beneficiaryBalance := v.StateProvider.GetBalance(beneficiary)
				newBeneficiaryBalance := new(uint256.Int).Add(beneficiaryBalance, currentBalance)
				v.StateProvider.SetBalance(beneficiary, newBeneficiaryBalance)
			}
			// If beneficiary is same as current address, ether is burned (balance set to 0)
		}

		// Delete the account (code, storage, nonce, balance)
		err = v.StateProvider.DeleteAccount(currentAddr)
		if err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}
	} else {
		// Transfer balance to beneficiary
		if !currentBalance.IsZero() && beneficiary != currentAddr {
			// Only transfer if beneficiary is different from current address
			// If same address, no net change in balance (ether is NOT burned)
			beneficiaryBalance := v.StateProvider.GetBalance(beneficiary)
			newBeneficiaryBalance := new(uint256.Int).Add(beneficiaryBalance, currentBalance)
			v.StateProvider.SetBalance(beneficiary, newBeneficiaryBalance)

			// Set current account balance to 0
			v.StateProvider.SetBalance(currentAddr, uint256.NewInt(0))
		}
	}

	v.Stopped = true
	return nil
}
