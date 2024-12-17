#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... |
grep -v /billing-loan-system/cmd |
grep -v /billing-loan-system/configs |
grep -v /billing-loan-system/pkg |
grep -v /billing-loan-system/gen |
grep -v /billing-loan-system/db |
grep -v /billing-loan-system/internal/ctx |
grep -v /billing-loan-system/internal/dto |
grep -v /billing-loan-system/tests |
grep -v /billing-loan-system/port |
grep -v /billing-loan-system/internal/dto |
grep -v /billing-loan-system/internal/billing_config/models |
grep -v /billing-loan-system/internal/payment/models |
grep -v /billing-loan-system/internal/loan/models
)

# Remove the coverage files directory, will keep this dir for sonarqube
rm -rf "$COVERAGE_DIR";

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR";

# Create a coverage file for each package
go test -covermode=count -coverprofile "${COVERAGE_DIR}/coverage.cov" ${PKG_LIST} ;

# Merge the coverage profile files
#echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov ;
#tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> "${COVERAGE_DIR}"/coverage.cov ;

# Display the global code coverage
go tool cover -func="${COVERAGE_DIR}"/coverage.cov ;

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o "${COVERAGE_DIR}"/coverage.html;
    open "${COVERAGE_DIR}"/coverage.html;
fi
