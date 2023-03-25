package cloudwatch

import (
	"fmt"

	"github.com/Zeus-Labs/ZeusCloud/rules/types"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type VPCChanges struct{}

func (VPCChanges) UID() string {
	return "cloudwatch/vpc_changes"
}

func (VPCChanges) Description() string {
	return "Log metric filter and alarm should exist for VPC changes"
}

func (VPCChanges) Severity() types.Severity {
	return types.Low
}

func (VPCChanges) RiskCategories() types.RiskCategoryList {
	return []types.RiskCategory{
		types.InsufficientMonitoring,
	}
}

func (VPCChanges) Execute(tx neo4j.Transaction) ([]types.Result, error) {
	records, err := tx.Run(
		`MATCH (a:AWSAccount{inscope: true})
		OPTIONAL MATCH (a)-[:RESOURCE]->(t:CloudTrail)-[:DELIVERS_TO]->(lg:CloudWatchLogGroup)-[:HAS_METRIC_FILTER]->(mf:CloudWatchMetricFilter)
		WHERE t.is_multi_region_trail AND
		t.is_logging AND
		mf.filter_pattern =~ '.*\$\.eventName\s*=\s*"?CreateVpc"?.+\$\.eventName\s*=\s*"?DeleteVpc"?.+\$\.eventName\s*=\s*"?ModifyVpcAttribute"?.+\$\.eventName\s*=\s*"?AcceptVpcPeeringConnection"?.+\$\.eventName\s*=\s*"?CreateVpcPeeringConnection"?.+\$\.eventName\s*=\s*"?DeleteVpcPeeringConnection"?.+\$\.eventName\s*=\s*"?RejectVpcPeeringConnection"?.+\$\.eventName\s*=\s*"?AttachClassicLinkVpc"?.+\$\.eventName\s*=\s*"?DetachClassicLinkVpc"?.+\$\.eventName\s*=\s*"?DisableVpcClassicLink"?.+\$\.eventName\s*=\s*"?EnableVpcClassicLink"?.*'
		OPTIONAL MATCH (al:CloudWatchAlarm)
		WHERE al.metric_name IS NOT NULL
		AND mf.metric_transformation_name = al.metric_name
		WITH a,
		COLLECT(CASE WHEN t IS NOT NULL AND lg IS NOT NULL AND mf IS NOT NULL AND al IS NOT NULL THEN t ELSE NULL END) as valid_ts,
		COLLECT(CASE WHEN t IS NOT NULL AND lg IS NOT NULL AND mf IS NOT NULL AND al IS NOT NULL THEN lg ELSE NULL END) as valid_lgs,
		COLLECT(CASE WHEN t IS NOT NULL AND lg IS NOT NULL AND mf IS NOT NULL AND al IS NOT NULL THEN mf ELSE NULL END) as valid_mfs,
		COLLECT(CASE WHEN t IS NOT NULL AND lg IS NOT NULL AND mf IS NOT NULL AND al IS NOT NULL THEN al ELSE NULL END) as valid_als,
		SUM(CASE WHEN t IS NOT NULL AND lg IS NOT NULL AND mf IS NOT NULL AND al IS NOT NULL THEN 1 ELSE 0 END) as num_valid
		RETURN a.id as resource_id,
		'AWSAccount' as resource_type,
		a.id as account_id,
		CASE WHEN 
			num_valid = 0 THEN 'failed'
			ELSE 'passed'
		END as status,
		CASE WHEN 
			num_valid = 0 THEN 'There is no log metric filter and alarm monitoring an active, multi-region CloudTrail trail for VPC changes.'
			ELSE 'The metric filter ' + valid_mfs[0].filter_name + ' and alarm ' + valid_als[0].arn + ' monitor VPC change events.'
		END as context`,
		nil,
	)
	if err != nil {
		return nil, err
	}
	var results []types.Result
	for records.Next() {
		record := records.Record()
		resourceID, _ := record.Get("resource_id")
		resourceIDStr, ok := resourceID.(string)
		if !ok {
			return nil, fmt.Errorf("resource_id %v should be of type string", resourceID)
		}
		resourceType, _ := record.Get("resource_type")
		resourceTypeStr, ok := resourceType.(string)
		if !ok {
			return nil, fmt.Errorf("resource_type %v should be of type string", resourceType)
		}
		accountID, _ := record.Get("account_id")
		accountIDStr, ok := accountID.(string)
		if !ok {
			return nil, fmt.Errorf("account_id %v should be of type string", accountID)
		}
		status, _ := record.Get("status")
		statusStr, ok := status.(string)
		if !ok {
			return nil, fmt.Errorf("status %v should be of type string", status)
		}
		context, _ := record.Get("context")
		contextStr, ok := context.(string)
		if !ok {
			return nil, fmt.Errorf("context %v should be of type string", context)
		}
		results = append(results, types.Result{
			ResourceID:   resourceIDStr,
			ResourceType: resourceTypeStr,
			AccountID:    accountIDStr,
			Status:       statusStr,
			Context:      contextStr,
		})
	}
	return results, nil
}

func (VPCChanges) ProduceRuleGraph(tx neo4j.Transaction, resourceId string) (neo4j.Result, error) {
	return nil, nil
}
