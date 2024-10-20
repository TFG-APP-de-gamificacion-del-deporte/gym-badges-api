#!/bin/bash

GREEN='\033[0;32m'
NC='\033[0m'
DIRS=("restapi/operations" "models" "mocks")
FILES=("restapi/doc.go" "restapi/embedded_spec.go" "restapi/server.go")
GOLANGCI_LINT_URL="github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0"
MOCKGEN_URL="go.uber.org/mock/mockgen@latest"

echo "Cleaning project..."

for dir in "${DIRS[@]}"; do
    if [ -d "$dir" ]; then
        rm -rf "$dir"
    fi
done

for file in "${FILES[@]}"; do
    rm -f "$file"
done

# SWAGGER
echo -e "${GREEN}Building swagger...${NC}"

swagger generate server -f ./swagger.yml --exclude-main -A gym-badges

# MOCKS
echo -e "${GREEN}Building mocks...${NC}"

interfaces=$(find . -type f -name "*_interface.go")

if ! mockgen --version 2> /dev/null; then
  echo -e "${GREEN}Installing mockgen...${NC}"
  go install $MOCKGEN_URL
fi

for interface in $interfaces; do
    echo -e "${GREEN}Building mock for: $interface ${NC}"
    package_name=$(basename "$interface" | cut -d'_' -f2)
    mockgen -source="$interface" -destination="mocks/$package_name/mock_$(basename "$interface")" -package="$package_name"
done

# GOLANGCI-LINT
if ! golangci-lint --version 2> /dev/null; then
  echo -e "${GREEN}Installing golangci-lint...${NC}"
  go install $GOLANGCI_LINT_URL
fi

echo -e "${GREEN}Executing golangci-lint...${NC}"

golangci-lint run --config golangci-lint.yml

echo -e "${GREEN}SUCCESS...${NC}"
