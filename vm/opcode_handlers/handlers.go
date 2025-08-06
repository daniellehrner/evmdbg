package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

var handlers = map[vm.OpCode]vm.Handler{
	vm.STOP:           &StopOpCode{},
	vm.ADD:            &AddOpCode{},
	vm.MUL:            &MulOpCode{},
	vm.SUB:            &SubOpCode{},
	vm.DIV:            &DivOpCode{},
	vm.SDIV:           &SDivOpCode{},
	vm.MOD:            &ModOpCode{},
	vm.SMOD:           &SModOpCode{},
	vm.ADDMOD:         &AddModOpCode{},
	vm.MULMOD:         &MulModOpCode{},
	vm.EXP:            &ExpOpCode{},
	vm.SIGNEXTEND:     &SignExtendOpCode{},
	vm.LT:             &LtOpCode{},
	vm.GT:             &GtOpCode{},
	vm.SLT:            &SLTOpCode{},
	vm.SGT:            &SGTOpCode{},
	vm.EQ:             &EqOpCode{},
	vm.ISZERO:         &IsZeroOpCode{},
	vm.AND:            &AndOpCode{},
	vm.OR:             &OrOpCode{},
	vm.XOR:            &XorOpCode{},
	vm.NOT:            &NotOpCode{},
	vm.BYTE:           &ByteOpCode{},
	vm.SHL:            &ShlOpCode{},
	vm.SHR:            &ShrOpCode{},
	vm.SAR:            &SarOpCode{},
	vm.SHA3:           &Sha3OpCode{},
	vm.ADDRESS:        &AddressOpCode{},
	vm.BALANCE:        &BalanceOpCode{},
	vm.ORIGIN:         &OriginOpCode{},
	vm.CALLER:         &CallerOpCode{},
	vm.CALLVALUE:      &CallValueOpCode{},
	vm.CALLDATALOAD:   &CallDataLoadOpCode{},
	vm.CALLDATASIZE:   &CallDataSizeOpCode{},
	vm.CALLDATACOPY:   &CallDataCopyOpCode{},
	vm.CODESIZE:       &CodeSizeOpCode{},
	vm.CODECOPY:       &CodeCopyOpCode{},
	vm.GASPRICE:       &GasPriceOpCode{},
	vm.RETURNDATASIZE: &ReturnDataSizeOpCode{},
	vm.RETURNDATACOPY: &ReturnDataCopyOpCode{},
	vm.COINBASE:       &CoinbaseOpCode{},
	vm.TIMESTAMP:      &TimestampOpCode{},
	vm.NUMBER:         &NumberOpCode{},
	vm.DIFFICULTY:     &DifficultyOpCode{},
	vm.GASLIMIT:       &GasLimitOpCode{},
	vm.CHAINID:        &ChainIdOpCode{},
	vm.SELFBALANCE:    &SelfBalanceOpCode{},
	vm.BASEFEE:        &BaseFeeOpCode{},
	vm.POP:            &PopOpCode{},
	vm.MLOAD:          &MLoadOpCode{},
	vm.MSTORE:         &MStoreOpCode{},
	vm.MSTORE8:        &MStore8OpCode{},
	vm.SLOAD:          &SLoadOpCode{},
	vm.SSTORE:         &SStoreOpCode{},
	vm.JUMP:           &JumpOpCode{},
	vm.JUMPI:          &JumpiOpCode{},
	vm.PC:             &PCOpCode{},
	vm.MSIZE:          &MSizeOpCode{},
	vm.GAS:            &GasOpCode{},
	vm.JUMPDEST:       &JumpDestOpCode{},
	vm.MCOPY:          &MCopyOpCode{},
	vm.PUSH0:          &Push0OpCode{},
	vm.LOG0:           &LogNOpCode{N: 0},
	vm.LOG1:           &LogNOpCode{N: 1},
	vm.LOG2:           &LogNOpCode{N: 2},
	vm.LOG3:           &LogNOpCode{N: 3},
	vm.LOG4:           &LogNOpCode{N: 4},
	vm.CALL:           &CallOpCode{},
	vm.CALLCODE:       &CallCodeOpCode{},
	vm.DELEGATECALL:   &DelegateCallOpCode{},
	vm.STATICCALL:     &StaticCallOpCode{},
	vm.RETURN:         &ReturnOpCode{},
	vm.REVERT:         &RevertOpCode{},
	vm.INVALID:        &InvalidOpCode{},
}

func init() {
	// PUSH1 (0x60) to PUSH32 (0x7f)
	for i := 0; i < 32; i++ {
		op := vm.OpCode(0x60 + i)
		handlers[op] = &PushNOpCode{N: i + 1}
	}

	// DUP1 (0x80) to DUP16 (0x8f)
	for i := 1; i <= 16; i++ {
		handlers[vm.OpCode(0x7f+i)] = &DupOpCode{N: i}
	}

	// SWAP1 (0x90) to SWAP16 (0x9f)
	for i := 1; i <= 16; i++ {
		handlers[vm.OpCode(0x8f+i)] = &SwapOpCode{N: i}
	}
}

func GetHandler(b byte) vm.Handler {
	return handlers[vm.OpCode(b)]
}
