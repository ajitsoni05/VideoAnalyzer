package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoFields struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Key    string             `json:"key" bson:"key"`
	Status string             `bson:"status" json:"status"`
}

type AnalysisResult struct {
	Key          string `json:"key" bson:"key"`
	PlatformName string `json:"platform_name" bson:"platform_name"`
	AssetID      string `json:"asset_id" bson:"asset_id"`
	FilePath     string `json:"file_path" bson:"file_path"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp"`
	Data         string `json:"data" bson:"data"`
	Status       string `json:"status" bson:"status"`
}

type GetStatusTestInputs struct {
	Key            string
	ExpectedStatus string
}

type AnalyzeVideoTestInputs struct {
	AnalyzeVideoReqBody AnalyzeVideoReqBody
	ExpectedStatus      string
}
type AnalyzeVideoReqBody struct {
	PlatformName string `json:"platform_name" bson:"platform_name"`
	AssetID      string `json:"asset_id" bson:"asset_id"`
	FilePath     string `json:"file_path" bson:"file_path"`
}
