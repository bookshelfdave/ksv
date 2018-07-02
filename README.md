# ksv


 decodes/encodes entire Kubernetes secrets files


# Installation

ksv uses [dep](https://github.com/golang/dep) to manage dependencies.

    go get github.com/metadave/ksv
    cd ${GOPATH}/src/github.com/metadave/ksv
    dep ensure
    go install
    # ksv will be installed in ${GOPATH}/bin
    


# Usage

### Base64 decoding secret values

    ksv < some_secrets_file_with_base64_encoded_data_values.yaml

or
    
    ksv decode < some_secrets_file_with_base64_encoded_data_values.yaml

> the default subcommand for ksv is `decode`

### Convert base64-encoded secret values to use K8s `stringData`

    ksv -s < some_secrets_file_with_base64_encoded_data_values.yaml


### Base64 encoding secret values

    ksv encode < some_secrets_file_with_plaintext_data_values.yaml

### Round trip

    ksv < test.yaml | ksv encode



# License

[Apache Software License 2.0](https://github.com/metadave/ksv/blob/master/LICENSE)
