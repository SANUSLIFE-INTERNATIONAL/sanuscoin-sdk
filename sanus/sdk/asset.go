package sdk

import (
	"fmt"

	"sanus/sanus-sdk/cc/asset"
	"sanus/sanus-sdk/entity"

	"github.com/btcsuite/btcd/txscript"
)

const (
	registeredAssetID = "La8e7WhGAEfiT9JGTmyPJopZhkRMwiEPz4uBEG"
)

func (w *BTCWallet) parseCCTx(cctx *asset.CCTransaction) error {
	for _, in := range cctx.Input {
		for _, _asset := range in.Assets {
			if _asset != nil {
				if err := w.db.AssetTransaction().Update(&entity.AssetTransactionRaw{
					AssetId: _asset.Id(),
					TxId:    cctx.Tx.TxHash().String(),
					Type:    cctx.Type,
				}); err != nil {
					return fmt.Errorf("can't update asset transaction. %v %v",
						cctx.Tx.TxHash().String(), err)
				}
			}
		}
	}
	assets, err := cctx.GetAssetOutput()
	if err != nil {
		return err
	}
	var index = 0
	for outIndex, _asset := range assets {
		if len(_asset) == 0 {
			continue
		}
		index = outIndex
		cctx.Output[outIndex].SetAssetsArray(_asset)
		if err := w.db.Utxo().Update(&entity.UtxoRaw{
			Assets: _asset,
			TxId:   cctx.Tx.TxHash().String(),
			Index:  outIndex,
		}); err != nil {
			return fmt.Errorf("can't update utxo. %v %v",
				cctx.Tx.TxHash().String(), err)
		}
		for _, oneAsset := range _asset {
			assetTxDB := w.db.AssetTransaction()
			if err := assetTxDB.Update(&entity.AssetTransactionRaw{
				AssetId: oneAsset.Id(),
				TxId:    cctx.Tx.TxHash().String(),
				Type:    cctx.Type,
			}); err != nil {
				return fmt.Errorf("can't update asset transaction. %v %v",
					cctx.Tx.TxHash().String(), err)
			}
			if err := w.db.AssetUtxo().Update(&entity.AssetUtxoRaw{
				AssetId: oneAsset.Id(),
				TxId:    cctx.Tx.TxHash().String(),
				Index:   outIndex,
			}); err != nil {
				return fmt.Errorf("can't update asset utxo. %v %v",
					cctx.Tx.TxHash().String(), err)
			}

			_, addrs, _, err := txscript.ExtractPkScriptAddrs(cctx.Output[outIndex].Output().PkScript, w.wlt.ChainParams())
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if err := w.db.AssetAddress().Update(addr.String(), oneAsset.Id()); err != nil {
					w.Logger.Errorf("error caused when trying to update into AssetAddress db %v", err.Error())
				}
			}
		}

	}
	for i := index + 1; i < len(cctx.Output); i++ {
		cctx.Output[i].SetAssetsArray(map[int]*asset.Asset{})
		utxoData := &entity.UtxoRaw{Index: i, TxId: cctx.Tx.TxHash().String()}

		if err := w.db.Utxo().ToEmptyByIndex(utxoData); err != nil {
			w.Logger.Errorf("error caused when trying to save into utxo db %v", err.Error())
		}
	}
	return nil
}
