# Hyperledger Fabric Chaincode for Insurance Use-Case

This project implements a Hyperledger Fabric chaincode for managing **Accident Reports** and **Insurance Policies**. It provides functionalities for creating, reading, deleting, and approving accident reports and insurance policies, along with access control mechanisms.

## Overview

This chaincode allows for the management of **Accident Reports** and **Insurance Policies** on a Hyperledger Fabric network. 

### Key Features:
1. **Accident Reports**:
   - Create, read, delete accident reports.
   - View the history of accident reports.
   - Query accident reports with pagination.

2. **Insurance Policies**:
   - Create and approve insurance policies.
   - Access the insurance policy status.
   - Control access to these functionalities based on the user's role.


## Steps to Deploy the chaincode using Minifab

#
```
sudo chmod -R 777 vars/
```
```
mkdir -p vars/chaincode/Insurance/go
```
```
cp -r ../Chaincode/* vars/chaincode/Insurance/go/
```
```
minifab ccup -n Insurance -l go -v 1.0 -d false -r false
```

#### Functions for Accident Report Contract
#### 1. CreateAccidentReport
```
minifab invoke -n Insurance -p '"CreateAccidentReport","report01","01/01/2024","AccidentReport","1234","Car"'
```
#### 2. ReadAccidentReport
```
minifab query -n Insurance -p '"ReadAccidentReport","report01"'
```
#### 3. DeleteAccidentReport
```
minifab invoke -n Insurance -p '"DeleteAccidentReport","report01"' -o government.insuranceclaim.com
```
#### 4. GetAccidentReportByRange
```
minifab query -n Insurance -p '"GetAccidentReportByRange","report01","report03"'
```
#### 5. GetAllAccidentReports
```
minifab query -n Insurance -p '"GetAllAccidentReports"'
```
#### 6. GetAccidentReportHistory
```
minifab query -n Insurance -p '"GetAccidentReportHistory","report01"'
```
#### 7. GetAccidentReportsWithPagination
```
minifab query -n Insurance -p '"GetAccidentReportsWithPagination","3",""'
```
#

#### Functions for Insurance Policy Contract created using PDC
#### 1. CreateInsurancePolicy
```
minifab invoke -n Insurance -p '"insurancePolicyContract:CreateInsurancePolicy","policy01","Amjad","Accident","100000"'
```
#### 2. ReadInsurancePolicy
```
minifab invoke -n Insurance -p   '"InsurancePolicyContract:ReadInsurancePolicy","policy01"'
```
#### 3. DeleteInsurancePolicy
```
minifab invoke -n Insurance -p   '"InsurancePolicyContract:DeleteInsurancePolicy","policy01"'
```
#### 4. GetInsurancePoliciesByRange
```
minifab query -n Insurance -p '"InsurancePolicyContract:GetInsurancePolicyByRange","policy01","policy03"'
```
#### 5. GetAllInsurancePolicies
```
minifab query -n Insurance -p '"InsurancePolicyContract:GetAllInsurancePolicies"'
```



