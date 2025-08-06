package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
)

type ExtCodeCopyOpCode struct{}

func (*ExtCodeCopyOpCode) Execute(v *vm.DebuggerVM) error {
	// EXTCODECOPY requires four values on the stack: address, destOffset, offset, size
	if err := v.RequireStack(4); err != nil {
		return err
	}

	// Pop the four items from the stack (in reverse order)
	addrInt, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	destOffset, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	offset, err := v.Stack().Pop()
	if err != nil {
		return err
	}
	size, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// If size is 0, do nothing
	if size.IsZero() {
		return nil
	}

	// Convert address to 20-byte format
	var addr [20]byte
	addrBytes := addrInt.Bytes()
	if len(addrBytes) > 20 {
		// Take only the last 20 bytes if longer
		copy(addr[:], addrBytes[len(addrBytes)-20:])
	} else {
		// Right-align if shorter
		copy(addr[20-len(addrBytes):], addrBytes)
	}

	var code []byte
	if v.StateProvider != nil {
		// Get code from state provider
		code = v.StateProvider.GetCode(addr)
	} else {
		// If no state provider, treat as empty code
		code = []byte{}
	}

	// Extract the requested portion of code
	codeOffset := int(offset.Uint64())
	copySize := int(size.Uint64())

	var dataToWrite []byte
	if codeOffset >= len(code) {
		// If offset is beyond code, fill with zeros
		dataToWrite = make([]byte, copySize)
	} else {
		// Copy available code and pad with zeros if needed
		availableSize := len(code) - codeOffset
		if availableSize >= copySize {
			dataToWrite = code[codeOffset : codeOffset+copySize]
		} else {
			dataToWrite = make([]byte, copySize)
			copy(dataToWrite, code[codeOffset:])
			// Remaining bytes are already zero-initialized
		}
	}

	// Write to memory at destOffset
	v.Memory().Write(int(destOffset.Uint64()), dataToWrite)

	return nil
}
