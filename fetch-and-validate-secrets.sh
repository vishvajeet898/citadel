#!/bin/bash
# Define the YAML file path
FILE="./config/prod/prod.yaml"  # Change this to your actual YAML file path in python project
CONFIG_FILE="/go/src/github.com/Orange-Health/citadel/config/config.yaml"
SERVICE_NAME="citadel"


if [ -n "$RDS_DB_CREDS" ]; then
    RDS_DB_CREDS=${RDS_DB_CREDS}
fi

if [ -n "$OTHER_DB_CREDS" ]; then
    OTHER_DB_CREDS=${OTHER_DB_CREDS}
fi

if [ -n "$KEY_VALUE_CREDS" ]; then
    KEY_VALUE_CREDS=${KEY_VALUE_CREDS}
fi


aws configure set aws_access_key_id ${S3_ACCESS_KEY}
aws configure set aws_secret_access_key ${S3_SECRET_KEY}
aws configure set region ${AWS_REGION}

# if env is prod, then use the prod secrets
if [ "$ENV" = "prod" ]; then
    cp /go/src/github.com/Orange-Health/${SERVICE_NAME}/config/prod/${ENV}.yaml /go/src/github.com/Orange-Health/${SERVICE_NAME}/config/config.yaml
else
    cp /go/src/github.com/Orange-Health/${SERVICE_NAME}/config/stag/${ENV}.yaml /go/src/github.com/Orange-Health/${SERVICE_NAME}/config/config.yaml
fi

if [ -n "$RDS_DB_CREDS" ]; then
    SECRET_JSON=$(aws secretsmanager get-secret-value --secret-id ${RDS_DB_CREDS} --region ${AWS_REGION} | jq -r '.SecretString')

    if [ -z "$SECRET_JSON" ]; then
        echo "Failed to retrieve secret ${RDS_DB_CREDS}"
        exit 1
    fi


    DB_HOST=$(echo $SECRET_JSON | jq -r '.host')
    DB_USER=$(echo $SECRET_JSON | jq -r '.username')
    DB_PASSWORD=$(echo $SECRET_JSON | jq -r '.password')

    # CREATE THE PLACEHOLDER USING SERVICE NAME
    sed -i "s|{{${ENV}_${SERVICE_NAME}_DbHost}}|${DB_HOST}|g" "$CONFIG_FILE"
    sed -i "s|{{${ENV}_${SERVICE_NAME}_DbUser}}|${DB_USER}|g" "$CONFIG_FILE"
    sed -i "s|{{${ENV}_${SERVICE_NAME}_DbPswd}}|${DB_PASSWORD}|g" "$CONFIG_FILE"
fi

# EXTRACT OTHER DB CREDENTIALS VALUES SEPARATED BY , AND ITERATE OVER THEM
if [ -n "$OTHER_DB_CREDS" ]; then
    echo "$OTHER_DB_CREDS" | tr ',' '\n' | while read -r element; do
        temp=${element#*_} # remove prefix
        temp=${temp%_*} # remove suffix

        SECRET_JSON=$(aws secretsmanager get-secret-value --secret-id ${element} --region ${AWS_REGION} | jq -r '.SecretString')

        if [ -z "$SECRET_JSON" ]; then
            echo "Failed to retrieve secret ${element}"
            exit 1
        fi

        DB_HOST=$(echo $SECRET_JSON | jq -r '.host')
        DB_USER=$(echo $SECRET_JSON | jq -r '.username')
        DB_PASSWORD=$(echo $SECRET_JSON | jq -r '.password')

        if [ "$ENV" = "prod" ]; then
            sed -i "s|{{${ENV}_${temp}_DbHost}}|$DB_HOST|g" "$CONFIG_FILE"
            sed -i "s|{{${ENV}_${temp}_DbUser}}|$DB_USER|g" "$CONFIG_FILE"
            sed -i "s|{{${ENV}_${temp}_DbPswd}}|$DB_PASSWORD|g" "$CONFIG_FILE"
        else
            sed -i "s|{{${CLUSTER_NAME}_${temp}_DbHost}}|$DB_HOST|g" "$CONFIG_FILE"
            sed -i "s|{{${CLUSTER_NAME}_${temp}_DbUser}}|$DB_USER|g" "$CONFIG_FILE"
            sed -i "s|{{${CLUSTER_NAME}_${temp}_DbPswd}}|$DB_PASSWORD|g" "$CONFIG_FILE"
        fi
    done
fi


# check if KEY_VALUE_CREDS is set
if [ -n "$KEY_VALUE_CREDS" ]; then
    # EXTRACT KEY VALUE CREDENTIALS VALUES SEPERATED BY , AND ITERATE OVER THEM
    echo "$KEY_VALUE_CREDS" | tr ',' '\n' | while read -r element; do

        SECRET_JSON=$(aws secretsmanager get-secret-value --secret-id ${element} --region ${AWS_REGION} | jq -r '.SecretString')
        if [ -z "$SECRET_JSON" ]; then
            echo "Failed to retrieve secret ${element}"
            exit 1
        fi

        # loop through the key value pairs and get the key and value
        echo "$SECRET_JSON" | jq -r 'to_entries|map("\(.key)=\(.value|tostring)")|.[]' | while IFS='=' read -r KEY VALUE; do
            # replace the key with the value
            sed -i "s|{{${KEY}}}|$(printf '%s' "$VALUE" | sed 's/[\/&]/\\&/g')|g" "$CONFIG_FILE"
        done
    done
fi

if [ "$ENV" = "ephemeral" ]; then
    sed -i "s|{{ephemeral_env}}|$NAMESPACE|g" "$CONFIG_FILE"
    sed -i "s|{{ephemeral_postgres_DbHost}}|$PSQL_DB_HOST|g" "$CONFIG_FILE"
    sed -i "s|{{ephemeral_postgres_DbUser}}|$PSQL_DB_USER|g" "$CONFIG_FILE"
    sed -i "s|{{ephemeral_postgres_DbPswd}}|$PSQL_DB_PASSWORD|g" "$CONFIG_FILE"
fi

if grep -q "{{" "$CONFIG_FILE"; then
    echo "One or more placeholders are still present in the config file:"
    grep -o "{{[^}]*}}" "$CONFIG_FILE"
    exit 1
fi

echo "Updated ${ENV} config.yaml with secrets!"
