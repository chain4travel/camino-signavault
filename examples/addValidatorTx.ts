import { Avalanche, BinTools, BN, Buffer } from "./dependencies/caminojs/dist/index"
import {
    PlatformVMAPI,
    KeyChain,
    UTXOSet,
    UnsignedTx,
    Tx,
    PlatformVMConstants
} from "./dependencies/caminojs/dist/apis/platformvm"
import {
    MultisigKeyChain,
    MultisigKeyPair,
    OutputOwners
} from "dependencies/caminojs/dist/common"
import {
    PrivateKeyPrefix,
    DefaultLocalGenesisPrivateKey,
    DefaultLocalGenesisPrivateKey2,
    PChainAlias
} from "dependencies/caminojs/dist/utils"
import { ExamplesConfig } from "./common/examplesConfig"
import {
    Configuration,
    DtoIssueTxResponse,
    ModelMultisigTx,
    MultisigApi
} from "@c4tplatform/signavaultjs"
import { AxiosResponse } from "axios"
import createHash from "create-hash";

const config: ExamplesConfig = require("./common/examplesConfig.json")
const bintools = BinTools.getInstance()
const avalanche: Avalanche = new Avalanche(
    config.host,
    config.port,
    config.protocol,
    config.networkID
)

const signavaultConfig: Configuration = new Configuration({
    basePath:
        config.signavaultProtocol +
        "://" +
        config.signavaultHost +
        ":" +
        config.signavaultPort +
        "/v1"
})
const signavault: MultisigApi = new MultisigApi(signavaultConfig)

const privKey1: string = `${PrivateKeyPrefix}${DefaultLocalGenesisPrivateKey}`
const privKey2: string = `${PrivateKeyPrefix}${DefaultLocalGenesisPrivateKey2}`
const nodePrivKey: string =
    "PrivateKey-2ZW6HUePBW2dP7dBGa5stjXe1uvK9LwEgrjebDwXEyL5bDMWWS"
const nodeID: string = "NodeID-D1LbWvUf9iaeEyUbTYYtYq4b7GaYR5tnJ"
const asOf: BN = new BN(0)
const threshold: number = 1
const msigAlias = "P-kopernikus1t5qgr9hcmf2vxj7k0hz77kawf9yr389cxte5j0"

// msig definition that needs to be present on chain in order for the example to work
// {
//   "memo": "222",
//     "alias": "X-kopernikus1t5qgr9hcmf2vxj7k0hz77kawf9yr389cxte5j0",
//     "addresses": [
//         "X-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
//         "X-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3"
//     ],
//     "threshold": 2
// }

let pchain: PlatformVMAPI
let pKeychain: KeyChain
let pAddresses: Buffer[]
let pAddressStrings: string[]

const InitAvalanche = async () => {
    await avalanche.fetchNetworkSettings()
    pchain = avalanche.PChain()
    pKeychain = pchain.keyChain()
    pKeychain.importKey(privKey1)
    pKeychain.importKey(privKey2)
    pKeychain.importKey(nodePrivKey)
    pAddresses = pchain.keyChain().getAddresses()
    pAddressStrings = pchain.keyChain().getAddressStrings()
}

