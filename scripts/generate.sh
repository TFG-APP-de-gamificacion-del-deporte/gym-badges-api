#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
DIRS=("restapi/operations" "models" "mocks")
FILES=("restapi/doc.go" "restapi/embedded_spec.go" "restapi/server.go")
GOLANGCI_LINT_URL="github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0"
MOCKGEN_URL="go.uber.org/mock/mockgen@latest"

CMD_NOT_FOUND_ERR="${YELLOW}Command not found.${NC} Trying .exe extension..."

echo "Cleaning project..."

for dir in "${DIRS[@]}"; do
    if [ -d "$dir" ]; then
        rm -rf "$dir"
    fi
done

for file in "${FILES[@]}"; do
    rm -f "$file"
done


# ===== SWAGGER =====
echo -e "${GREEN}Building swagger...${NC}"

swagger generate server -f ./swagger.yml --exclude-main -A gym-badges
# If command not found, assume the script is running in WSL and add ".exe"
if [[ $? -ne 0 ]]; then
  echo -e "${CMD_NOT_FOUND_ERR}"
  swagger.exe generate server -f ./swagger.yml --exclude-main -A gym-badges
fi
echo -e "${GREEN}Done!\n${NC}"


# ===== MOCKS =====
echo -e "${GREEN}Building mocks...${NC}"

interfaces=$(find . -type f -name "*_interface.go")

if ! mockgen --version 2> /dev/null; then
  echo -e "${GREEN}Installing mockgen...${NC}"

  go install $MOCKGEN_URL
  # If command not found, assume the script is running in WSL and add ".exe"
  if [[ $? -ne 0 ]]; then
    echo -e "${CMD_NOT_FOUND_ERR}"
    go.exe install $MOCKGEN_URL
  fi
  echo -e "${GREEN}Done!\n${NC}"
fi

for interface in $interfaces; do
    echo -e "${GREEN}Building mock for: $interface ${NC}"
    package_name=$(basename "$interface" | cut -d'_' -f2)
    mockgen -source="$interface" -destination="mocks/$package_name/mock_$(basename "$interface")" -package="$package_name"

    # If command not found, assume the script is running in WSL and add ".exe"
    if [[ $? -ne 0 ]]; then
      echo -e "${CMD_NOT_FOUND_ERR}"
      mockgen.exe -source="$interface" -destination="mocks/$package_name/mock_$(basename "$interface")" -package="$package_name"
    fi
    echo -e "${GREEN}Done!\n${NC}"
done


# ===== GOLANGCI-LINT =====
if ! golangci-lint --version 2> /dev/null; then
  echo -e "${GREEN}Installing golangci-lint...${NC}"

  go install $GOLANGCI_LINT_URL
  # If command not found, assume the script is running in WSL and add ".exe"
  if [[ $? -ne 0 ]]; then
    echo -e "${CMD_NOT_FOUND_ERR}"
    go.exe install $GOLANGCI_LINT_URL
  fi
  echo -e "${GREEN}Done!\n${NC}"
fi

echo -e "${GREEN}Executing golangci-lint...${NC}"

golangci-lint run --config golangci-lint.yml
# If command not found, assume the script is running in WSL and add ".exe"
if [[ $? -ne 0 ]]; then
  echo -e "${CMD_NOT_FOUND_ERR}"
  golangci-lint.exe run --config golangci-lint.yml
fi
echo -e "${GREEN}Done!\n${NC}"

echo -e "${GREEN}SUCCESS...${NC}"
