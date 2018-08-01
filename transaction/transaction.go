// Loslly following https://en.bitcoin.it/wiki/Transaction

package transaction

// The default reward for our loyal miner
//TODO: add some algorithm to change this as time passes by :D
const minerReward = 50

//Each transaction is constructed of several inputs and outputs
type Transaction struct {
	inCounter  int
	outCounter int
	in         []TXInput
	out        []TXOutput
}

func New(inputs []TXInput, outputs []TXOutput) *Transaction{
	//transaction.SetID()
	return &Transaction{len(inputs), len(outputs), inputs, outputs}
}

type TXInput struct {
	prevTransactionHash []byte
	prevTxOutIndex      int
	txInScript          []byte
	txInScriptLength    int
}

type TXOutput struct {
	value             int
	txOutScript       []byte
	txOutScriptLength int
}
//Simple coinbase (first block transaction) generation transaction generator with no regard to script
func NewCoinbaseTransaction() *Transaction {
	txin := TXInput{[]byte{}, -1, []byte(""),len("")}
	txout := TXOutput{minerReward, []byte(""), len("")}
	transaction := New([]TXInput{txin}, []TXOutput{txout})

	return transaction
}