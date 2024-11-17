import yaml
import json
from collections import defaultdict

# Load the OpenAPI file
openapi_file_path = './openAPI.yaml'
output_file = './postman_collection.json'

with open(openapi_file_path, 'r') as file:
    openapi_data = yaml.safe_load(file)

# Helper function to convert schema to example JSON
def schema_to_example(schema):
    if schema.get("type") == "object":
        return {key: "" for key in schema.get("properties", {}).keys()}
    return {}

# Helper function to extract JSON body schema from components
def extract_schema(schema_ref, components):
    if not schema_ref.startswith('#/components/schemas/'):
        return {}
    schema_name = schema_ref.split('/')[-1]
    schema = components.get(schema_name, {})
    return schema_to_example(schema)

# Helper function to extract parameters
def extract_parameters(parameters):
    query_params = []
    path_params = []
    headers = []

    for param in parameters:
        param_type = param.get("in", "")
        param_example = param.get("example", "") or ""
        param_name = param.get("name", "")

        if param_type == "query":
            query_params.append({
                "key": param_name,
                "value": param_example,
                "description": param.get("description", ""),
                "disabled": False
            })
        elif param_type == "path":
            path_params.append(param_example)  # Used directly in URL generation
        elif param_type == "header":
            headers.append({
                "key": param_name,
                "value": param_example,
                "description": param.get("description", ""),
                "disabled": False
            })

    return query_params, path_params, headers

# Prepare the Postman Collection
postman_collection = {
    "info": {
        "name": openapi_data.get("info", {}).get("title", "Generated Collection"),
        "description": openapi_data.get("info", {}).get("description", ""),
        "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
    },
    "item": []
}

# Components schemas
components = openapi_data.get("components", {}).get("schemas", {})

# Group endpoints by tags
endpoints_by_tag = defaultdict(list)

for path, methods in openapi_data.get("paths", {}).items():
    for method, details in methods.items():
        tags = details.get("tags", ["General"])
        for tag in tags:
            endpoints_by_tag[tag].append({
                "path": path,
                "method": method.upper(),
                "summary": details.get("summary", ""),
                "description": details.get("description", ""),
                "parameters": details.get("parameters", []),
                "requestBody": details.get("requestBody", {}),
            })

# Create folders for tags and add requests
for tag, endpoints in endpoints_by_tag.items():
    folder = {
        "name": tag,
        "item": []
    }
    for endpoint in endpoints:
        # Extract parameters
        query_params, path_params, headers = extract_parameters(endpoint["parameters"])

        # Handle request body
        body = {}
        if endpoint["requestBody"]:
            content = endpoint["requestBody"].get("content", {})
            for content_type, content_details in content.items():
                if content_type == "application/json":
                    schema_ref = content_details.get("schema", {}).get("$ref", "")
                    if schema_ref:
                        body = extract_schema(schema_ref, components)

        # Construct Postman request
        request = {
            "name": f"{endpoint['method']} {endpoint['path']}",
            "request": {
                "method": endpoint["method"],
                "header": headers,
                "body": {
                    "mode": "raw",
                    "raw": json.dumps(body, indent=4)
                } if body else {},
                "url": {
                    "raw": "{{baseUrl}}" + endpoint["path"],
                    "host": ["{{baseUrl}}"],
                    "path": endpoint["path"].lstrip('/').split('/'),
                    "query": query_params,
                    "variable": [{"key": param, "value": param} for param in path_params]
                },
                "description": endpoint["description"]
            }
        }
        folder["item"].append(request)
    postman_collection["item"].append(folder)

# Save Postman collection to file
with open(output_file, 'w') as file:
    json.dump(postman_collection, file, indent=4)

print(f"Postman collection saved at: {output_file}")
