package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type StaticCallOpCode struct{}

func (*StaticCallOpCode) Execute(v *vm.DebuggerVM) error {
	// STATICCALL requires 6 values on the stack (no value parameter):
	// gas, address, argsOffset, argsSize, retOffset, retSize
	if err := v.RequireStack(6); err != nil {
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

	// STATICCALL is a read-only call - no state changes allowed
	// It executes external code in a new context with static flag set
	newFrame := vm.MessageFrame{
		Code:         targetCode,
		PC:           0,
		Stack:        vm.NewStack(),
		Memory:       vm.NewMemory(),
		ReturnData:   nil,
		Gas:          gas.Uint64(),
		CallType:     vm.CallTypeStaticCall,
		IsStatic:     true, // Important: static calls cannot modify state
		CodeMetadata: vm.ScanCodeMetadata(targetCode),
	}

	// For STATICCALL, create new context with no value transfer (value = 0)
	oldContext := v.Context
	newContext := &vm.ExecutionContext{
		Caller:   oldContext.Address, // Current contract is the caller
		Address:  addr,               // Target address
		Origin:   oldContext.Origin,  // Origin remains the same
		Value:    uint256.NewInt(0),  // No value transfer in static call
		CallData: callData,
		GasPrice: oldContext.GasPrice,
		Gas:      gas.Uint64(),
		Balance:  v.StateProvider.GetBalance(addr),
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

	// Execute the call frame (static call - no state changes allowed)
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
