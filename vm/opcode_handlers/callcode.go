package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type CallCodeOpCode struct{}

func (*CallCodeOpCode) Execute(v *vm.DebuggerVM) error {
	// CALLCODE requires 7 values on the stack:
	// gas, address, value, argsOffset, argsSize, retOffset, retSize
	if err := v.RequireStack(7); err != nil {
		return err
	}

	// Pop all arguments from the stack
	gas, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	address, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	value, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	argsOffset, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	argsSize, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	retOffset, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	retSize, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Extract address as bytes
	var addr [20]byte
	addressBytes := address.Bytes()
	if len(addressBytes) <= 20 {
		copy(addr[20-len(addressBytes):], addressBytes)
	}

	// For now, if no StateProvider is set, return success but do nothing
	if v.StateProvider == nil {
		// Push success result (1) onto stack
		return v.Push(uint256.NewInt(1))
	}

	// Check if the account exists
	if !v.StateProvider.AccountExists(addr) {
		// Push failure result (0) onto stack
		return v.Push(uint256.NewInt(0))
	}

	// Get the target code
	targetCode := v.StateProvider.GetCode(addr)
	if len(targetCode) == 0 {
		// Empty code means successful call with no execution
		// Clear return data area if specified
		if !retSize.IsZero() && !retOffset.IsZero() {
			retSizeInt := int(retSize.Uint64())
			retOffsetInt := int(retOffset.Uint64())
			v.Memory().Write(retOffsetInt, make([]byte, retSizeInt))
		}
		// Push success result (1) onto stack
		return v.Push(uint256.NewInt(1))
	}

	// Prepare call data
	var callData []byte
	if !argsSize.IsZero() && !argsOffset.IsZero() {
		argsSizeInt := int(argsSize.Uint64())
		argsOffsetInt := int(argsOffset.Uint64())
		callData = v.Memory().Read(argsOffsetInt, argsSizeInt)
	}

	// CALLCODE executes external code in the current context
	// This means storage writes go to the current contract
	newFrame := vm.MessageFrame{
		Code:         targetCode,
		PC:           0,
		Stack:        vm.NewStack(),
		Memory:       vm.NewMemory(),
		ReturnData:   nil,
		Gas:          gas.Uint64(),
		CallType:     vm.CallTypeCallCode,
		IsStatic:     false,
		CodeMetadata: vm.ScanCodeMetadata(targetCode),
	}

	// For CALLCODE, context keeps same address (current contract)
	// but we execute the code from the target address
	oldContext := v.Context
	newContext := &vm.ExecutionContext{
		Caller:   oldContext.Caller,  // Same caller
		Address:  oldContext.Address, // Same address (current contract)
		Origin:   oldContext.Origin,  // Same origin
		Value:    new(uint256.Int).Set(value),
		CallData: callData,
		GasPrice: oldContext.GasPrice,
		Gas:      gas.Uint64(),
		Balance:  oldContext.Balance, // Same balance (current contract)
		Block:    oldContext.Block,
	}

	// Push the new frame
	if err := v.PushFrame(newFrame); err != nil {
		return err
	}

	// Set new context
	v.Context = newContext

	// Execute the called code
	success := uint256.NewInt(1)

	// Execute the call frame
	err = v.ExecuteCall()
	if err != nil {
		// If there was an error, mark as failure
		success = uint256.NewInt(0)
		// Continue with cleanup
	}

	// Restore context and pop frame
	v.Context = oldContext
	if popErr := v.PopFrame(); popErr != nil {
		return popErr
	}

	// Handle return data if specified
	returnData := v.ReturnData()
	if !retSize.IsZero() && !retOffset.IsZero() {
		retSizeInt := int(retSize.Uint64())
		retOffsetInt := int(retOffset.Uint64())

		// Copy return data to memory, truncating if necessary
		if len(returnData) > retSizeInt {
			returnData = returnData[:retSizeInt]
		}
		v.Memory().Write(retOffsetInt, returnData)
	}

	// Push success result onto stack
	return v.Push(success)
}
