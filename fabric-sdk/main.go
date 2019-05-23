package main

import (
    "log"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
    "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

    mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

func main() {
    cp := config.FromFile("./config.yaml")
    sdk, err := fabsdk.New(cp)
    if err != nil {
        log.Fatalln(err)
    }

    //createChannel(sdk)

    joinChannel(sdk)
}

func createChannel(sdk *fabsdk.FabricSDK) {
    clientCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Aincrad"))
    resMgmtClient, err :=resmgmt.New(clientCtx)
    if err != nil {
        log.Fatalln("resmgmt new error: ", err)
    }
    mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("Kirito"))
    if err != nil {
        log.Fatalln("msp new error: ", err)
    }
    adminIdentity, err := mspClient.GetSigningIdentity("Admin")
    req := resmgmt.SaveChannelRequest{
        ChannelID: "test",
        ChannelConfigPath: "/data/swordartonline-network/channel-artifacts/test.tx",
        SigningIdentities: []msp.SigningIdentity{adminIdentity},
    }
    txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.aincrad.moe"))
    if err != nil {
        log.Fatalln("save channel error: ", err)
    }

    log.Println("create channel txID: ", txID.TransactionID)
}

func joinChannel(sdk *fabsdk.FabricSDK) {
    clientCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Kirito"))
    resMgmtClient, err :=resmgmt.New(clientCtx)
    if err != nil {
       log.Fatalln("resmgmt new error: ", err)
    }
    if err := resMgmtClient.JoinChannel("master", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.aincrad.moe")); err != nil {
       log.Fatalln("join channel error: ", err)
    }
}
