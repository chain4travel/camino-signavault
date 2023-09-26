/* tslint:disable */
/* eslint-disable */
/**
 * Signavault API
 * This is the signavault API.
 *
 * The version of the OpenAPI document: 1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import type { Configuration } from './configuration';
import type { AxiosPromise, AxiosInstance, AxiosRequestConfig } from 'axios';
import globalAxios from 'axios';
// Some imports not used depending on template conditions
// @ts-ignore
import { DUMMY_BASE_URL, assertParamExists, setApiKeyToObject, setBasicAuthToObject, setBearerAuthToObject, setOAuthToObject, setSearchParams, serializeDataIfNeeded, toPathString, createRequestFunction } from './common';
import type { RequestArgs } from './base';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, BaseAPI, RequiredError } from './base';

/**
 * 
 * @export
 * @interface DtoAddSignatureArgs
 */
export interface DtoAddSignatureArgs {
    /**
     * 
     * @type {string}
     * @memberof DtoAddSignatureArgs
     */
    'address': string;
    /**
     * 
     * @type {string}
     * @memberof DtoAddSignatureArgs
     */
    'depositOfferID': string;
    /**
     * 
     * @type {string}
     * @memberof DtoAddSignatureArgs
     */
    'signature': string;
    /**
     * 
     * @type {number}
     * @memberof DtoAddSignatureArgs
     */
    'timestamp'?: number;
}
/**
 * 
 * @export
 * @interface DtoCancelTxArgs
 */
export interface DtoCancelTxArgs {
    /**
     * 
     * @type {string}
     * @memberof DtoCancelTxArgs
     */
    'id': string;
    /**
     * 
     * @type {string}
     * @memberof DtoCancelTxArgs
     */
    'signature': string;
    /**
     * 
     * @type {string}
     * @memberof DtoCancelTxArgs
     */
    'timestamp': string;
}
/**
 * 
 * @export
 * @interface DtoIssueTxArgs
 */
export interface DtoIssueTxArgs {
    /**
     * 
     * @type {string}
     * @memberof DtoIssueTxArgs
     */
    'signature': string;
    /**
     * 
     * @type {string}
     * @memberof DtoIssueTxArgs
     */
    'signedTx': string;
}
/**
 * 
 * @export
 * @interface DtoIssueTxResponse
 */
export interface DtoIssueTxResponse {
    /**
     * 
     * @type {string}
     * @memberof DtoIssueTxResponse
     */
    'txID': string;
}
/**
 * 
 * @export
 * @interface DtoMultisigTxArgs
 */
export interface DtoMultisigTxArgs {
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'alias': string;
    /**
     * 
     * @type {number}
     * @memberof DtoMultisigTxArgs
     */
    'expiration'?: number;
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'metadata'?: string;
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'outputOwners': string;
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'parentTransaction'?: string;
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'signature': string;
    /**
     * 
     * @type {string}
     * @memberof DtoMultisigTxArgs
     */
    'unsignedTx': string;
}
/**
 * 
 * @export
 * @interface DtoSignTxArgs
 */
export interface DtoSignTxArgs {
    /**
     * 
     * @type {string}
     * @memberof DtoSignTxArgs
     */
    'signature': string;
}
/**
 * 
 * @export
 * @interface DtoSignavaultError
 */
export interface DtoSignavaultError {
    /**
     * 
     * @type {string}
     * @memberof DtoSignavaultError
     */
    'error': string;
    /**
     * 
     * @type {string}
     * @memberof DtoSignavaultError
     */
    'message': string;
}
/**
 * 
 * @export
 * @interface ModelDepositOfferSig
 */
export interface ModelDepositOfferSig {
    /**
     * 
     * @type {string}
     * @memberof ModelDepositOfferSig
     */
    'address'?: string;
    /**
     * 
     * @type {string}
     * @memberof ModelDepositOfferSig
     */
    'depositOfferID'?: string;
    /**
     * 
     * @type {string}
     * @memberof ModelDepositOfferSig
     */
    'signature'?: string;
}
/**
 * 
 * @export
 * @interface ModelMultisigTx
 */
export interface ModelMultisigTx {
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'alias': string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'chainId': string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'expiration'?: string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'id': string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'metadata'?: string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'outputOwners': string;
    /**
     * 
     * @type {Array<ModelMultisigTxOwner>}
     * @memberof ModelMultisigTx
     */
    'owners': Array<ModelMultisigTxOwner>;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'parentTransaction'?: string;
    /**
     * 
     * @type {number}
     * @memberof ModelMultisigTx
     */
    'threshold': number;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'timestamp': string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'transactionId'?: string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTx
     */
    'unsignedTx': string;
}
/**
 * 
 * @export
 * @interface ModelMultisigTxOwner
 */
