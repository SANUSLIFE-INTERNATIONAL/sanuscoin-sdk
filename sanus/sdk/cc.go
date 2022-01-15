package sdk

import (
	"fmt"
	"time"

	"sanus/sanus-sdk/cc/asset"
	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/entity"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	issuanceTypeName = "issuance"
	transferTypeName = "Transfer"

	defaultStartBlock = 1610203
)

func (w *BTCWallet) Scan() error {
	asset.InitENCLookup()
	var startBlock = defaultStartBlock
	lastBlockDB := w.db.LastBlockDB()
	lastBlockentity, err := lastBlockDB.GetLastIndex()
	if err != nil {
		w.Logger.Errorf("error caused when trying to get last processed block %v", err.Error())
	} else {
		startBlock = lastBlockentity.Index
	}
	blockchain, err := w.rpcClient.GetBlockChainInfo()
	if err != nil {
		return err
	}
	w.Logger.Infof("Scanning started from %v and will finished at %v\n", startBlock, blockchain.Blocks)

	for blockHeight := startBlock; blockHeight < int(blockchain.Blocks); blockHeight++ {
		hash, err := w.rpcClient.GetBlockHash(int64(blockHeight))
		if err != nil {
			return err
		}
		block, err := w.rpcClient.GetBlock(hash)
		if err != nil {
			return err
		}

		for _, transaction := range block.Transactions {
			cctx := w.toCCTransaction(transaction)
			if cctx.Type == "Transfer" || cctx.Type == "issuance" {
				if err = w.parseCCTx(cctx); err != nil {
					w.Logger.Errorf("error caused when trying to parse CCTX", err)
					continue
				}
				rawTxEntity := entity.NewRawTransactionEntity(cctx)

				if err := w.db.RawTransaction().Update(rawTxEntity); err != nil {
					w.Logger.Errorf("error caused when trying to save raw transaction")
					continue
				} else {
					w.Logger.Infof("%v transaction saved", transaction.TxHash().String())
				}

			}
		}
		lastBlock := &entity.LastBlockEntity{
			Index: blockHeight,
			Hash:  "",
		}
		if err = w.db.LastBlockDB().Update(lastBlock); err != nil {
			w.Logger.Errorf("error caused when trying to save last block %v", err.Error())
		}

	}
	// waiting for a minute to restart scanning
	time.Sleep(1 * time.Minute)

	return nil
}

func (w *BTCWallet) toCCTransaction(tx *wire.MsgTx) *asset.CCTransaction {
	cctx := &asset.CCTransaction{
		Tx: tx, Type: "null",
		Output: map[int]*asset.CCVout{},
		Input:  map[int]*asset.CCVin{},
	}

	for _, in := range tx.TxIn {
		cctx.AppendInput(&asset.CCVin{Input: in, Assets: map[int]*asset.Asset{}})
	}
	for _, out := range tx.TxOut {
		cctx.AppendOutput(&asset.CCVout{Out: out, Assets: map[int]*asset.Asset{}})
		if len(out.PkScript) > 2 && out.PkScript[0] == txscript.OP_RETURN {
			script := out.PkScript[2:]

			if len(script) < 4 {
				continue
			}
			encodeType := script[3]
			encoder, ok := asset.EncLookup[encodeType]
			if !ok {
				continue
			}
			switch encoder {
			case "issuance":
				ccData, err := issuance.Decode(script)
				if err != nil {
					continue
				}
				cctx.Issuance = ccData
				cctx.Type = issuanceTypeName
			case "Transfer":
				fmt.Println("transfer decode")
				ccData, err := transfer.Decode(script)
				if err != nil {
					fmt.Println(err)
					continue
				}
				cctx.Transfer = ccData
				cctx.Type = transferTypeName
			}
		}
	}
	return cctx
}
