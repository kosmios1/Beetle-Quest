import yaml
import json

file_path = "./openAPI.yaml"
with open(file_path, "r") as file:
    openapi_spec = yaml.safe_load(file)

def generate_postman_collection(spec):
    collection = {
        "info": {
            "name": spec.get("info", {}).get("title", "Generated Postman Collection"),
            "description": spec.get("info", {}).get("description", "A collection generated from the OpenAPI specification."),
            "_postman_id": "unique-id-placeholder",
            "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
        },
        "item": []
    }

    paths = spec.get("paths", {})
    for path, methods in paths.items():
        for method, details in methods.items():
            item = {
                "name": details.get("summary", path),
                "request": {
                    "method": method.upper(),
                    "header": [],
                    "url": {
                        "raw": "{{baseUrl}}" + path,
                        "host": ["{{baseUrl}}"],
                        "path": path.lstrip("/").split("/")
                    },
                    "description": details.get("description", ""),
                }
            }
            # Add parameters
            params = details.get("parameters", [])
            query_params = [
                {"key": param["name"], "value": "", "description": param.get("description", ""), "disabled": not param.get("required", False)}
                for param in params if param.get("in") == "query"
            ]
            if query_params:
                item["request"]["url"]["query"] = query_params

            # Add body if applicable
            if "requestBody" in details:
                content = details["requestBody"].get("content", {})
                if "application/json" in content:
                    item["request"]["body"] = {
                        "mode": "raw",
                        "raw": json.dumps(content["application/json"].get("example", {}), indent=2),
                        "options": {"raw": {"language": "json"}}
                    }

            collection["item"].append(item)

    return collection

postman_collection = generate_postman_collection(openapi_spec)

output_file = "./generated_postman_beetle_quest.json"
with open(output_file, "w") as file:
    json.dump(postman_collection, file, indent=2)
