# Remote execution

The remote execution offering works pretty much the same way. You need to define your transactions, as you would with the CLI, and then using the APId CLI issue a `remote` command. This works as follows:

- The CLI reads the transactions from the specified directory or file, as per ususal
- It will execute each transaction, sending the necessary `STEP` information to the cloud for remote execution
- It will wait for each `STEP` result to come back and continue with the other steps and transactions

One major benefit of this workflow is that all the interpolation is done locally (on the machine running the CLI), thus you have control over the environment it runs in. This means you can invoke any custom executables.

## Usage

In order to use the power of the cloud you will need a personal access key. To generate one, you will have to:

- Head over to https://console.getapid.com and sign up
- Go to the dashboard and create a new access key

Once you have your key you will need to [install the APId CLI](../installation/cli.md) (if you haven't already) or use our [official docker image](../installation/docker.md).

A reference on how use the CLI after installation for remote execution can be found [here](../cli/cloud/check.md).

## Timeouts

The execution timeout is set to 30 seconds. If your API does not respond within that time an error is returned.
