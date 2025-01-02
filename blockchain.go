package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

type Block struct {
	PreviousHash [32]byte
	Timestamp    int64
	Transactions []*Transactions
	MerkleRoot   [32]byte
}

func NewBlock(previousHash [32]byte, transactions []*Transactions) *Block {
	b := new(Block)
	b.Timestamp = time.Now().UnixNano()
	b.PreviousHash = previousHash
	b.Transactions = transactions
	b.MerkleRoot = CalculateMerkleRoot(transactions)
	return b
}

func (b *Block) Print() {
	fmt.Printf("Timestamp:       %d\n", b.Timestamp)
	fmt.Printf("Previous Hash:   %x\n", b.PreviousHash)
	fmt.Printf("Merkle Root:     %x\n", b.MerkleRoot)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		PreviousHash [32]byte       `json:"previous_hash"`
		MerkleRoot   [32]byte       `json:"merkle_root"`
		Transactions []*Transactions `json:"transactions"`
	}{
		Timestamp:    b.Timestamp,
		PreviousHash: b.PreviousHash,
		MerkleRoot:   b.MerkleRoot,
		Transactions: b.Transactions,
	})
}

type Blockchain struct {
	TransactionPool []*Transactions
	Chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(b.Hash()) // Genesis block
	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.TransactionPool = append(bc.TransactionPool, t)
}

type Transactions struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transactions {
	return &Transactions{SenderBlockchainAddress: sender, RecipientBlockchainAddress: recipient, Value: value}
}

func (t *Transactions) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("Sender Blockchain Address:    %s\n", t.SenderBlockchainAddress)
	fmt.Printf("Recipient Blockchain Address: %s\n", t.RecipientBlockchainAddress)
	fmt.Printf("Value:                        %.1f\n", t.Value)
}

func (t *Transactions) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchainaddress"`
		Recipient string  `json:"recipient_blockchainaddress"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) CreateBlock(previousHash [32]byte) *Block {
	b := NewBlock(previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*Transactions{}
	return b
}

func CalculateMerkleRoot(transactions []*Transactions) [32]byte {
	if len(transactions) == 0 {
		return [32]byte{}
	}

	hashes := make([][32]byte, len(transactions))
	for i, t := range transactions {
		hashes[i] = t.Hash()
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		newLevel := make([][32]byte, 0)
		for i := 0; i < len(hashes); i += 2 {
			h := sha256.Sum256(append(hashes[i][:], hashes[i+1][:]...))
			newLevel = append(newLevel, h)
		}
		hashes = newLevel
	}

	return hashes[0]
}

func (t *Transactions) Hash() [32]byte {
	m, _ := json.Marshal(t)
	return sha256.Sum256(m)
}

func main() {
	blockchain := NewBlockchain()
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 1.0)
	previousHash := blockchain.LastBlock().Hash()
	blockchain.CreateBlock(previousHash)
	blockchain.Print()

	blockchain.AddTransaction("C", "D", 2.0)
	previousHash = blockchain.LastBlock().Hash()
	blockchain.CreateBlock(previousHash)
	blockchain.Print()
}
