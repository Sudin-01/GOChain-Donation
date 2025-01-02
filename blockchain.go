package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func init(){
	log.SetPrefix("Blockchain: ")
}


type Block struct {
	nounce int
	previousHash [32]byte
	timestamp int64
	transactions []*Transactions
}

func NewBlock(nounce int, previousHash [32]byte, transactions []*Transactions) *Block{ //NewBlock function to create a new block
	
	b := new(Block)
	b.timestamp= time.Now().UnixNano()
	b.nounce=nounce
	b.previousHash=previousHash
	b.transactions=transactions
	return b
	
}

func (b *Block) Print(){ //Print function to print the block
	fmt.Printf("timestamp        %d\n", b.timestamp)	
	fmt.Printf("nounce           %d\n", b.nounce)	
	fmt.Printf("previous_hash    %x\n", b.previousHash)	
	// fmt.Printf("transactions     %s\n", b.transactions)
	for _,t :=range b.transactions{
		t.Print()
	}
}

func (b *Block) Hash() [32]byte{
	m,_ :=json.Marshal(b) //Marshal is used to convert struct to json
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte,error){ //CUsotm Marshal function so as to show the json in a better way
	return json.Marshal(struct{
		Timestamp int64 			`json:"timestamp"`
		Nounce int 					`json:"nounce"`
		PreviousHash [32]byte 		`json:"previous_hash"`
		Transactions []*Transactions `json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nounce: b.nounce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

type Blockchain struct{ //Blockchain struct
	transactionPool []*Transactions
	chain 			[]*Block
}

func NewBlockchain() *Blockchain{ //NewBlockchain function to create a new blockchain
	b:=&Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash()) //Geneisi block
	return bc
}

func (bc *Blockchain) LastBlock() *Block{ //LastBlock function to get the last block
	return bc.chain[len(bc.chain)-1]
}

func(bc *Blockchain) AddTransaction(sender string, receipient string, value float32){ //AddTransaction function to add a transaction to the blockchain
	t := NewTransaction(sender,receipient,value)
	bc.transactionPool = append(bc.transactionPool,t)
}



type Transactions struct{
	senderBlockChainAddress string
	receipientBlockChainAddress string
	value 						float32

}

func NewTransaction(sender string, receipient string, value float32) *Transactions{
	return &Transactions{sender,receipient,value}
}

func (t *Transactions) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("Sender Blockchain Address:    %s\n", t.senderBlockChainAddress)
	fmt.Printf("Recipient Blockchain Address: %s\n", t.receipientBlockChainAddress)
	fmt.Printf("Value:                        %.1f\n", t.value)
}


func (t *Transactions) MarshalJSON() ([]byte,error){
	return json.Marshal(struct{
		Sender string 		`json:"sender_blockchainaddress"`
		Receipient string 	`json:"receipient_blockchainaddress"`
		Value float32 		`json:"value"`
	}{
		Sender: 	t.senderBlockChainAddress,
		Receipient: t.receipientBlockChainAddress,
		Value: 		t.value,
	})
}


func (bc *Blockchain) Print(){ //Print function to print the blockchain
	for i,block := range bc.chain{
		fmt.Printf("%s Chain %d %s \n", strings.Repeat("=",25),i,
		strings.Repeat("=",25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*",25))
}

func (bc *Blockchain) CreateBlock(nounce int,previousHash [32]byte) *Block{
	b := NewBlock(nounce,previousHash,bc.transactionPool)
	bc.chain = append(bc.chain,b)
	bc.transactionPool = []*Transactions{}
	return b
}

func main(){


	blockchain := NewBlockchain()
	blockchain.Print()
 
	blockchain.AddTransaction("A","B",1.0)
	previousHash:=blockchain.LastBlock().Hash()
	blockchain.CreateBlock(5,previousHash)
	blockchain.Print()

	blockchain.AddTransaction("X","Y",2.0)
	blockchain.AddTransaction("P","Q",3.0)
	previousHash=blockchain.LastBlock().Hash()
	blockchain.CreateBlock(2,previousHash)
	blockchain.Print()


}