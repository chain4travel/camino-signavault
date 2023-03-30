# Signavault
Signavault is a standalone off chain service which is responsible for collecting signatures for multisig alias transactions. It provides an API with four endpoits:

 - `CreateMultisigTx`: creates a new multisig transaction.
 - `GetAllMultisigTxForAlias`: gets all the pending multisig transactions for a given alias.
 - `SignMultisigTx`: signs an already existing multisig transaction.
 - `IssueMultisigTx`: issues a multisig transaction to the network if threshold of signatures is reached.

# Requirements
To run Signavault locally either you need `docker-compose` installed or you could set up `mysql` and the migration scripts manually. .
In addition, Signavault is dependent to a Camino network. You can either run a local network or connect to a remote one.

# Installation
To install Signavault, follow these steps:

- Clone the repository: `git clone https://github.com/chain4travel/camino-signavault`.
- Open the `config.yml` file and fill in the appropriate values:
  - `listenerAddress`: the address where the service will listen for requests (e.g., `:8080`).
  - `caminoNode`: the URL of the Camino node that signavault will connect to (e.g., `http://localhost:9650`).
  - `networkID`: the ID of the running Camino network (e.g., `1002`).
  - `database.dsn`: the connection string for the database (e.g., `root:password@tcp(mysql:3306)/signavault?parseTime=true`).
- Go to the `docker/local` directory: `cd docker/local`.
- Run `docker-compose up`. This will start the database and the migration scripts.
- In a new terminal window, go to the `cmd/camino-signavault` directory.
- Run `go run main.go`. This will start the Signavault service.

# Usage
Once Signavault is running, you can use the API endpoints to create, sign, and issue multisignature transactions. The API documentation is available at https://c4t.atlassian.net/wiki/spaces/TECH/pages/292257793/SignaVault+API.

# Client SDK
Signavault also provides a TypeScript client SDK that can be used in front-end apps to communicate with the Signavault API. The SDK is available in the `signavaultjs` directory and can be installed as an npm package:
`npm install @c4tplatform/signavaultjs`. The SDK implements all Signavault endpoints and provides TypeScript types for the API responses.

# Examples 
Signavault provides frontend examples that demonstrate how to use the TypeScript client SDK to interact with the Signavault API. The examples are available in the `examples` directory and can be run with the following steps:
- Go to the `examples` directory: `cd examples`.
- Install the dependencies: `npm install`
- Go to the `examples/dependencies/caminojs` directory: `cd examples/dependencies/caminojs`. This is the Camino go client which is also a dependency of the examples.
- Install and build the dependency: `npm install`
- Run the example: `npx ts-node -r tsconfig-paths/register addValidatorTx.ts`.

Currently, there is one example which demonstrates the full flow of a multisignature transaction using all endpoints and simulating a multisig alias with two users. The example demonstrates the registerNodeTx and addValidatorTx endpoints.
