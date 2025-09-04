import yaml
import re
from collections import defaultdict
import os
import boto3

def extract_placeholders_from_yaml(file_path):
    with open(file_path, "r") as file:
        yaml_data = yaml.safe_load(file)

    placeholders = defaultdict(int)
    pattern = re.compile(r"{{(.*?)}}")

    def find_placeholders(data):
        if isinstance(data, dict):
            for value in data.values():
                find_placeholders(value)
        elif isinstance(data, list):
            for item in data:
                find_placeholders(item)
        elif isinstance(data, str):
            matches = pattern.findall(data)
            for match in matches:
                if "_Db" not in match:
                    placeholders[match] = 0

    find_placeholders(yaml_data)
    return placeholders

def fetch_secrets():
    secret_keys = [key.strip() for key in os.environ.get("SECRET_KEYS", "").split(",")]
    secrets = {}

    for key in secret_keys:
        if key:
            secret_value = fetch_secret_from_manager(key)
            if secret_value:
                secret_dict = yaml.safe_load(secret_value)
                for secret_key in secret_dict.keys():
                    secrets[secret_key] = 0

    return secrets

def fetch_secret_from_manager(secret_name):
    client = boto3.client('secretsmanager')
    try:
        response = client.get_secret_value(SecretId=secret_name)
        if 'SecretString' in response:
            return response['SecretString']
        else:
            return None
    except Exception as e:
        print(f"Error fetching secret for {secret_name}: {e}")
        return None

if __name__ == "__main__":
    file_path = "config/prod/prod.yaml"
    extracted_placeholders = extract_placeholders_from_yaml(file_path)

    secrets = fetch_secrets()

    for key in secrets.keys():
        if key in extracted_placeholders:
            extracted_placeholders[key] = 1

    missing_secrets = [key for key, value in extracted_placeholders.items() if value == 0]

    print("Missing secrets:", missing_secrets)

    if missing_secrets:
        exit(1)
    else:
        exit(0)