const sendRegisterNodeTx = async (): Promise<any> => {
    // Those are not serialized back and forth because
    // its so simple and has no methods
    let signatures: [string, string][] = []

    // these are serialized to test if their methods are
    // working properly
    let outputOwnersHex: string = ""

    // simulate tx creation

    const owner = await pchain.getMultisigAlias(msigAlias)

    const platformVMUTXOResponse: any = await pchain.getUTXOs([msigAlias])
    const utxoSet: UTXOSet = platformVMUTXOResponse.utxos

    const unsignedTx: UnsignedTx = await pchain.buildRegisterNodeTx(
        utxoSet,
        [[msigAlias], pAddressStrings],
        [msigAlias],
        undefined,
        nodeID,
        msigAlias,
        [[0, msigAlias]],
        undefined,
        asOf,
        owner.threshold
    )

    // turn it into a hex blob
    const txbuff = unsignedTx.toBuffer()
    const msg: Buffer = Buffer.from(
        createHash("sha256").update(txbuff).digest()
    )

    outputOwnersHex = OutputOwners.toArray(
        unsignedTx.getTransaction().getOutputOwners()
    ).toString("hex")

    // simulate signing
    {
        for (let address of pAddresses) {
            // We need the keychain for signing
            const keyPair = pKeychain.getKey(address)
            // The signature
            const signature = keyPair.sign(msg)
            // save the signature
            signatures.push([keyPair.getAddressString(), signature.toString("hex")])
        }
    }

    // simulate the first signer
    // create multisig tx call to signavault
    try {
        await signavault.createMultisigTx({
            alias: msigAlias,
            unsignedTx: txbuff.toString('hex'),
            signature: signatures[0][1],
            outputOwners: outputOwnersHex,
            // we send node's signature as metadata so it can be used form the issuer
            metadata: signatures[2][1]
        })
    } catch (e) {
        console.log(e.response.data)
        return
    }

    // simulate the second signer
    // compose signature from alias and timestamp
    const timestamp = Math.floor(Date.now() / 1000).toString()
    const keyPair1 = pKeychain.getKey(pAddresses[1])
    const signatureAliasTimestamp = keyPair1
        .sign(
            Buffer.from(
                createHash("sha256")
                    .update(Buffer.from(msigAlias + timestamp))
                    .digest()
            )
        )
        .toString("hex")

    // get all pending multisig txs from signavault
    let pendingMultisigTxs: AxiosResponse<Array<ModelMultisigTx>>
    try {
        pendingMultisigTxs = await signavault.getAllMultisigTxForAlias(
            msigAlias,
            signatureAliasTimestamp,
            timestamp
        )
    } catch (e) {
        console.log(e.response.data)
        return
    }

    try {
        await signavault.signMultisigTx(pendingMultisigTxs.data[0].id, {
            signature: signatures[1][1]
        })
    } catch (e) {
        console.log(e.response.data)
        return
    }

    // get all pending multisig txs from signavault
    try {
        pendingMultisigTxs = await signavault.getAllMultisigTxForAlias(
            msigAlias,
            signatureAliasTimestamp,
            timestamp
        )
    } catch (e) {
        console.log(e.response.data)
        return
    }
    const pendingMultisigTx = pendingMultisigTxs.data[0]

    // reconstruciton
    {
        // load msig configuration from node
        const msigAliasBuffer = pchain.parseAddress(msigAlias)
        const owner = await pchain.getMultisigAlias(msigAlias)

        // deserialize
        let unsignedTx = new UnsignedTx()
        unsignedTx.fromBuffer(Buffer.from(pendingMultisigTx.unsignedTx, "hex"))

        // parse and set output owners - are requried for msig resolving
        let parsedOwners: OutputOwners[] = OutputOwners.fromArray(
            Buffer.from(pendingMultisigTx.outputOwners, "hex")
        )
        unsignedTx.getTransaction().setOutputOwners(parsedOwners)

        const txbuff = unsignedTx.toBuffer()
        const msg: Buffer = Buffer.from(
            createHash("sha256").update(txbuff).digest()
        )

        // create MSKeychein to create proper signidx
        const msKeyChain = new MultisigKeyChain(
            avalanche.getHRP(),
            PChainAlias,
            msg,
            PlatformVMConstants.SECPMULTISIGCREDENTIAL,
            unsignedTx.getTransaction().getOutputOwners(),
            new Map([
                [
                    msigAliasBuffer.toString("hex"),
                    new OutputOwners(
                        owner.addresses.map((a) => bintools.parseAddress(a, "P")),
                        new BN(owner.locktime),
                        owner.threshold
                    )
                ]
            ])
        )

        for (let Owner of pendingMultisigTx.owners) {
            let address = pchain.parseAddress(Owner.address)
            let signature = Buffer.from(Owner.signature, "hex")
            msKeyChain.addKey(new MultisigKeyPair(msKeyChain, address, signature))
        }

        // add node's signature taken from metadata
        let nodeAddress = pchain.parseAddress(signatures[2][0])
        let nodeSignature = Buffer.from(pendingMultisigTx.metadata, "hex")
        msKeyChain.addKey(
            new MultisigKeyPair(msKeyChain, nodeAddress, nodeSignature)
        )

        // build signature indices
        msKeyChain.buildSignatureIndices()

        // sign the tx
        const tx: Tx = unsignedTx.sign(msKeyChain)

        // send tx to node through signavault
        // compose signature from signedTx
        const keyPair1 = pKeychain.getKey(pAddresses[1])
        const signatureOfSignedTx = keyPair1
            .sign(Buffer.from(createHash("sha256").update(tx.toBuffer()).digest()))
            .toString("hex")

        // issue multisig tx call to signavault
        let txIDResponse: AxiosResponse<DtoIssueTxResponse>
        try {
            txIDResponse = await signavault.issueMultisigTx({
                signature: signatureOfSignedTx,
                signedTx: tx.toBuffer().toString("hex")
            })
        } catch (e) {
            console.log(e.response.data)
            return
        }

        console.log(`Success! TXID: ${txIDResponse.data.txID}`)
    }

}

