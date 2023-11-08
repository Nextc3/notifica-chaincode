package chaincode

import (
	"encoding/json"
	"fmt"
	notifica_model "github.com/Nextc3/notifica-model"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"strconv"
)

/**
Esse Chaincode foi criado no assetransfer original tentando fazer o minímo de modificaçẽos e preservando
métodos originais
**/

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
/**
type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}
**/

var ultimoId int

//gerarObjetoTeste retorna dois objetos com valores de teste sendo que o segundo pode ser ignorado

func (s *SmartContract) GetUltimoId() int {
	return ultimoId
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	asset := notifica_model.Asset{
		Id:                  0,
		DocType:             "notificacao",
		DataNascimento:      "28/06/1988",
		Sexo:                "Masculino",
		Endereco:            "Av Cardeal Avelar Brandão Villela 6",
		Bairro:              "Jardim Santo Inácio",
		Cidade:              "Salvador",
		Estado:              "Bahia",
		Pais:                "Brasil",
		Doenca:              "Chagas",
		DataInicioSintomas:  "01/11/2023",
		DataDiagnostico:     "04/11/2023",
		DataNotificacao:     "08/11/2023",
		InformacoesClinicas: "tá ruim",
	}

	aEmBytes, _ := json.Marshal(asset)
	ultimoId++

	return s.CreateAsset(ctx, string(aEmBytes))

}

/*
*
// CreateAsset issues a new asset to the world state with given details.

	func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
		exists, err := s.AssetExists(ctx, id)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("the asset %s already exists", id)
		}

		asset := notifica_model.Asset{
			ID:             id,
			Color:          color,
			Size:           size,
			Owner:          owner,
			AppraisedValue: appraisedValue,
		}
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		return ctx.GetStub().PutState(id, assetJSON)
	}

*
*/
func (s *SmartContract) CreateAsset(contexto contractapi.TransactionContextInterface, asset string) error {

	/**exists, err := s.AssetExists(contexto, strconv.Itoa(id))
	  if err != nil {
	  	return err
	  }
	  if exists {
	  	return fmt.Errorf("the asset %s already exists", id)
	  }
	  **/

	assetEmBytes := []byte(asset)
	var a notifica_model.Asset
	_ = json.Unmarshal(assetEmBytes, &a)

	//Chave do estado é Asset + Id da asset
	//Cuidado para não salvar uma Asset com mesmo Id pois são utilizados para salvar na ledger
	err := contexto.GetStub().PutState("Asset"+strconv.Itoa(a.Id), assetEmBytes)
	if err != nil {
		log.Fatalf("Erro ao salvar na ledger %s", err)
	}

	return err
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(contexto contractapi.TransactionContextInterface, idAsset string) (*notifica_model.Asset, error) {
	assetEmBytes, err := contexto.GetStub().GetState("Asset" + idAsset)

	if err != nil {
		return nil, fmt.Errorf("Falha em consultar em Notificação na Ledger com GetState %s", err.Error())
	}

	if assetEmBytes == nil {
		return nil, fmt.Errorf("Notificacao%s não existe", idAsset)
	}

	asset := new(notifica_model.Asset)
	_ = json.Unmarshal(assetEmBytes, asset)

	return asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
/**
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}
**/
// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(contexto contractapi.TransactionContextInterface, idAsset string) (bool, error) {
	assetEmBytes, err := contexto.GetStub().GetState("Asset" + idAsset)
	if err != nil {
		return false, fmt.Errorf("falhou em consultar a existência da Notificacao: %v", err)
	}

	return assetEmBytes != nil, nil
}

/**
// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}
**/
// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*notifica_model.Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*notifica_model.Asset
	ultimoId = 0
	for resultsIterator.HasNext() {
		ultimoId++
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset notifica_model.Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
