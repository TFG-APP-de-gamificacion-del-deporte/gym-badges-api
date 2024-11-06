#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'
DIRS=("restapi/operations" "models" "mocks")
FILES=("restapi/doc.go" "restapi/embedded_spec.go" "restapi/server.go")
GOLANGCI_LINT_URL="github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0"
MOCKGEN_URL="go.uber.org/mock/mockgen@latest"

# Recieves the command to execute as arguments
try_command () {
  # Try running command as it is
  $*

  # If command not found, assume the script is running in WSL and add ".exe"
  if [[ $? -ne 0 ]]; then
    echo -e "${YELLOW}Command not found. Trying .exe extension...${NC}"
    cmd=$1
    shift 1
    $cmd.exe $*
  fi

  # If command fails, stop script
  if [[ $? -ne 0 ]]; then
    echo -e "${RED}Error on command: ${cmd}.exe $*${YELLOW}"
    exit
  fi

  echo -e "${GREEN}Done!\n${NC}"
}


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

try_command swagger generate server -f ./swagger.yml --exclude-main -A gym-badges


# ===== MOCKS =====
echo -e "${GREEN}Building mocks...${NC}"

interfaces=$(find . -type f -name "*_interface.go")

if ! mockgen --version 2> /dev/null; then
  echo -e "${GREEN}Installing mockgen...${NC}"

  try_command go install $MOCKGEN_URL
fi

for interface in $interfaces; do
    echo -e "${GREEN}Building mock for: $interface ${NC}"
    package_name=$(basename "$interface" | cut -d'_' -f2)

    try_command mockgen -source="$interface" -destination="mocks/$package_name/mock_$(basename "$interface")" -package="$package_name"
done


# ===== GOLANGCI-LINT =====
if ! golangci-lint --version 2> /dev/null; then
  echo -e "${GREEN}Installing golangci-lint...${NC}"

  try_command go install $GOLANGCI_LINT_URL
fi

echo -e "${GREEN}Executing golangci-lint...${NC}"

try_command golangci-lint run --config golangci-lint.yml

echo -e "${GREEN}SUCCESS...${NC}"
