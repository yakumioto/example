package main

import (
    "log"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
    "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"

    mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

func main() {
    cp := config.FromFile("/data/swordartonline-network/config.yaml")
    sdk, err := fabsdk.New(cp)
    if err != nil {
        log.Fatalln(err)
    }

    //createChannel(sdk)

    //joinChannel(sdk)

    //createCC(sdk)

    queryCC(sdk)
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
        ChannelID: "master",
        ChannelConfigPath: "/data/swordartonline-network/channel-artifacts/master.tx",
        SigningIdentities: []msp.SigningIdentity{adminIdentity},
    }
    txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.aincrad.svc.cluster.local"))
    if err != nil {
        log.Fatalln("save channel error: ", err)
    }

    log.Println("create channel txID: ", txID.TransactionID)
}

func joinChannel(sdk *fabsdk.FabricSDK) {
    kiritoClientCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Kirito"))
    kiritoResMgmtClient, err :=resmgmt.New(kiritoClientCtx)
    if err != nil {
        log.Fatalln("resmgmt new error: ", err)
    }
    if err := kiritoResMgmtClient.JoinChannel("master", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.aincrad.svc.cluster.local")); err != nil {
        log.Fatalln("join channel error: ", err)
    }

    asunaClientCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Asuna"))
    asunaResMgmtClient, err := resmgmt.New(asunaClientCtx)
    if err != nil {
       log.Fatalln("resmgmt new error: ", err)
    }
    if err := asunaResMgmtClient.JoinChannel("master", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.aincrad.svc.cluster.local")); err != nil {
       log.Fatalln("join channel error: ", err)
    }
}

func createCC(sdk *fabsdk.FabricSDK) {
    kiritoClientCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Kirito"))
    kiritoResMgmtClient, err :=resmgmt.New(kiritoClientCtx)
    if err != nil {
        log.Fatalln("resmgmt new error: ", err)
    }

    ccpkg, err := gopackager.NewCCPackage("chaincode-example", "/data/swordartonline-network")
    ccreq := resmgmt.InstallCCRequest{
        Name: "example",
        Path: "chaincode-example",
        Version: "0",
        Package: ccpkg,
    }

    _, err = kiritoResMgmtClient.InstallCC(ccreq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
    if err != nil {
        log.Fatalln("install chaincode error: ", err)
    }

    ccPolicy := cauthdsl.SignedByAnyMember([]string{"Kirito"})
    resp, err := kiritoResMgmtClient.InstantiateCC("master",
        resmgmt.InstantiateCCRequest{
            Name: "example",
            Path: "chaincode-example",
            Version: "0",
            Args: [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")},
            Policy: ccPolicy,
        },
        resmgmt.WithRetry(retry.DefaultResMgmtOpts))
    if err != nil {
        log.Fatalln("instantiate chaincode error: ", err)
    }

    log.Println("create channel txID: ", resp.TransactionID)
}

func queryCC(sdk *fabsdk.FabricSDK) {
    kiritoClientCtx := sdk.ChannelContext("master", fabsdk.WithUser("Admin"), fabsdk.WithOrg("Kirito"))
    channelClient, err := channel.New(kiritoClientCtx)
    if err != nil {
        log.Fatalln("new channel error: ", err)
    }
    req := channel.Request{
        ChaincodeID: "example",
        Fcn: "query",
        Args: [][]byte{[]byte("a")},
    }
    res, err := channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints("peer1.kirito.svc.cluster.local"))
    if err != nil {
        log.Fatalln("query err: ", err)
    }

    log.Println("query cc:", string(res.Payload))
}
