package interview

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/v2/x/cert/keeper"
	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

func InitDefaultGenesis(ctx sdk.Context, k keeper.Keeper) {
	InitGenesis(ctx, k, *types.DefaultGenesisState())
}

// InitGenesis initialize default parameters and the keeper's address to pubkey map.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	certifiers := data.Certifiers
	platforms := data.Platforms
	certificates := data.Certificates
	libraries := data.Libraries
	nextCertificateID := data.NextCertificateId

	for _, certifier := range certifiers {
		k.SetCertifier(ctx, certifier)
	}
	if len(certifiers) > 0 {
		certifierAddr, err := sdk.AccAddressFromBech32(certifiers[0].Address)
		if err != nil {
			panic(err)
		}
		for _, platform := range platforms {
			pk, ok := platform.ValidatorPubkey.GetCachedValue().(cryptotypes.PubKey)
			if !ok {
				panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", platform.ValidatorPubkey))
			}

			_ = k.CertifyPlatform(ctx, certifierAddr, pk, platform.Description)
		}
	}
	for _, certificate := range certificates {
		k.SetCertificate(ctx, certificate)
	}
	for _, library := range libraries {
		libAddr, err := sdk.AccAddressFromBech32(library.Address)
		if err != nil {
			panic(err)
		}
		publisherAddr, err := sdk.AccAddressFromBech32(library.Publisher)
		if err != nil {
			panic(err)
		}
		k.SetLibrary(ctx, libAddr, publisherAddr)
	}
	k.SetNextCertificateID(ctx, nextCertificateID)
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	certifiers := k.GetAllCertifiers(ctx)
	platforms := k.GetAllPlatforms(ctx)
	certificates := k.GetAllCertificates(ctx)
	libraries := k.GetAllLibraries(ctx)
	nextCertificateID := k.GetNextCertificateID(ctx)

	return &types.GenesisState{
		Certifiers:        certifiers,
		Platforms:         platforms,
		Certificates:      certificates,
		Libraries:         libraries,
		NextCertificateId: nextCertificateID,
	}
}