export interface ModelMultisigTxOwner {
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTxOwner
     */
    'address': string;
    /**
     * 
     * @type {string}
     * @memberof ModelMultisigTxOwner
     */
    'signature'?: string;
}

/**
 * DepositOfferApi - axios parameter creator
 * @export
 */
export const DepositOfferApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Adds a signature mapped to a deposit offer id and an address
         * @param {DtoAddSignatureArgs} addSignatureArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        addSignature: async (addSignatureArgs: DtoAddSignatureArgs, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'addSignatureArgs' is not null or undefined
            assertParamExists('addSignature', 'addSignatureArgs', addSignatureArgs)
            const localVarPath = `/deposit-offer`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(addSignatureArgs, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Retrieves all signatures for an address only for authorized calls.
         * @param {string} address Address for which to retrieve all signatures
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSignatures: async (address: string, signature: string, timestamp: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'address' is not null or undefined
            assertParamExists('getSignatures', 'address', address)
            // verify required parameter 'signature' is not null or undefined
            assertParamExists('getSignatures', 'signature', signature)
            // verify required parameter 'timestamp' is not null or undefined
            assertParamExists('getSignatures', 'timestamp', timestamp)
            const localVarPath = `/deposit-offer/{address}`
                .replace(`{${"address"}}`, encodeURIComponent(String(address)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (signature !== undefined) {
                localVarQueryParameter['signature'] = signature;
            }

            if (timestamp !== undefined) {
                localVarQueryParameter['timestamp'] = timestamp;
            }


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * DepositOfferApi - functional programming interface
 * @export
 */
export const DepositOfferApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = DepositOfferApiAxiosParamCreator(configuration)
    return {
        /**
         * 
         * @summary Adds a signature mapped to a deposit offer id and an address
         * @param {DtoAddSignatureArgs} addSignatureArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async addSignature(addSignatureArgs: DtoAddSignatureArgs, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<void>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.addSignature(addSignatureArgs, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * 
         * @summary Retrieves all signatures for an address only for authorized calls.
         * @param {string} address Address for which to retrieve all signatures
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getSignatures(address: string, signature: string, timestamp: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelDepositOfferSig>>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getSignatures(address, signature, timestamp, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
    }
};

/**
 * DepositOfferApi - factory interface
 * @export
 */
export const DepositOfferApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = DepositOfferApiFp(configuration)
    return {
        /**
         * 
         * @summary Adds a signature mapped to a deposit offer id and an address
         * @param {DtoAddSignatureArgs} addSignatureArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        addSignature(addSignatureArgs: DtoAddSignatureArgs, options?: any): AxiosPromise<void> {
            return localVarFp.addSignature(addSignatureArgs, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Retrieves all signatures for an address only for authorized calls.
         * @param {string} address Address for which to retrieve all signatures
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSignatures(address: string, signature: string, timestamp: string, options?: any): AxiosPromise<Array<ModelDepositOfferSig>> {
            return localVarFp.getSignatures(address, signature, timestamp, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * DepositOfferApi - object-oriented interface
 * @export
 * @class DepositOfferApi
 * @extends {BaseAPI}
 */
export class DepositOfferApi extends BaseAPI {
    /**
     * 
     * @summary Adds a signature mapped to a deposit offer id and an address
     * @param {DtoAddSignatureArgs} addSignatureArgs The input parameters for the multisig transaction
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DepositOfferApi
     */
    public addSignature(addSignatureArgs: DtoAddSignatureArgs, options?: AxiosRequestConfig) {
        return DepositOfferApiFp(this.configuration).addSignature(addSignatureArgs, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Retrieves all signatures for an address only for authorized calls.
     * @param {string} address Address for which to retrieve all signatures
     * @param {string} signature Signature for the request
     * @param {string} timestamp Timestamp for the request
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DepositOfferApi
     */
    public getSignatures(address: string, signature: string, timestamp: string, options?: AxiosRequestConfig) {
        return DepositOfferApiFp(this.configuration).getSignatures(address, signature, timestamp, options).then((request) => request(this.axios, this.basePath));
    }
}



/**
 * MultisigApi - axios parameter creator
 * @export
 */
export const MultisigApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Cancel a multisig transaction by setting the expiration to the current time
         * @param {DtoCancelTxArgs} cancelTxArgs CancelTxArgs object that contains the parameters for the multisig transaction to be canceled
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        cancelMultisigTx: async (cancelTxArgs: DtoCancelTxArgs, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'cancelTxArgs' is not null or undefined
            assertParamExists('cancelMultisigTx', 'cancelTxArgs', cancelTxArgs)
            const localVarPath = `/multisig/cancel`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(cancelTxArgs, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Create a new multisig transaction
         * @param {DtoMultisigTxArgs} multisigTxArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createMultisigTx: async (multisigTxArgs: DtoMultisigTxArgs, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'multisigTxArgs' is not null or undefined
            assertParamExists('createMultisigTx', 'multisigTxArgs', multisigTxArgs)
            const localVarPath = `/multisig`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(multisigTxArgs, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Retrieves all multisig transactions for a given alias
         * @param {string} alias Alias of the multisig account
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getAllMultisigTxForAlias: async (alias: string, signature: string, timestamp: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'alias' is not null or undefined
            assertParamExists('getAllMultisigTxForAlias', 'alias', alias)
            // verify required parameter 'signature' is not null or undefined
            assertParamExists('getAllMultisigTxForAlias', 'signature', signature)
            // verify required parameter 'timestamp' is not null or undefined
            assertParamExists('getAllMultisigTxForAlias', 'timestamp', timestamp)
            const localVarPath = `/multisig/{alias}`
                .replace(`{${"alias"}}`, encodeURIComponent(String(alias)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (signature !== undefined) {
                localVarQueryParameter['signature'] = signature;
            }

            if (timestamp !== undefined) {
                localVarQueryParameter['timestamp'] = timestamp;
            }


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Issue a new multisig transaction
         * @param {DtoIssueTxArgs} issueTxArgs IssueTxArgs object that contains the parameters for the multisig transaction to be issued
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        issueMultisigTx: async (issueTxArgs: DtoIssueTxArgs, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'issueTxArgs' is not null or undefined
            assertParamExists('issueMultisigTx', 'issueTxArgs', issueTxArgs)
            const localVarPath = `/multisig/issue`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(issueTxArgs, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Signs a multisig transaction
         * @param {string} id Multisig transaction ID
         * @param {DtoSignTxArgs} signTxArgs Signer details
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        signMultisigTx: async (id: string, signTxArgs: DtoSignTxArgs, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            assertParamExists('signMultisigTx', 'id', id)
            // verify required parameter 'signTxArgs' is not null or undefined
            assertParamExists('signMultisigTx', 'signTxArgs', signTxArgs)
            const localVarPath = `/multisig/{id}`
                .replace(`{${"id"}}`, encodeURIComponent(String(id)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'PUT', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(signTxArgs, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * MultisigApi - functional programming interface
 * @export
 */
export const MultisigApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = MultisigApiAxiosParamCreator(configuration)
    return {
        /**
         * 
         * @summary Cancel a multisig transaction by setting the expiration to the current time
         * @param {DtoCancelTxArgs} cancelTxArgs CancelTxArgs object that contains the parameters for the multisig transaction to be canceled
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async cancelMultisigTx(cancelTxArgs: DtoCancelTxArgs, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<void>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.cancelMultisigTx(cancelTxArgs, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * 
         * @summary Create a new multisig transaction
         * @param {DtoMultisigTxArgs} multisigTxArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createMultisigTx(multisigTxArgs: DtoMultisigTxArgs, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<ModelMultisigTx>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.createMultisigTx(multisigTxArgs, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * 
         * @summary Retrieves all multisig transactions for a given alias
         * @param {string} alias Alias of the multisig account
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getAllMultisigTxForAlias(alias: string, signature: string, timestamp: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelMultisigTx>>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getAllMultisigTxForAlias(alias, signature, timestamp, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * 
         * @summary Issue a new multisig transaction
         * @param {DtoIssueTxArgs} issueTxArgs IssueTxArgs object that contains the parameters for the multisig transaction to be issued
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async issueMultisigTx(issueTxArgs: DtoIssueTxArgs, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<DtoIssueTxResponse>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.issueMultisigTx(issueTxArgs, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * 
         * @summary Signs a multisig transaction
         * @param {string} id Multisig transaction ID
         * @param {DtoSignTxArgs} signTxArgs Signer details
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async signMultisigTx(id: string, signTxArgs: DtoSignTxArgs, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<ModelMultisigTx>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.signMultisigTx(id, signTxArgs, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
    }
};

/**
 * MultisigApi - factory interface
 * @export
 */
export const MultisigApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = MultisigApiFp(configuration)
    return {
        /**
         * 
         * @summary Cancel a multisig transaction by setting the expiration to the current time
         * @param {DtoCancelTxArgs} cancelTxArgs CancelTxArgs object that contains the parameters for the multisig transaction to be canceled
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        cancelMultisigTx(cancelTxArgs: DtoCancelTxArgs, options?: any): AxiosPromise<void> {
            return localVarFp.cancelMultisigTx(cancelTxArgs, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Create a new multisig transaction
         * @param {DtoMultisigTxArgs} multisigTxArgs The input parameters for the multisig transaction
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createMultisigTx(multisigTxArgs: DtoMultisigTxArgs, options?: any): AxiosPromise<ModelMultisigTx> {
            return localVarFp.createMultisigTx(multisigTxArgs, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Retrieves all multisig transactions for a given alias
         * @param {string} alias Alias of the multisig account
         * @param {string} signature Signature for the request
         * @param {string} timestamp Timestamp for the request
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getAllMultisigTxForAlias(alias: string, signature: string, timestamp: string, options?: any): AxiosPromise<Array<ModelMultisigTx>> {
            return localVarFp.getAllMultisigTxForAlias(alias, signature, timestamp, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Issue a new multisig transaction
         * @param {DtoIssueTxArgs} issueTxArgs IssueTxArgs object that contains the parameters for the multisig transaction to be issued
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        issueMultisigTx(issueTxArgs: DtoIssueTxArgs, options?: any): AxiosPromise<DtoIssueTxResponse> {
            return localVarFp.issueMultisigTx(issueTxArgs, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Signs a multisig transaction
         * @param {string} id Multisig transaction ID
         * @param {DtoSignTxArgs} signTxArgs Signer details
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        signMultisigTx(id: string, signTxArgs: DtoSignTxArgs, options?: any): AxiosPromise<ModelMultisigTx> {
            return localVarFp.signMultisigTx(id, signTxArgs, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * MultisigApi - object-oriented interface
 * @export
 * @class MultisigApi
 * @extends {BaseAPI}
 */
export class MultisigApi extends BaseAPI {
    /**
     * 
     * @summary Cancel a multisig transaction by setting the expiration to the current time
     * @param {DtoCancelTxArgs} cancelTxArgs CancelTxArgs object that contains the parameters for the multisig transaction to be canceled
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MultisigApi
     */
    public cancelMultisigTx(cancelTxArgs: DtoCancelTxArgs, options?: AxiosRequestConfig) {
        return MultisigApiFp(this.configuration).cancelMultisigTx(cancelTxArgs, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Create a new multisig transaction
     * @param {DtoMultisigTxArgs} multisigTxArgs The input parameters for the multisig transaction
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MultisigApi
     */
    public createMultisigTx(multisigTxArgs: DtoMultisigTxArgs, options?: AxiosRequestConfig) {
        return MultisigApiFp(this.configuration).createMultisigTx(multisigTxArgs, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Retrieves all multisig transactions for a given alias
     * @param {string} alias Alias of the multisig account
     * @param {string} signature Signature for the request
     * @param {string} timestamp Timestamp for the request
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MultisigApi
     */
    public getAllMultisigTxForAlias(alias: string, signature: string, timestamp: string, options?: AxiosRequestConfig) {
        return MultisigApiFp(this.configuration).getAllMultisigTxForAlias(alias, signature, timestamp, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Issue a new multisig transaction
     * @param {DtoIssueTxArgs} issueTxArgs IssueTxArgs object that contains the parameters for the multisig transaction to be issued
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MultisigApi
     */
    public issueMultisigTx(issueTxArgs: DtoIssueTxArgs, options?: AxiosRequestConfig) {
        return MultisigApiFp(this.configuration).issueMultisigTx(issueTxArgs, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Signs a multisig transaction
     * @param {string} id Multisig transaction ID
     * @param {DtoSignTxArgs} signTxArgs Signer details
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MultisigApi
     */
    public signMultisigTx(id: string, signTxArgs: DtoSignTxArgs, options?: AxiosRequestConfig) {
        return MultisigApiFp(this.configuration).signMultisigTx(id, signTxArgs, options).then((request) => request(this.axios, this.basePath));
    }
}