const sendAddValidatorTx = async (): Promise<any> => {
    // Those are not serialized back and forth because
    // its so simple and has no methods
    let signatures: [string, string][] = []

    // these are serialized to test if their methods are
    // working properly
    let outputOwnersHex: string = ""

    // simulate tx creation

    const owner = await pchain.getMultisigAlias(msigAlias)

    const platformVMUTXOResponse: any = await pchain.getUTXOs([msigAlias])
    const utxoSet: UTXOSet = platformVMUTXOResponse.utxos

    let startDate = new Date(Date.now() + 0.5 * 60 * 1000).getTime() / 1000
    let endDate = startDate + 60 * 60 * 24 * 10

    const unsignedTx: UnsignedTx = await pchain.buildAddValidatorTx(
        utxoSet,
        [msigAlias],
        [[msigAlias], pAddressStrings],
        [msigAlias],
        nodeID,
        new BN(startDate),
        new BN(endDate),
        new BN(2000000000000),
        [msigAlias],
        0, // delegation fee
        undefined,
        threshold,
        undefined,
        asOf,
        owner.threshold,
        owner.threshold
    )

    // turn it into a hex blob
    const txbuff = unsignedTx.toBuffer()
    const msg: Buffer = Buffer.from(
        createHash("sha256").update(txbuff).digest()
    )

    outputOwnersHex = OutputOwners.toArray(
        unsignedTx.getTransaction().getOutputOwners()
    ).toString("hex")

    // simulate signing
    {
        for (let address of pAddresses) {
            // We need the keychain for signing
            const keyPair = pKeychain.getKey(address)
            // The signature
            const signature = keyPair.sign(msg)
            // save the signature
            signatures.push([keyPair.getAddressString(), signature.toString("hex")])
        }
    }

    // simulate the first signer
    // create multisig tx call to signavault
    try {
        await signavault.createMultisigTx({
            alias: msigAlias,
            unsignedTx: txbuff.toString('hex'),
            signature: signatures[0][1],
            outputOwners: outputOwnersHex,
            // we send node's signature as metadata so it can be used form the issuer
            metadata: signatures[2][1]
        })
    } catch (e) {
        console.log(e.response.data)
        return
    }

    // simulate the second signer
    // compose signature from alias and timestamp
    const timestamp = Math.floor(Date.now() / 1000).toString()
    const keyPair1 = pKeychain.getKey(pAddresses[1])
    const signatureAliasTimestamp = keyPair1
        .sign(
            Buffer.from(
                createHash("sha256")
                    .update(Buffer.from(msigAlias + timestamp))
                    .digest()
            )
        )
        .toString("hex")

    // get all pending multisig txs from signavault
    let pendingMultisigTxs: AxiosResponse<Array<ModelMultisigTx>>
    try {
        pendingMultisigTxs = await signavault.getAllMultisigTxForAlias(
            msigAlias,
            signatureAliasTimestamp,
            timestamp
        )
    } catch (e) {
        console.log(e.response.data)
        return
    }

    try {
        await signavault.signMultisigTx(pendingMultisigTxs.data[0].id, {
            signature: signatures[1][1]
        })
    } catch (e) {
        console.log(e.response.data)
        return
    }

    // get all pending multisig txs from signavault
    try {
        pendingMultisigTxs = await signavault.getAllMultisigTxForAlias(
            msigAlias,
            signatureAliasTimestamp,
            timestamp
        )
    } catch (e) {
        console.log(e.response.data)
        return
    }
    const pendingMultisigTx = pendingMultisigTxs.data[0]


    // simulate reconstruciton
    {
        // load msig configuration from node
        const msigAliasBuffer = pchain.parseAddress(msigAlias)
        const owner = await pchain.getMultisigAlias(msigAlias)

        // deserialize
        let unsignedTx = new UnsignedTx()
        unsignedTx.fromBuffer(Buffer.from(pendingMultisigTx.unsignedTx, "hex"))

        // parse and set output owners - are requried for msig resolving
        let parsedOwners: OutputOwners[] = OutputOwners.fromArray(
            Buffer.from(outputOwnersHex, "hex")
        )
        unsignedTx.getTransaction().setOutputOwners(parsedOwners)

        const txbuff = unsignedTx.toBuffer()
        const msg: Buffer = Buffer.from(
            createHash("sha256").update(txbuff).digest()
        )

        // create MSKeychein to create proper signidx
        const msKeyChain = new MultisigKeyChain(
            avalanche.getHRP(),
            PChainAlias,
            msg,
            PlatformVMConstants.SECPMULTISIGCREDENTIAL,
            unsignedTx.getTransaction().getOutputOwners(),
            new Map([
                [
                    msigAliasBuffer.toString("hex"),
                    new OutputOwners(
                        owner.addresses.map((a) => bintools.parseAddress(a, "P")),
                        new BN(owner.locktime),
                        owner.threshold
                    )
                ]
            ])
        )

        for (let Owner of pendingMultisigTx.owners) {
            let address = pchain.parseAddress(Owner.address)
            let signature = Buffer.from(Owner.signature, "hex")
            msKeyChain.addKey(new MultisigKeyPair(msKeyChain, address, signature))
        }

        // add node's signature taken from metadata
        let nodeAddress = pchain.parseAddress(signatures[2][0])
        let nodeSignature = Buffer.from(pendingMultisigTx.metadata, "hex")
        msKeyChain.addKey(
            new MultisigKeyPair(msKeyChain, nodeAddress, nodeSignature)
        )

        // build signature indices
        msKeyChain.buildSignatureIndices()

        // sign the tx
        const tx: Tx = unsignedTx.sign(msKeyChain)

        // send tx to node through signavault
        // compose signature from signedTx
        const keyPair1 = pKeychain.getKey(pAddresses[1])
        const signatureOfSignedTx = keyPair1
            .sign(Buffer.from(createHash("sha256").update(tx.toBuffer()).digest()))
            .toString("hex")

        // issue multisig tx call to signavault
        let txIDResponse: AxiosResponse<DtoIssueTxResponse>
        try {
            txIDResponse = await signavault.issueMultisigTx({
                signature: signatureOfSignedTx,
                signedTx: tx.toBuffer().toString("hex")
            })
        } catch (e) {
            console.log(e.response.data)
            return
        }

        console.log(`Success! TXID: ${txIDResponse.data.txID}`)
    }
}


const main = async (): Promise<any> => {
    await InitAvalanche()
    try {
        await sendRegisterNodeTx()
        await sendAddValidatorTx()
    } catch (e) {
        console.log(e)
    }
}

main()
